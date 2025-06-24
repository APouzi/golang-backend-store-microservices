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


func (routes *InventoryRoutesTray) GetProductInventoriesFromParameter(w http.ResponseWriter, r *http.Request) {
	locationID := r.URL.Query().Get("location-id")
	productID := r.URL.Query().Get("product_id")
	var query string
	var queried_var string
	if locationID != "" {
		query = "SELECT inventory_id, quantity, product_id, location_id, description FROM tblInventoryProductDetail WHERE location_id = ?"
		queried_var = locationID
	}else if productID != "" {
		query = "SELECT inventory_id, quantity, product_id, location_id, description FROM tblInventoryProductDetail WHERE product_id = ?"
		queried_var = productID
	}else{
		http.Error(w,"No Parameter Found",http.StatusBadRequest)
	}
	
	tx, err := routes.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start DB transaction", http.StatusInternalServerError)
		log.Println("Begin transaction error:", err)
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(query, queried_var)
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

func (routes *InventoryRoutesTray) GetAllInventoryProductDetails(w http.ResponseWriter, r *http.Request) {
	
	tx, err := routes.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start DB transaction", http.StatusInternalServerError)
		log.Println("Begin transaction error:", err)
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT inventory_id, quantity_at_location, product_id, location_id, description FROM tblInventoryProductDetail WHERE ")
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


func (routes *InventoryRoutesTray) GetAllInventoryProductDetailsByProduct(w http.ResponseWriter, r *http.Request) {
	inventory_id := chi.URLParam(r, "inventory-id")
	if inventory_id == "" {
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

	rows, err := tx.Query("SELECT inventory_id, quantity, product_id, location_id, description FROM tblInventoryProductDetail WHERE inventory_id = ?", inventory_id)
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



func (routes *InventoryRoutesTray) GetAllInventoryShelfDetail(w http.ResponseWriter, r *http.Request) {
	
	tx, err := routes.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start DB transaction", http.StatusInternalServerError)
		log.Println("Begin transaction error:", err)
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT inventory_shelf_id, quantity_at_shelf, product_id, inventory_id, shelf FROM tblInventoryShelfDetail")
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


func (routes *InventoryRoutesTray) GetInventoryShelfDetailByInventoryID(w http.ResponseWriter, r *http.Request) {
	inventory_id := chi.URLParam(r, "product-variation-id")
	if inventory_id == "" {
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

	rows, err := tx.Query("SELECT inventory_shelf_id, quantity_at_shelf, product_id, inventory_id, shelf FROM tblInventoryShelfDetail")
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

func (routes *InventoryRoutesTray) GetInventoryShelfDetailByParameter(w http.ResponseWriter, r *http.Request) {
	inventoryID := r.URL.Query().Get("inventory_id")
	productID := r.URL.Query().Get("product_id")
	shelf := r.URL.Query().Get("shelf")
	var query string
	var queried_var string
	if inventoryID != "" {
		query = "SELECT inventory_id, quantity, product_id, location_id, description FROM tblInventoryProductDetail WHERE location_id = ?"
		queried_var = inventoryID
	}else if productID != "" {
		query = "SELECT inventory_id, quantity, product_id, location_id, description FROM tblInventoryProductDetail WHERE product_id = ?"
		queried_var = productID
	}else if shelf != ""{
		query = "SELECT inventory_id, quantity, product_id, location_id, description FROM tblInventoryProductDetail WHERE product_id = ?"
		queried_var = productID
	}else{
		http.Error(w,"No Parameter Found",http.StatusBadRequest)
	}
	
	tx, err := routes.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start DB transaction", http.StatusInternalServerError)
		log.Println("Begin transaction error:", err)
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(query, queried_var)
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