package inventory

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	rows, err := tx.Query("SELECT inventory_id, quantity_at_location, product_size_id, location_id, description FROM tblInventoryProductDetail")
	if err != nil {
		http.Error(w, "Failed to fetch Inventory Product Detail", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var inventory_product_detail []InventoryProductDetail
	for rows.Next() {
		var invprddetail InventoryProductDetail
		err := rows.Scan(&invprddetail.InventoryID, &invprddetail.Quantity, &invprddetail.SizeID, &invprddetail.LocationID, &invprddetail.Description)
		if err != nil {
			http.Error(w, "Failed to parse result", http.StatusInternalServerError)
			log.Println("Row scan error:", err)
			return
		}
		inventory_product_detail = append(inventory_product_detail, invprddetail)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(inventory_product_detail); err != nil {
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

	var inventory_product_details []InventoryProductDetail
	for rows.Next() {
		var invprddet InventoryProductDetail
		err := rows.Scan(&invprddet.InventoryID, &invprddet.Quantity, &invprddet.SizeID, &invprddet.LocationID, &invprddet.Description, &invprddet.Variation_Name)
		if err != nil {
			helpers.ErrorJSON(w,errors.New("failed to parse result"), http.StatusInternalServerError)
			log.Println("row scan error:", err)
			return
		}
		inventory_product_details = append(inventory_product_details, invprddet)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	if len(inventory_product_details) == 0{
		helpers.WriteJSON(w,http.StatusAccepted, []InventoryProductDetail{})	
		return
	}

	helpers.WriteJSON(w,http.StatusAccepted, inventory_product_details)

}

func (routes *InventoryRoutesTray) GetInventoryProductDetailFromParameter(w http.ResponseWriter, r *http.Request) {
	locationID := r.URL.Query().Get("location_id")
	productID := r.URL.Query().Get("product_size_id")
	var query string
	var queried_var string
	if locationID != "" {
		query = "SELECT inventory_id, quantity_at_location, product_size_id, location_id, description FROM tblInventoryProductDetail WHERE location_id = ?"
		queried_var = locationID
	}else if productID != "" {
		query = "SELECT inventory_id, quantity_at_location, product_size_id, location_id, description FROM tblInventoryProductDetail WHERE product_size_id = ?"
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

	var inventory_product_details []InventoryProductDetail
	for rows.Next() {
		var invprddet InventoryProductDetail
		err := rows.Scan(&invprddet.InventoryID, &invprddet.Quantity,&invprddet.SizeID,&invprddet.LocationID,&invprddet.Description)
		if err != nil {
			http.Error(w, "Failed to parse result", http.StatusInternalServerError)
			log.Println("Row scan error:", err)
			return
		}
		inventory_product_details = append(inventory_product_details, invprddet)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	helpers.WriteJSON(w,http.StatusAccepted,inventory_product_details)
}




//   ___                     _       
//  |_ _|_ __  ___  ___ _ __| |_ ___ 
//   | || '_ \/ __|/ _ \ '__| __/ __|
//   | || | | \__ \  __/ |  | |_\__ \
//  |___|_| |_|___/\___|_|   \__|___/    
                                                                          
func (adminProdRoutes *InventoryRoutesTray) CreateInventoryProductDetail(w http.ResponseWriter, r *http.Request) {

	ProductID := chi.URLParam(r, "ProductID")
	fmt.Println("productID", ProductID)
	inventory_product_detail := InventoryProductDetail{}
	helpers.ReadJSON(w,r, &inventory_product_detail)
	varitCrt := Confirmation{}
		
	varit, err := adminProdRoutes.DB.Exec("INSERT INTO tblInventoryProductDetail(product_size_id, quantity_at_location, location_id, description) VALUES(?,?,?,?)", inventory_product_detail.SizeID,inventory_product_detail.Quantity, inventory_product_detail.Description, inventory_product_detail.Description)
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


func (adminProdRoutes *InventoryRoutesTray) CreateInventoryLocationTransfer(w http.ResponseWriter, r *http.Request) {

	invShelfDtl := InventoryLocationTransfer{}
	helpers.ReadJSON(w,r, &invShelfDtl)
	createConfirm := Confirmation{}
		
	varit, err := adminProdRoutes.DB.Exec("INSERT INTO tblInventoryLocationTransfers(transfers_id, source_location_id, destination_location_id, product_id, quantity, transfer_date, description, status) VALUES(?,?,?,?,?,?,?,?)", invShelfDtl.TransfersID,invShelfDtl.SourceLocationID, invShelfDtl.DestinationLocationID, invShelfDtl.ProductID, invShelfDtl.Quantity, invShelfDtl.TransferDate, invShelfDtl.Description, invShelfDtl.Status)
	if err != nil{
		log.Println("insert into tblInventoryLocationTransfers failed")
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblInventoryLocationTransfers failed"),400)
		return
	}
	createConfirm.Createdid, err = varit.LastInsertId()
	if err != nil{
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblInventoryLocationTransfers failed, could not retrieve varitation id"),400)
		return
	}	
	helpers.WriteJSON(w, http.StatusCreated,createConfirm)
}


//  ___          _           _          
// (  _`\       (_ )        ( )_        
// | | ) |   __  | |    __  | ,_)   __  
// | | | ) /'__`\| |  /'__`\| |   /'__`\
// | |_) |(  ___/| | (  ___/| |_ (  ___/
// (____/'`\____|___)`\____)`\__)`\____)

func (adminProdRoutes *InventoryRoutesTray) DeleteInventoryLocationTransfer(w http.ResponseWriter, r *http.Request) {

	invShelfDtl := InventoryLocationTransfer{}
	helpers.ReadJSON(w,r, &invShelfDtl)
	createConfirm := Confirmation{}
		
	varit, err := adminProdRoutes.DB.Exec("DELETE FROM tblInventoryLocationTransfers WHERE transfers_id = ?", invShelfDtl.TransfersID)
	if err != nil{
		log.Println("insert into tblInventoryLocationTransfers failed")
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblInventoryLocationTransfers failed"),400)
		return
	}
	createConfirm.Createdid, err = varit.LastInsertId()
	if err != nil{
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblInventoryLocationTransfers failed, could not retrieve varitation id"),400)
		return
	}	
	helpers.WriteJSON(w, http.StatusCreated,createConfirm)
}


func (adminProdRoutes *InventoryRoutesTray) DeleteInventoryShelfDetail(w http.ResponseWriter, r *http.Request) {

	invShelfDtl := InventoryShelfDetail{}
	helpers.ReadJSON(w,r, &invShelfDtl)
	createConfirm := Confirmation{}
		
	varit, err := adminProdRoutes.DB.Exec("DELETE FROM tblInventoryShelfDetail WHERE inventory_shelf_id = ?", invShelfDtl.InventoryShelfID)
	if err != nil{
		log.Println("insert into tblInventoryLocationTransfers failed")
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

func (adminProdRoutes *InventoryRoutesTray) DeleteInventoryProductDetail(w http.ResponseWriter, r *http.Request) {

	invShelfDtl := InventoryShelfDetail{}
	helpers.ReadJSON(w,r, &invShelfDtl)
	createConfirm := Confirmation{}
		
	varit, err := adminProdRoutes.DB.Exec("DELETE FROM tblInventoryProductDetail WHERE inventory_id = ?", invShelfDtl.InventoryShelfID)
	if err != nil{
		log.Println("insert into tblInventoryProductDetail failed")
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblInventoryProductDetail failed"),400)
		return
	}
	createConfirm.Createdid, err = varit.LastInsertId()
	if err != nil{
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblInventoryProductDetail failed, could not retrieve varitation id"),400)
		return
	}	
	helpers.WriteJSON(w, http.StatusCreated,createConfirm)
}


func (adminProdRoutes *InventoryRoutesTray) DeleteLocation(w http.ResponseWriter, r *http.Request) {

	invShelfDtl := Location{}
	helpers.ReadJSON(w,r, &invShelfDtl)
	createConfirm := Confirmation{}
		
	varit, err := adminProdRoutes.DB.Exec("DELETE FROM tblLocation WHERE location_id = ?", invShelfDtl.LocationID)
	if err != nil{
		log.Println("insert into tblLocation failed")
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblLocation failed"),400)
		return
	}
	createConfirm.Createdid, err = varit.LastInsertId()
	if err != nil{
		log.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblInventoryProductDetail failed, could not retrieve varitation id"),400)
		return
	}	
	helpers.WriteJSON(w, http.StatusCreated,createConfirm)
}


//  __  __              __           __             
// /\ \/\ \            /\ \         /\ \__          
// \ \ \ \ \  _____    \_\ \     __ \ \ ,_\    __   
//  \ \ \ \ \/\ '__`\  /'_` \  /'__`\\ \ \/  /'__`\ 
//   \ \ \_\ \ \ \L\ \/\ \L\ \/\ \L\.\\ \ \_/\  __/ 
//    \ \_____\ \ ,__/\ \___,_\ \__/.\_\ \__\ \____\
//     \/_____/\ \ \/  \/__,_ /\/__/\/_/\/__/\/____/
//              \ \_\                               
//               \/_/

func (adminProdRoutes *InventoryRoutesTray) UpdateInventoryShelfDetail(w http.ResponseWriter, r *http.Request) {
    // Use a consistent path param name, e.g. "inventory-shelf-id"
    inventoryShelfID := chi.URLParam(r, "inventory-id")
    if inventoryShelfID == "" {
        helpers.ErrorJSON(w, errors.New("missing inventory-id parameter"), http.StatusBadRequest)
        return
    }

    var payload InventoryShelfDetail

    // Decode into a POINTER and return on error
    if err := helpers.ReadJSON(w, r, &payload); err != nil {
        helpers.ErrorJSON(w, errors.New("invalid JSON payload"), http.StatusBadRequest)
        return
    }

    // Build SET dynamically; presence in JSON (non-nil) => update
	setClauses := []string{}
	args := []any{}

    if payload.QuantityAtShelf != nil {
        setClauses = append(setClauses, "quantity_at_shelf = ?")
        args = append(args, *payload.QuantityAtShelf)
    }
    if payload.ProductID != nil {
        setClauses = append(setClauses, "product_id = ?")
        args = append(args, *payload.ProductID)
    }
    if payload.Shelf != nil {
        setClauses = append(setClauses, "shelf = ?")
        args = append(args, *payload.Shelf)
    }

    if len(setClauses) == 0 {
        helpers.ErrorJSON(w, errors.New("no fields provided to update"), http.StatusBadRequest)
        return
    }

    idInt, err := strconv.ParseInt(inventoryShelfID, 10, 64)
    if err != nil {
        helpers.ErrorJSON(w, errors.New("invalid inventory-shelf-id format"), http.StatusBadRequest)
        return
    }

    // WHERE param last
    args = append(args, idInt)

    query := fmt.Sprintf(
        "UPDATE tblInventoryShelfDetail SET %s WHERE inventory_shelf_id = ?",
        strings.Join(setClauses, ", "),
    )
    fmt.Println("Query to be executed:", query, "args:", args)

    result, err := adminProdRoutes.DB.Exec(query, args...)
    if err != nil {
        log.Println("update failed:", err)
        helpers.ErrorJSON(w, errors.New("update failed"), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Println("RowsAffected failed:", err)
        helpers.ErrorJSON(w, errors.New("could not determine update result"), http.StatusInternalServerError)
        return
    }
    if rowsAffected == 0 {
        helpers.ErrorJSON(w, errors.New("no record found to update"), http.StatusNotFound)
        return
    }

    helpers.WriteJSON(w, http.StatusOK, Confirmation{Createdid: idInt})
}


func (adminProdRoutes *InventoryRoutesTray) UpdateInventoryProductDetail(w http.ResponseWriter, r *http.Request) {
	inventoryID := chi.URLParam(r, "inventory-id")
	if inventoryID == "" {
		helpers.ErrorJSON(w, errors.New("missing inventory-id parameter"), http.StatusBadRequest)
		return
	}
	fmt.Println("Updating Inventory Product Detail for ID:", inventoryID)

	var payload InventoryProductDetail

	if err := helpers.ReadJSON(w, r, &payload); err != nil {
		helpers.ErrorJSON(w, errors.New("invalid JSON payload"), http.StatusBadRequest)
		return
	}

	fmt.Println("Payload to update:", payload)

	setClauses := []string{}
	args := []interface{}{}

	if payload.Quantity != nil {
		setClauses = append(setClauses, "quantity_at_location = ?")
		args = append(args, *payload.Quantity)
	}
	if payload.SizeID != nil {
		setClauses = append(setClauses, "product_size_id = ?")
		args = append(args, *payload.SizeID)
	}
	if payload.LocationID != nil {
		setClauses = append(setClauses, "location_id = ?")
		args = append(args, *payload.LocationID)
	}
	if payload.Description != nil {
		setClauses = append(setClauses, "description = ?")
		args = append(args, *payload.Description)
	}

	if len(setClauses) == 0 {
		helpers.ErrorJSON(w, errors.New("no fields provided to update"), http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseInt(inventoryID, 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("invalid inventory-id format"), http.StatusBadRequest)
		return
	}

	args = append(args, idInt)

	query := fmt.Sprintf(
		"UPDATE tblInventoryProductDetail SET %s WHERE inventory_id = ?",
		strings.Join(setClauses, ", "),
	)
	fmt.Println("Query to be executed:", query, "args:", args)

	result, err := adminProdRoutes.DB.Exec(query, args...)
	if err != nil {
		log.Println("update failed:", err)
		helpers.ErrorJSON(w, errors.New("update failed"), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("RowsAffected failed:", err)
		helpers.ErrorJSON(w, errors.New("could not determine update result"), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		helpers.ErrorJSON(w, errors.New("no record found to update"), http.StatusNotFound)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, Confirmation{Createdid: idInt})
}



