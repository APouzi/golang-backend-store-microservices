package customer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/APouzi/MerchantMachinee/routes/helpers"
	"github.com/go-chi/chi/v5"
)



func (cr *CustomerRoutes) GetCustomerWishList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetCustomerWishList: started")
	userProfileID := chi.URLParam(r, "userProfileID")
	wishlistID := chi.URLParam(r, "wishlistID")
	if wishlistID == "" {
		fmt.Println("GetCustomerWishList: missing wishlistID")
		helpers.ErrorJSON(w, fmt.Errorf("wishlistID is required"), http.StatusBadRequest)
		return
	}
	if userProfileID == "" {
		fmt.Println("GetCustomerWishList: missing userProfileID")
		helpers.ErrorJSON(w, fmt.Errorf("userProfileID is required"), http.StatusBadRequest)
		return
	}

	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}
	fmt.Println("GetCustomerWishList: dbURL =", dbURL)

	targetURL := fmt.Sprintf("%s/users/%s/wishlists/%s", dbURL, userProfileID, wishlistID)
	fmt.Println("GetCustomerWishList: targetURL =", targetURL)
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, targetURL, nil)
	if err != nil {
		fmt.Println("GetCustomerWishList: failed to build request:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to build request"), http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("GetCustomerWishList: failed to reach dblayer:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to reach dblayer: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	fmt.Println("GetCustomerWishList: received response from dblayer")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("GetCustomerWishList: failed to read dblayer response:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to read dblayer response"), http.StatusInternalServerError)
		return
	}
	fmt.Println("GetCustomerWishList: successfully read response body")

	var decoded interface{}
	if err := json.Unmarshal(body, &decoded); err != nil {
		fmt.Println("GetCustomerWishList: failed to decode dblayer response:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to decode dblayer response"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusAccepted, decoded)
	fmt.Println("GetCustomerWishList: response sent")
}

func (cr *CustomerRoutes) CreateCustomerWishList(w http.ResponseWriter, r *http.Request) {
	userProfileID := chi.URLParam(r, "userProfileID")
	if userProfileID == "" {
		helpers.ErrorJSON(w, fmt.Errorf("userProfileID is required"), http.StatusBadRequest)
		return
	}

	// Parse wishlist creation payload
	var payload struct {
		WishlistName string `json:"wishlist_name"`
		IsDefault    *bool  `json:"is_default,omitempty"`
	}
	falseValue := false
	payload.IsDefault = &falseValue
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid request payload"), http.StatusBadRequest)
		return
	}

	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}

	targetURL := fmt.Sprintf("%s/users/%s/wishlists", dbURL, userProfileID)
	reqBody, err := json.Marshal(payload)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to marshal request body"), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, targetURL, bytes.NewBuffer(reqBody))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to read dblayer response"), http.StatusInternalServerError)
		return
	}

	var decoded interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &decoded); err != nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				helpers.WriteJSON(w, resp.StatusCode, string(body))
				return
			}
			helpers.ErrorJSON(w, fmt.Errorf("dblayer error: %s", string(body)), resp.StatusCode)
			return
		}
	} else {
		decoded = map[string]interface{}{}
	}

	helpers.WriteJSON(w, resp.StatusCode, decoded)
}

func (cr *CustomerRoutes) GetAllCustomerWishLists(w http.ResponseWriter, r *http.Request) {
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

	var decoded interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &decoded); err != nil {
			// If dblayer returned non-JSON but a valid HTTP status, forward raw body for visibility.
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				helpers.WriteJSON(w, resp.StatusCode, string(body))
				return
			}
			helpers.ErrorJSON(w, fmt.Errorf("dblayer error: %s", string(body)), resp.StatusCode)
			return
		}
	} else {
		decoded = map[string]interface{}{}
	}

	helpers.WriteJSON(w, resp.StatusCode, decoded)
}

func (cr *CustomerRoutes) AddProductToWishListEndpoint(w http.ResponseWriter, r *http.Request) {
	wishListID := chi.URLParam(r, "wishlistID")

	// Parse size ID from request body
	var payload struct {
		SizeID int `json:"size_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid request payload"), http.StatusBadRequest)
		return
	}
	
	// Forward request to DBLayer
	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}
	
	targetURL := fmt.Sprintf("%s/wishlists/%s/products", dbURL, wishListID)
	reqBody, err := json.Marshal(payload)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to marshal request body"), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, targetURL, bytes.NewBuffer(reqBody))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to read dblayer response"), http.StatusInternalServerError)
		return
	}

	var decoded interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &decoded); err != nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				helpers.WriteJSON(w, resp.StatusCode, string(body))
				return
			}
			helpers.ErrorJSON(w, fmt.Errorf("dblayer error: %s", string(body)), resp.StatusCode)
			return
		}
	} else {
		decoded = map[string]interface{}{}
	}

	helpers.WriteJSON(w, resp.StatusCode, decoded)
}

func (cr *CustomerRoutes) AddProductToDefaultWishListEndpoint(w http.ResponseWriter, r *http.Request) {
	userProfileID := chi.URLParam(r, "userProfileID")

	// Parse size ID from request body
	var payload struct {
		SizeID int `json:"size_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid request payload"), http.StatusBadRequest)
		return
	}
	
	// Forward request to DBLayer
	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}
	
	targetURL := fmt.Sprintf("%s/users/%s/wishlists/default/products", dbURL, userProfileID)
	reqBody, err := json.Marshal(payload)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to marshal request body"), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, targetURL, bytes.NewBuffer(reqBody))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to read dblayer response"), http.StatusInternalServerError)
		return
	}

	var decoded interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &decoded); err != nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				helpers.WriteJSON(w, resp.StatusCode, string(body))
				return
			}
			helpers.ErrorJSON(w, fmt.Errorf("dblayer error: %s", string(body)), resp.StatusCode)
			return
		}
	} else {
		decoded = map[string]interface{}{}
	}

	helpers.WriteJSON(w, resp.StatusCode, decoded)
}

func (cr *CustomerRoutes) RemoveProductFromWishListEndpoint(w http.ResponseWriter, r *http.Request) {
	wishListID := chi.URLParam(r, "wishlistID")

	// Parse size ID from request body
	var payload struct {
		SizeID int `json:"size_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid request payload"), http.StatusBadRequest)
		return
	}
	
	// Forward request to DBLayer
	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}
	
	targetURL := fmt.Sprintf("%s/wishlists/%s/products", dbURL, wishListID)
	reqBody, err := json.Marshal(payload)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to marshal request body"), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodDelete, targetURL, bytes.NewBuffer(reqBody))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to read dblayer response"), http.StatusInternalServerError)
		return
	}

	var decoded interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &decoded); err != nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				helpers.WriteJSON(w, resp.StatusCode, string(body))
				return
			}
			helpers.ErrorJSON(w, fmt.Errorf("dblayer error: %s", string(body)), resp.StatusCode)
			return
		}
	} else {
		decoded = map[string]interface{}{}
	}

	helpers.WriteJSON(w, resp.StatusCode, decoded)
}