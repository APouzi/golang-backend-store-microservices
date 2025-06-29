package inventory

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/APouzi/DBLayer/helpers"
	"github.com/go-chi/chi/v5"
)


func (routes *InventoryRoutesTray) GetInventoryShelfDetailsByParameter(w http.ResponseWriter, r *http.Request) {
	inventoryID := r.URL.Query().Get("inventory-shelf-id")
	productID := r.URL.Query().Get("product-id")
	shelf := r.URL.Query().Get("shelf")
	var query string
	var queried_var string
	if inventoryID != "" {
		query = "SELECT inventory_shelf_id, inventory_id,quantity_at_shelf, product_id, shelf FROM tblInventoryShelfDetail WHERE inventory_id = ?"
		queried_var = inventoryID
	}else if productID != "" {
		query = "SELECT inventory_shelf_id, inventory_id,quantity_at_shelf, product_id, shelf FROM tblInventoryShelfDetail WHERE product_id = ?"
		queried_var = productID
	}else if shelf != ""{
		shelf += "%"
		query = "SELECT inventory_shelf_id, inventory_id,quantity_at_shelf, product_id, shelf FROM tblInventoryShelfDetail WHERE shelf LIKE ?"
		queried_var = shelf
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

	var shelves []InventoryShelfDetail
	for rows.Next() {
		var shelve InventoryShelfDetail
		err := rows.Scan(&shelve.InventoryShelfID, &shelve.InventoryID,&shelve.QuantityAtShelf,&shelve.ProductID,&shelve.Shelf)
		if err != nil {
			http.Error(w, "Failed to parse result", http.StatusInternalServerError)
			log.Println("Row scan error:", err)
			return
		}
		shelves = append(shelves, shelve)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	helpers.WriteJSON(w,http.StatusAccepted,shelves)
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

	var shelves []InventoryShelfDetail
	for rows.Next() {
		var shelve InventoryShelfDetail
		err := rows.Scan(&shelve.InventoryShelfID, &shelve.InventoryID,&shelve.ProductID,&shelve.QuantityAtShelf,&shelve.Shelf)
		if err != nil {
			http.Error(w, "Failed to parse result", http.StatusInternalServerError)
			log.Println("Row scan error:", err)
			return
		}
		shelves = append(shelves, shelve)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	helpers.WriteJSON(w,http.StatusAccepted,shelves)
}

func (routes *InventoryRoutesTray) GetInventoryShelfDetailByInventoryShelfID(w http.ResponseWriter, r *http.Request) {
	inventory_id := chi.URLParam(r, "inventory-shelf-id")
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

	row := tx.QueryRow("SELECT inventory_shelf_id, inventory_id, quantity_at_shelf, product_id, shelf FROM tblInventoryShelfDetail WHERE inventory_shelf_id = ?",inventory_id)
	if err != nil {
		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}

	
	var shelve InventoryShelfDetail
	err = row.Scan(&shelve.InventoryShelfID, &shelve.InventoryID,&shelve.QuantityAtShelf,&shelve.ProductID,&shelve.Shelf)
	if err != nil {
		http.Error(w, "Failed to parse result", http.StatusInternalServerError)
		log.Println("Row scan error:", err)
		return
	}
	
	if err == sql.ErrNoRows{
		helpers.ErrorJSON(w,errors.New("no rows for record:" + inventory_id),http.StatusInternalServerError)
	}

	if err := tx.Commit(); err != nil {
		helpers.ErrorJSON(w,errors.New("failed to commit transaction"), http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	helpers.WriteJSON(w,http.StatusAccepted,shelve)
}