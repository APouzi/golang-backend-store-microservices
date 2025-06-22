package retrival

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)


func InventoryLocation(w http.ResponseWriter, r *http.Request) {
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


func InventoryLocationAll(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	log.Printf("Received request for InventoryLocationAll ID: %s\n", id)

	response := map[string]string{
		"message": "Received ID",
		"id":      id,
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode([]any{response,response,response,response}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}		
}


func HandleSearch(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("search")

	log.Printf("Search query: %s\n", searchTerm)

	response := map[string]string{
		"search": searchTerm,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}








