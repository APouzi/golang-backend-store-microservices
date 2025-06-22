package retrival

import (
	"net/http"

	"github.com/APouzi/inventory-management/helpers"
	"github.com/go-chi/chi/v5"
)

func AllTransfers(w http.ResponseWriter, r *http.Request) {

	response := map[string]bool{
		"AllTransfers": true,
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, 200, response)
}

func TransfersByID(w http.ResponseWriter, r *http.Request) {

	id_param := chi.URLParam(r, "id")

	response := map[string]string{
		"TransfersByID": id_param,
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, 200, response)
}

func TransfersSearchByProduct(w http.ResponseWriter, r *http.Request) {

	searchTerm := r.URL.Query().Get("search")

	response := map[string]string{
		"TransfersSearchByProduct": searchTerm,
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, 200, response)
}