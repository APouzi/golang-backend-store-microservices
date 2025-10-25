package updates

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func UpdateLocationProductIntermediaryTable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	log.Printf("Received request for ID: %s\n", id)

	response := map[string]string{
		"message": "Received ID",
		"id":      id,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}		
}


func UpdateLocationProductDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	log.Printf("Received request for ID: %s\n", id)

	response := map[string]string{
		"message": "Received ID",
		"id":      id,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}		
}