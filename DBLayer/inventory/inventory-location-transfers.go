package inventory

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/APouzi/DBLayer/helpers"
	"github.com/go-chi/chi/v5"
)


func (routes *InventoryRoutesTray) GetAllLocationTransfers(w http.ResponseWriter, r *http.Request) {

	tx, err := routes.DB.Begin()
	if err != nil {
		helpers.ErrorJSON(w,errors.New("failed to start DB transaction"), http.StatusInternalServerError)
		log.Println("Could not start db transation:", err)
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT transfers_id, quantity, source_location_id, destination_location_id, product_id, transfer_date, description, status FROM tblInventoryLocationTransfers")
	if err != nil {
		helpers.ErrorJSON(w,errors.New("failed to fetch DB transaction"), http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var transfers []InventoryLocationTransfer
	for rows.Next() {
		var transfer InventoryLocationTransfer
		var description sql.NullString
		var status sql.NullString
		err := rows.Scan(
			&transfer.TransfersID,
			&transfer.Quantity,
			&transfer.ProductID,
			&transfer.SourceLocationID,
			&transfer.DestinationLocationID,
			&transfer.TransferDate,
			&description,
			&status,
		)
		if err != nil {
			helpers.ErrorJSON(w,errors.New("failed to parse result"), http.StatusInternalServerError)
			log.Println("Row scan error:", err)
			return
		}

		if description.Valid {
			transfer.Description = &status.String
		}
		if status.Valid {
			transfer.Status = &status.String
		}

		transfers = append(transfers, transfer)
	}

	if err := tx.Commit(); err != nil {
		helpers.ErrorJSON(w,errors.New("failed to commit transaction"), http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	helpers.WriteJSON(w,200,transfers)
}




func (routes *InventoryRoutesTray) GetInventoryLocationTransfersById(w http.ResponseWriter, r *http.Request) {
	transfers_id := chi.URLParam(r, "transfers-id")
	if transfers_id == "" {
		helpers.ErrorJSON(w,errors.New("missing transfers_id parameter"), http.StatusInternalServerError)
		return
	}
	
	tx, err := routes.DB.Begin()
	if err != nil {
		helpers.ErrorJSON(w,errors.New("failed to start DB transaction"), http.StatusInternalServerError)
		log.Println("Begin transaction error:", err)
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT transfers_id, quantity, product_id, source_location_id, destination_location_id, transfer_date, description, status FROM tblInventoryLocationTransfers WHERE transfers_id = ?", transfers_id)
	
	

	var transfers []InventoryLocationTransfer
	
	var transfer InventoryLocationTransfer
	var description sql.NullString
	var status sql.NullString

	err = row.Scan(
		&transfer.TransfersID,
		&transfer.Quantity,
		&transfer.ProductID,
		&transfer.SourceLocationID,
		&transfer.DestinationLocationID,
		&transfer.TransferDate,
		&description,
		&status,
	)
	if err == sql.ErrNoRows{
		helpers.ErrorJSON(w,errors.New("No records for Transfers with id of " + transfers_id), http.StatusInternalServerError)
		return
	}
	if err != nil {
		helpers.ErrorJSON(w,errors.New("failed to parse result"), http.StatusInternalServerError)
		log.Println("Row scan error:", err)
		return
	}

	if description.Valid {
		transfer.Description = &description.String
	}
	if status.Valid {
		transfer.Status = &status.String
	}

	transfers = append(transfers, transfer)
	

	if err := tx.Commit(); err != nil {
		helpers.ErrorJSON(w,errors.New("failed to commit transaction"), http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	helpers.WriteJSON(w,200,transfers)

}



func (routes *InventoryRoutesTray) GetLocationTransfersByParam(w http.ResponseWriter, r *http.Request) {

	sourceLocationID := r.URL.Query().Get("source_location_id")
	destinationLocationID := r.URL.Query().Get("destination_location_id")
	transferDate := r.URL.Query().Get("transfer_date")
	status := r.URL.Query().Get("status")
	productID := r.URL.Query().Get("product_id")
	var query, queried_var string


	if sourceLocationID != "" {
		query = "SELECT transfers_id, quantity, product_id, source_location_id, destination_location_id, transfer_date, description, status FROM tblInventoryLocationTransfers WHERE source_location_id = ?"
		queried_var = sourceLocationID
	}else if productID != "" {
		query = "SELECT transfers_id, quantity, product_id, source_location_id, destination_location_id,  transfer_date, description, status FROM tblInventoryLocationTransfers WHERE product_id = ?"
		queried_var = productID
	}else if destinationLocationID != ""{
		query = "SELECT transfers_id, quantity, product_id, source_location_id, destination_location_id, transfer_date, description, status FROM tblInventoryLocationTransfers WHERE destination_location_id = ?"
		queried_var = destinationLocationID
	}else if transferDate != ""{
		query = "SELECT transfers_id, quantity, product_id, source_location_id, destination_location_id, transfer_date, description, status FROM tblInventoryLocationTransfers WHERE transfer_date = ?"
		queried_var = transferDate
	}else if status != ""{
		query = "SELECT transfers_id, quantity, product_id, source_location_id, destination_location_id, transfer_date, description, status FROM tblInventoryLocationTransfers WHERE status = ?"
		queried_var = status
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

	var transfers []InventoryLocationTransfer
	for rows.Next() {
		var transfer InventoryLocationTransfer
		// var description sql.NullString
		// var status sql.NullString

		err := rows.Scan(
			&transfer.TransfersID,
			&transfer.Quantity,
			&transfer.ProductID,
			&transfer.SourceLocationID,
			&transfer.DestinationLocationID,
			&transfer.TransferDate,
			&transfer.Description,
			&transfer.Status,
		)
		if err != nil {
			http.Error(w, "Failed to parse result", http.StatusInternalServerError)
			log.Println("Row scan error:", err)
			return
		}

		// if description.Valid {
		// 	transfer.Description = &description.String
		// }
		// if status.Valid {
		// 	transfer.Status = &status.String
		// }

		transfers = append(transfers, transfer)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	helpers.WriteJSON(w,200,transfers)
}