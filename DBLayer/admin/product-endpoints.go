package admin

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/APouzi/DBLayer/helpers"
	"github.com/go-chi/chi/v5"
)

func (adminProdRoutes *ProductRoutesTray) CreateProductMultiChain(w http.ResponseWriter, r *http.Request) {
	transaction, err := adminProdRoutes.DB.Begin()
	if err != nil {
		log.Println("Error creating a transation in CreateProduct")
		log.Println(err)
	}

	productRetrieve := &ProductCreate{}

	helpers.ReadJSON(w, r, &productRetrieve)
	fmt.Println("product retrieve at db:", productRetrieve)
	tRes, err := transaction.Exec("INSERT INTO tblProducts(Product_Name, Product_Description) VALUES(?,?)", productRetrieve.Name, productRetrieve.Description)
	if err != nil {
		fmt.Println("transaction at tblProduct has failed")
		fmt.Println(err)
		transaction.Rollback()
	}
	prodID, err := tRes.LastInsertId()
	if err != nil {
		fmt.Println("retrieval of LastInsertID of tblProduct has failed")
		fmt.Println(err)
		transaction.Rollback()
		return
	}
	var ProdVarID int64
	for _, variation := range productRetrieve.Variations {
		tRes, err = transaction.Exec("INSERT INTO tblProductVariation(Product_ID,Variation_Name, Variation_Description) VALUES(?,?,?)", prodID, variation.Name, variation.Description)
		if err != nil {
			fmt.Println("transaction at tblProductVariation has failed")
			fmt.Println(err)
			transaction.Rollback()
			return
		}

		ProdVarID, err = tRes.LastInsertId()
		if err != nil {
			fmt.Println("retrieval of LastInsertID of tblProductVariation has failed")
			fmt.Println(err)
			transaction.Rollback()
			return
		}

		// Insert default size with price
		_, err = transaction.Exec("INSERT INTO tblProductSize(Size_Name, Variation_ID, Variation_Price, Price) VALUES(?, ?, ?, ?)", "Standard", ProdVarID, variation.Price, variation.Price)
		if err != nil {
			fmt.Println("transaction at tblProductSize has failed")
			fmt.Println(err)
			transaction.Rollback()
			return
		}

		tRes, err = transaction.Exec("INSERT INTO tblProductInventoryLocation(Variation_ID, Quantity, Location_AT) VALUES(?,?,?)", ProdVarID, variation.VariationQuantity, variation.LocationAt)
		if err != nil {
			fmt.Println("transaction at tblProductInventory has failed")
			fmt.Println(err)
		}
	}

	PCR := ProductCreateRetrieve{
		ProductID: prodID,
		VarID:     ProdVarID,
	}

	err = transaction.Commit()
	if err != nil {
		fmt.Println(err)
		transaction.Rollback()
		return
	}
	helpers.WriteJSON(w, http.StatusAccepted, &PCR)

	invID, err := tRes.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}
	PCR.ProdInvLoc = invID
	err = transaction.Commit()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("print of prod created", PCR)
	helpers.WriteJSON(w, http.StatusAccepted, PCR)
}

func (adminProdRoutes *ProductRoutesTray) CreateProductVariation(w http.ResponseWriter, r *http.Request) {

	ProductID := chi.URLParam(r, "ProductID")
	fmt.Println("productID", ProductID)
	variation := VariationCreate{PrimaryImage: ""}
	helpers.ReadJSON(w, r, &variation)
	varitCrt := variCrtd{}
	if variation.PrimaryImage != "" {

		varit, err := adminProdRoutes.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description) VALUES(?,?,?)", ProductID, variation.Name, variation.Description)
		if err != nil {
			log.Println("insert into tblProductVariation failed")
			log.Println(err)
			helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed"), 400)
			return
		}
		varitCrt.VariationID, err = varit.LastInsertId()
		if err != nil {
			log.Println(err)
			helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed, could not retrieve varitation id"), 400)
			return
		}

		_, err = adminProdRoutes.DB.Exec("INSERT INTO tblProductSize(Size_Name, Variation_ID, Variation_Price, Price) VALUES(?, ?, ?, ?)", "Standard", varitCrt.VariationID, variation.Price, variation.Price)
		if err != nil {
			log.Println("insert into tblProductSize failed", err)
		}

		helpers.WriteJSON(w, http.StatusCreated, varitCrt)
	}
	prodid, err := strconv.Atoi(ProductID)
	varit, err := adminProdRoutes.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description) VALUES(?,?,?)", prodid, variation.Name, variation.Description)
	if err != nil {
		fmt.Println("insert into tblProductVariation failed")
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed, could not retrieve varitation id"), 400)
		return
	}
	varitCrt.VariationID, err = varit.LastInsertId()

	_, err = adminProdRoutes.DB.Exec("INSERT INTO tblProductSize(Size_Name, Variation_ID, Variation_Price, Price) VALUES(?, ?, ?, ?)", "Standard", varitCrt.VariationID, variation.Price, variation.Price)
	if err != nil {
		fmt.Println("insert into tblProductSize failed", err)
	}

	helpers.WriteJSON(w, http.StatusCreated, varitCrt)
}

func (route *ProductRoutesTray) EditProduct(w http.ResponseWriter, r *http.Request) {
	ProdID := chi.URLParam(r, "ProductID")
	prodEdit := ProductEdit{}
	if err := helpers.ReadJSON(w, r, &prodEdit); err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	var buf strings.Builder
	buf.WriteString("UPDATE tblProducts SET")
	var count int = 0
	Varib := []any{}
	if prodEdit.Name != nil {
		if count > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" Product_Name = ?")
		Varib = append(Varib, *prodEdit.Name)
		count++
	}
	if prodEdit.Description != nil {
		if count > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" Product_Description = ?")
		Varib = append(Varib, *prodEdit.Description)
		count++
	}
	if prodEdit.PrimaryImage != nil {
		if count > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" PRIMARY_IMAGE = ?")
		Varib = append(Varib, *prodEdit.PrimaryImage)
		count++
	}
	if count == 0 {
		helpers.ErrorJSON(w, errors.New("no fields provided for update: expected at least one of product_name, product_description, primary_image"), http.StatusBadRequest)
		return
	}

	buf.WriteString(", Modified_Date = ? WHERE Product_ID = ?")
	Varib = append(Varib, time.Now().UTC(), ProdID)
	_, err := route.DB.Exec(buf.String(), Varib...)
	if err != nil {
		log.Println("err with exec Edit Product Update:", err)
		helpers.ErrorJSON(w, errors.New("failed to update product"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, &prodEdit)
}

func (route *ProductRoutesTray) EditVariation(w http.ResponseWriter, r *http.Request) {
	r.Header.Get("Authorization")
	VarID := chi.URLParam(r, "VariationID")
	VaritEdit := VariationEdit{}
	helpers.ReadJSON(w, r, &VaritEdit)
	fmt.Println("the recieved payload:",VaritEdit)
	var buf strings.Builder
	Varib := []any{}
	buf.WriteString("UPDATE tblProductVariation SET ")
	var count int = 0
	if VaritEdit.VariationName != "" {
		if count > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("Variation_Name = ?")
		Varib = append(Varib, VaritEdit.VariationName)
		count++
	}
	if VaritEdit.VariationDescription != "" {
		if count > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("Variation_Description = ?")
		Varib = append(Varib, VaritEdit.VariationDescription)
		count++
	}

	if count > 0 {
		buf.WriteString(" WHERE Variation_ID = ?")
		Varib = append(Varib, VarID)
		_, err := route.DB.Exec(buf.String(), Varib...)
		if err != nil {
			fmt.Println(err)
		}
	}

	if VaritEdit.VariationPrice != 0 {
		_, err := route.DB.Exec("UPDATE tblProductSize SET Variation_Price = ?, Price = ? WHERE Variation_ID = ?", VaritEdit.VariationPrice, VaritEdit.VariationPrice, VarID)
		if err != nil {
			fmt.Println("Error updating price in tblProductSize:", err)
		}
	}

	helpers.WriteJSON(w, http.StatusAccepted, VaritEdit)
}

func (route *ProductRoutesTray) DeleteVariation(w http.ResponseWriter, r *http.Request) {
	r.Header.Get("Authorization")
	VarID := chi.URLParam(r, "VariationID")
	fmt.Println("variation id to delete:", VarID)
	VaritDelete := VariationEdit{}
	helpers.ReadJSON(w, r, &VaritDelete)

	_, err := route.DB.Exec("DELETE FROM tblProductVariation WHERE Variation_ID = ?", VarID)
	if err != nil {
		fmt.Println("Error deleting variation:", err)
		helpers.ErrorJSON(w, errors.New("there was an error deleting the variation"), 500)
		return
	}

	helpers.WriteJSON(w, http.StatusAccepted, VaritDelete)
}

func (route *ProductRoutesTray) CreateProductSize(w http.ResponseWriter, r *http.Request) {
	var prdSize ProductSize

	helpers.ReadJSON(w, r, &prdSize)
	fmt.Println("product size payload:", prdSize)
	// Ensure DateCreated and ModifiedDate are initialized to now
	now := time.Now().UTC()
	prdSize.DateCreated = &now
	prdSize.ModifiedDate = &now

	// Use SQL column names consistent with DB schema - singular table name and Date_Created/Modified_Date
	sql, err := route.DB.Exec("INSERT INTO tblProductSize (Size_Name, Size_Description, Variation_ID, Variation_Price, SKU, UPC,Price, PRIMARY_IMAGE, Date_Created, Modified_Date) VALUES(?,?,?,?,?,?,?,?,?,?)",
		prdSize.SizeName, prdSize.SizeDescription, prdSize.VariationID, prdSize.VariationPrice, prdSize.SKU, prdSize.UPC, prdSize.Price, prdSize.PrimaryImage, prdSize.DateCreated, prdSize.ModifiedDate)
	if err != nil {
		fmt.Println("There was an error inserting product size:", err)
		helpers.ErrorJSON(w, errors.New("there was an error inserting product size"), 500)
		return
	}

	sizeID, err := sql.LastInsertId()
	if err != nil {
		fmt.Println("There was an error getting last insert id for product size:", err)
		helpers.ErrorJSON(w, errors.New("there was an error getting last insert id for product size"), 500)
		return
	}

	prdSize.SizeID = &sizeID

	helpers.WriteJSON(w, http.StatusCreated, prdSize)
}

func (route *ProductRoutesTray) EditProductSize(w http.ResponseWriter, r *http.Request) {
	prodsizeID := chi.URLParam(r, "ProductSizeID")
	if prodsizeID == "" {
		helpers.ErrorJSON(w, errors.New("please input ProductSizeID"), http.StatusBadRequest)
		return
	}

	sizeID, err := strconv.ParseInt(prodsizeID, 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("invalid ProductSizeID"), http.StatusBadRequest)
		return
	}

	var prdSize ProductSize
	if err := helpers.ReadJSON(w, r, &prdSize); err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	fmt.Println("edit product size payload:", prdSize)

	setClauses := []string{}
	args := []any{}

	if prdSize.SizeName != nil {
		setClauses = append(setClauses, "Size_Name = ?")
		args = append(args, *prdSize.SizeName)
	}
	if prdSize.SizeDescription != nil {
		setClauses = append(setClauses, "Size_Description = ?")
		args = append(args, *prdSize.SizeDescription)
	}
	if prdSize.VariationPrice != nil {
		setClauses = append(setClauses, "Variation_Price = ?")
		args = append(args, *prdSize.VariationPrice)
	}
	if prdSize.SKU != nil {
		setClauses = append(setClauses, "SKU = ?")
		args = append(args, *prdSize.SKU)
	}
	if prdSize.UPC != nil {
		setClauses = append(setClauses, "UPC = ?")
		args = append(args, *prdSize.UPC)
	}
	if prdSize.Price != nil {
		setClauses = append(setClauses, "Price = ?")
		args = append(args, *prdSize.Price)
	}
	if prdSize.PrimaryImage != nil {
		setClauses = append(setClauses, "PRIMARY_IMAGE = ?")
		args = append(args, *prdSize.PrimaryImage)
	}

	if len(setClauses) == 0 {
		helpers.ErrorJSON(w, errors.New("no fields provided for update"), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	setClauses = append(setClauses, "Modified_Date = ?")
	args = append(args, now, sizeID)

	query := "UPDATE tblProductSize SET " + strings.Join(setClauses, ", ") + " WHERE Size_ID = ?"
	result, err := route.DB.Exec(query, args...)
	if err != nil {
		fmt.Println("There was an error updating product size:", err)
		helpers.ErrorJSON(w, errors.New("there was an error updating product size"), 500)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		helpers.ErrorJSON(w, errors.New("product size not found"), http.StatusNotFound)
		return
	}

	prdSize.SizeID = &sizeID
	prdSize.ModifiedDate = &now
	helpers.WriteJSON(w, http.StatusOK, prdSize)
}

func (route *ProductRoutesTray) AddAttribute(w http.ResponseWriter, r *http.Request) {
	VarID := chi.URLParam(r, "VariationID")
	if VarID == "" {
		helpers.ErrorJSON(w, errors.New("please input VariationID"), 400)
		return
	}
	att := Attribute{}

	err := helpers.ReadJSON(w, r, &att)
	if err != nil {
		helpers.ErrorJSON(w, err, 500)
		return
	}
	sql, err := route.DB.Exec("INSERT INTO tblProductAttribute (Variation_ID, AttributeName) VALUES(?,?)", VarID, att.Attribute)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		return
	}
	var id int64
	id, err = sql.LastInsertId()
	if err != nil {
		helpers.ErrorJSON(w, errors.New("failed attribute LastInsertID"))
		return
	}
	sendBack := AddedSendBack{IDSendBack: id}
	helpers.WriteJSON(w, 200, sendBack)
}

func (route *ProductRoutesTray) DeleteAttribute(w http.ResponseWriter, r *http.Request) {
	VarID := chi.URLParam(r, "VariationID")
	AttName := chi.URLParam(r, "AttributeName")
	if VarID == "" {
		helpers.ErrorJSON(w, errors.New("please input VariationID"), 400)
		return
	}

	sql, err := route.DB.Exec("DELETE FROM tblProductAttribute WHERE Variation_ID = ? AND AttributeName = ?", VarID, AttName)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		return
	}

	nRows, _ := sql.RowsAffected()
	if nRows < 1 {
		helpers.WriteJSON(w, 200, "Not Deleted")
		return
	}

	helpers.WriteJSON(w, 200, "Deleted")
}

func (route *ProductRoutesTray) DeleteProductSize(w http.ResponseWriter, r *http.Request) {
	VarID := chi.URLParam(r, "ProductSizeID")

	if VarID == "" {
		helpers.ErrorJSON(w, errors.New("please input ProductSizeID"), 400)
		return
	}

	sql, err := route.DB.Exec("DELETE FROM tblProductSize WHERE Size_ID = ?", VarID)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		return
	}

	nRows, _ := sql.RowsAffected()
	if nRows < 1 {
		helpers.WriteJSON(w, 200, "Not Deleted")
		return
	}

	helpers.WriteJSON(w, 200, "Deleted")
}

func (route *ProductRoutesTray) UpdateAttribute(w http.ResponseWriter, r *http.Request) {
	VarID := chi.URLParam(r, "VariationID")
	AttName := chi.URLParam(r, "AttributeName")
	if VarID == "" {
		helpers.ErrorJSON(w, errors.New("please input VariationID"), 400)
		return
	}
	AttRead := Attribute{}
	helpers.ReadJSON(w, r, &AttRead)
	fmt.Println(AttRead.Attribute, "variatio id and atttribute name", VarID, AttName)
	sql, err := route.DB.Exec("UPDATE tblProductAttribute SET AttributeName = ? WHERE Variation_ID = ? AND AttributeName = ?", AttRead.Attribute, VarID, AttName)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		return
	}

	nRows, _ := sql.RowsAffected()
	if nRows < 1 {
		helpers.WriteJSON(w, 200, "No Updated Happened")
		return
	}

	helpers.WriteJSON(w, 200, "Attribute Updated")
}
