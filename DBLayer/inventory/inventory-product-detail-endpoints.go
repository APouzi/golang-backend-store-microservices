package inventory

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/APouzi/DBLayer/helpers"
	"github.com/go-chi/chi/v5"
)

//   ____        _ _
//  |  _ \ _   _| | |___
//  | |_) | | | | | / __|
//  |  __/| |_| | | \__ \
//  |_|    \__,_|_|_|___/

func (routes *InventoryRoutesTray) GetAllInventoryProductDetails(w http.ResponseWriter, r *http.Request) {

	tx, err := routes.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start DB transaction", http.StatusInternalServerError)
		log.Println("Begin transaction error:", err)
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT inventory_id, quantity_at_location, product_id, location_id, description FROM tblInventoryProductDetail")
	if err != nil {
		http.Error(w, "Failed to fetch Inventory Product Detail", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var locations []InventoryProductDetail
	for rows.Next() {
		var loc InventoryProductDetail
		err := rows.Scan(&loc.InventoryID, &loc.Quantity, &loc.ProductID, &loc.LocationID, &loc.Description)
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

func (routes *InventoryRoutesTray) GetAllInventoryProductDetailsByID(w http.ResponseWriter, r *http.Request) {
	inventory_id := chi.URLParam(r, "inventory-id")
	if inventory_id == "" {
		http.Error(w, "Missing inventory-id parameter", http.StatusBadRequest)
		return
	}

	tx, err := routes.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start DB transaction", http.StatusInternalServerError)
		log.Println("Begin transaction error:", err)
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT ipd.inventory_id, ipd.quantity_at_location, ipd.product_id, ipd.location_id, ipd.description, pv.Variation_Name FROM tblInventoryProductDetail ipd JOIN tblProductVariation pv ON pv.Variation_ID = ipd.product_id WHERE inventory_id = ? ", inventory_id)
	if err != nil {
		helpers.ErrorJSON(w,errors.New("failed to fetch locations:" + err.Error()), http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var locations []InventoryProductDetail
	for rows.Next() {
		var loc InventoryProductDetail
		err := rows.Scan(&loc.InventoryID, &loc.Quantity, &loc.ProductID, &loc.LocationID, &loc.Description, &loc.Variation_Name)
		if err != nil {
			helpers.ErrorJSON(w,errors.New("failed to parse result"), http.StatusInternalServerError)
			log.Println("row scan error:", err)
			return
		}
		locations = append(locations, loc)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	if len(locations) == 0{
		helpers.WriteJSON(w,http.StatusAccepted, []InventoryProductDetail{})	
		return
	}

	helpers.WriteJSON(w,http.StatusAccepted, locations)

}

func (routes *InventoryRoutesTray) GetInventoryProductDetailFromParameter(w http.ResponseWriter, r *http.Request) {
	locationID := r.URL.Query().Get("location_id")
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

	helpers.WriteJSON(w,http.StatusAccepted,locations)
}




//   ___                     _       
//  |_ _|_ __  ___  ___ _ __| |_ ___ 
//   | || '_ \/ __|/ _ \ '__| __/ __|
//   | || | | \__ \  __/ |  | |_\__ \
//  |___|_| |_|___/\___|_|   \__|___/    
                                                                          
func (adminProdRoutes *InventoryRoutesTray) CreateInventoryProductDetail(w http.ResponseWriter, r *http.Request) {

	ProductID := chi.URLParam(r, "ProductID")
	fmt.Println("productID", ProductID)
	variation := InventoryProductDetail{}
	helpers.ReadJSON(w,r, &variation)
	varitCrt := Confirmation{}
		
	varit, err := adminProdRoutes.DB.Exec("INSERT INTO tblInventoryProductDetail(product_id, quantity_at_location, location_id, description) VALUES(?,?,?,?)", variation.ProductID,variation.Quantity, variation.Description, variation.Description)
	if err != nil{
		log.Println("insert into tblProductVariation failed")
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed"),400)
		return
	}
	varitCrt.Createdid, err = varit.LastInsertId()
	if err != nil{
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed, could not retrieve varitation id"),400)
		return
	}	
	helpers.WriteJSON(w, http.StatusCreated,varitCrt)
}


func (adminProdRoutes *InventoryRoutesTray) CreateLocation(w http.ResponseWriter, r *http.Request) {

	location := Location{}
	helpers.ReadJSON(w,r, &location)
	createConfirm := Confirmation{}
		
	varit, err := adminProdRoutes.DB.Exec("INSERT INTO Location(product_id, LocationID, Name, Description, Latitude, Longitude) VALUES(?,?,?,?)", location.LocationID,location.Name, location.Description, location.Latitude, )
	if err != nil{
		log.Println("insert into tblProductVariation failed")
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed"),400)
		return
	}
	createConfirm.Createdid, err = varit.LastInsertId()
	if err != nil{
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed, could not retrieve varitation id"),400)
		return
	}	
	helpers.WriteJSON(w, http.StatusCreated,createConfirm)
}

func (adminProdRoutes *InventoryRoutesTray) CreateInventoryShelfDetail(w http.ResponseWriter, r *http.Request) {

	invShelfDtl := InventoryShelfDetail{}
	helpers.ReadJSON(w,r, &invShelfDtl)
	createConfirm := Confirmation{}
		
	varit, err := adminProdRoutes.DB.Exec("INSERT INTO tblInventoryShelfDetail(inventory_shelf_id, inventory_id, quantity_at_shelf, product_id, shelf) VALUES(?,?,?,?,?)", invShelfDtl.InventoryShelfID,invShelfDtl.InventoryID, invShelfDtl.QuantityAtShelf, invShelfDtl.ProductID, invShelfDtl.Shelf)
	if err != nil{
		log.Println("insert into tblInventoryShelfDetail failed")
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblInventoryShelfDetail failed"),400)
		return
	}
	createConfirm.Createdid, err = varit.LastInsertId()
	if err != nil{
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblInventoryShelfDetail failed, could not retrieve varitation id"),400)
		return
	}	
	helpers.WriteJSON(w, http.StatusCreated,createConfirm)
}



