package insertion

import (
	"net/http"

	"github.com/APouzi/inventory-management/helpers"
	"github.com/go-chi/chi/v5"
)

func InsertProductDetail(w http.ResponseWriter, r *http.Request) {

	location_param := chi.URLParam(r, "location")

	response := map[string]string{
		"allProducts": location_param,
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, 200, response)
}
