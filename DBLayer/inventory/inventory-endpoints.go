package inventory

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)


func GetInventoryRoutesTrayInstance (databaseInstance *sql.DB) *InventoryRoutesTray{
	return &InventoryRoutesTray{DB: databaseInstance}
}

type InventoryRoutesTray struct{
	DB *sql.DB
}



func (routes *InventoryRoutesTray) GetAllLocations(w http.ResponseWriter, r *http.Request) {
	tx, err := routes.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start DB transaction", http.StatusInternalServerError)
		log.Println("Begin transaction error:", err)
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT location_id, name, description, latitude, longitude, street_address FROM tblLocation")
	if err != nil {
		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var loc Location
		err := rows.Scan(&loc.LocationID, &loc.Name,&loc.Description,&loc.Latitude,&loc.Longitude,&loc.StreetAddress,)
		if err != nil {
			http.Error(w, "Failed to parse result", http.StatusInternalServerError)
			log.Println("Row scan error:", err)
			return
		}
		locations = append(locations, loc)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locations); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		log.Println("JSON encode error:", err)
		return
	}
}


func (routes *InventoryRoutesTray) GetAllProductInventoriesFromLocation(w http.ResponseWriter, r *http.Request) {
	locationID := chi.URLParam(r, "product-variation-id")
	if locationID == "" {
		http.Error(w, "Missing product-variation-id parameter", http.StatusBadRequest)
		return
	}
	
	tx, err := routes.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start DB transaction", http.StatusInternalServerError)
		log.Println("Begin transaction error:", err)
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT inventory_id, quantity, product_id, location_id, description FROM tblInventoryProductDetail WHERE location_id = ?", locationID)
	if err != nil {
		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var locations []InventoryProductDetail
	for rows.Next() {
		var loc InventoryProductDetail
		err := rows.Scan(&loc.InventoryID, &loc.Quantity,&loc.ProductID,&loc.LocationID,&loc.Description)
		if err != nil {
			http.Error(w, "Failed to parse result", http.StatusInternalServerError)
			log.Println("Row scan error:", err)
			return
		}
		locations = append(locations, loc)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locations); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		log.Println("JSON encode error:", err)
		return
	}
}