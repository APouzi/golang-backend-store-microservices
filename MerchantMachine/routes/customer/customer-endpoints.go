package customer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/APouzi/MerchantMachinee/authorization"
	"github.com/APouzi/MerchantMachinee/routes/helpers"
	"github.com/go-chi/chi/v5"
)

type CustomerRoutes struct {
	authClient *auth.Client
}

func InstanceCustomerRoutes(app *firebase.App) *CustomerRoutes {
	client, err := app.Auth(context.Background())
	if err != nil {
		fmt.Printf("error getting Auth client: %v\n", err)
		return nil
	}
	return &CustomerRoutes{
		authClient: client,
	}
}



func (cr *CustomerRoutes) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("started here")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Try to recover token from body (frontend may send it there when middleware is off)
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // restore body for downstream
		var bodyToken struct {
			Token string `json:"token"`
		}
		_ = json.Unmarshal(bodyBytes, &bodyToken)
		fmt.Printf("No Authorization header, body token=%s\n", bodyToken.Token)
		if bodyToken.Token != "" {
			authHeader = "Bearer " + bodyToken.Token
			r.Header.Set("Authorization", authHeader)
		} else {
			helpers.ErrorJSON(w, fmt.Errorf("missing authorization header"), http.StatusUnauthorized)
			return
		}
	}
	fmt.Printf("Made it past header; Authorization=%s\n", authHeader)

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	fmt.Printf("Made it past tokenstring len=%d\n", len(tokenString))
	token, err := cr.authClient.VerifyIDToken(r.Context(), tokenString)
	fmt.Printf("made it to token: %v\n", token != nil)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid token: %v", err), http.StatusUnauthorized)
		return
	}
	fmt.Println("claims started")
	// Extract user information from claims
	claims := token.Claims
	fmt.Printf("claims keys: ")
	for k := range claims {
		fmt.Printf("%s ", k)
	}
	fmt.Println()
	
	customer := CustomerRegistration{
		UID: token.UID,
	}

	if email, ok := claims["email"].(string); ok {
		customer.Email = email
	}
	// Prefer structured name claims when available (OIDC uses given_name / family_name)
	if given, ok := claims["given_name"].(string); ok {
		fmt.Printf("given_name claim found: %s\n", given)
		customer.FirstName = given
	}
	if family, ok := claims["family_name"].(string); ok {
		fmt.Printf("family_name claim found: %s\n", family)
		customer.LastName = family
	}
	if name, ok := claims["name"].(string); ok {
		fmt.Printf("name claim found: %s\n", name)
		customer.DisplayName = name
		// If first/last were missing, try to split display name
		if customer.FirstName == "" || customer.LastName == "" {
			parts := strings.Fields(name)
			if len(parts) >= 1 && customer.FirstName == "" {
				customer.FirstName = parts[0]
			}
			if len(parts) >= 2 && customer.LastName == "" {
				customer.LastName = strings.Join(parts[1:], " ")
			}
		}
	}
	if picture, ok := claims["picture"].(string); ok {
		customer.PhotoURL = picture
	}
	
	// Get user record for more details if needed
	userRecord, err := cr.authClient.GetUser(r.Context(), token.UID)
	if err == nil {
		customer.PhoneNumber = userRecord.PhoneNumber
		customer.ProviderID = userRecord.ProviderID
		if customer.DisplayName == "" {
			customer.DisplayName = userRecord.DisplayName
		}
		if customer.PhotoURL == "" {
			customer.PhotoURL = userRecord.PhotoURL
		}
		if customer.Email == "" {
			customer.Email = userRecord.Email
		}
	}

	fmt.Printf("Registering Customer: %+v\n", customer)

	// Forward to DBLayer to create user + profile
	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}

	profilePayload := DBUserProfileRequest{
		Email:                 customer.Email,
		FirstName:             customer.FirstName,
		LastName:              customer.LastName,
		PhoneNumberMobileE164: customer.PhoneNumber,
	}

	body, err := json.Marshal(profilePayload)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to marshal profile payload"), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Forwarding to DBLayer payload=%+v\n", profilePayload)
	proxyReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, dbURL+"/users/profile", bytes.NewBuffer(body))
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to build profile request"), http.StatusInternalServerError)
		return
	}
	proxyReq.Header.Set("Content-Type", "application/json")

	proxyResp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to reach dblayer: %v", err), http.StatusBadGateway)
		return
	}
	defer proxyResp.Body.Close()

	if proxyResp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(proxyResp.Body)
		fmt.Printf("dblayer error status=%d body=%s\n", proxyResp.StatusCode, string(respBody))
		for key, vals := range proxyResp.Header {
			for _, v := range vals {
				w.Header().Add(key, v)
			}
		}
		w.WriteHeader(proxyResp.StatusCode)
		w.Write(respBody)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, customer)
}

func (cr *CustomerRoutes) GetCustomerProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetCustomerProfile started")
	email, ok := authorization.GetEmailFromContext(r.Context())
	if !ok || email == "" {
		email = r.URL.Query().Get("email")
	}

	if email == "" {
		var req struct {
			Email string `json:"email"`
		}
		if err := helpers.ReadJSON(w, r, &req); err == nil {
			email = req.Email
		}
	}
	if email == "" {
		helpers.ErrorJSON(w, fmt.Errorf("email is required"), http.StatusBadRequest)
		return
	}

	// Forward request to DBLayer to fetch customer profile by email
	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}
	fmt.Printf("GetCustomerProfile using email=%s target=%s\n", email, dbURL)

	targetURL := dbURL + "/users/profile?email=" + url.QueryEscape(email)
	proxyReq, err := http.NewRequestWithContext(r.Context(), http.MethodGet, targetURL, nil)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to build request"), http.StatusInternalServerError)
		return
	}

	proxyResp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to reach dblayer: %v", err), http.StatusBadGateway)
		return
	}
	defer proxyResp.Body.Close()

	respBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to read dblayer response"), http.StatusInternalServerError)
		return
	}

	for key, vals := range proxyResp.Header {
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(proxyResp.StatusCode)
	w.Write(respBody)
}

func (cr *CustomerRoutes) UpdateCustomerProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update has been started")
	var payload DBUserProfileRequest

	if err := helpers.ReadJSON(w, r, &payload); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	email, ok := authorization.GetEmailFromContext(r.Context())
	if !ok || email == "" {
		email = r.URL.Query().Get("email")
	}
	if email == "" {
		email = payload.Email
	}
	
	// Ensure payload has the email if we found it elsewhere
	if payload.Email == "" {
		payload.Email = email
	}

	if email == "" {
		helpers.ErrorJSON(w, fmt.Errorf("email is required"), http.StatusBadRequest)
		return
	}

	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}

	body, err := json.Marshal(payload)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to marshal patch payload"), http.StatusInternalServerError)
		return
	}

	urlWithQuery := dbURL + "/users/profile?email=" + url.QueryEscape(email)
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPatch, urlWithQuery, bytes.NewBuffer(body))
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to build request"), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to reach dblayer: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to read dblayer response"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (cr *CustomerRoutes) GetCustomerWishList(w http.ResponseWriter, r *http.Request) {
	userProfileID := chi.URLParam(r, "userProfileID")
	if userProfileID == "" {
		helpers.ErrorJSON(w, fmt.Errorf("userProfileID is required"), http.StatusBadRequest)
		return
	}

	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}

	targetURL := fmt.Sprintf("%s/users/%s/wishlists", dbURL, userProfileID)
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, targetURL, nil)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to build request"), http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to reach dblayer: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to read dblayer response"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func (cr *CustomerRoutes) DeleteCustomerProfile(w http.ResponseWriter, r *http.Request) {
	email, ok := authorization.GetEmailFromContext(r.Context())
	if !ok || email == "" {
		email = r.URL.Query().Get("email")
	}
	if email == "" {
		var req struct {
			Email string `json:"email"`
		}
		if err := helpers.ReadJSON(w, r, &req); err == nil {
			email = req.Email
		}
	}
	if email == "" {
		helpers.ErrorJSON(w, fmt.Errorf("email is required"), http.StatusBadRequest)
		return
	}

	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}

	targetURL := dbURL + "/users/profile?email=" + url.QueryEscape(email)
	req, err := http.NewRequestWithContext(r.Context(), http.MethodDelete, targetURL, nil)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to build request"), http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to reach dblayer: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to read dblayer response"), http.StatusInternalServerError)
		return
	}

	for key, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (cr *CustomerRoutes) AddAddressToProfileEndpoint(w http.ResponseWriter, r *http.Request) {
	userProfileID := chi.URLParam(r, "userProfileID")
	if userProfileID == "" {
		helpers.ErrorJSON(w, fmt.Errorf("userProfileID is required"), http.StatusBadRequest)
		return
	}

	var payload struct {
		Label        string `json:"label"`
		AddressLine1 string `json:"address_line1"`
		AddressLine2 string `json:"address_line2,omitempty"`
		City         string `json:"city"`
		State        string `json:"state"`
		PostalCode   string `json:"postal_code"`
		Country      string `json:"country"`
		IsDefault    bool   `json:"is_default"`
	}
	if err := helpers.ReadJSON(w, r, &payload); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}

	body, err := json.Marshal(payload)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to marshal address payload"), http.StatusInternalServerError)
		return
	}

	targetURL := fmt.Sprintf("%s/users/%s/addresses", dbURL, userProfileID)
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, targetURL, bytes.NewBuffer(body))
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to build request"), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to reach dblayer: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to read dblayer response"), http.StatusInternalServerError)
		return
	}

	for key, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}