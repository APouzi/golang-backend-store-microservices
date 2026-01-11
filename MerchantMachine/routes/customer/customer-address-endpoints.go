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