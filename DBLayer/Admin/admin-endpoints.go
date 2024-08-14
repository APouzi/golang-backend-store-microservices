package admin

import (
	"database/sql"
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

type ProductRoutesTray struct{
	DB *sql.DB
}

func GetProductRouteInstance(dbInst *sql.DB) *ProductRoutesTray{
	
	// routeMap := prepareProductRoutes(dbInst)
	return &ProductRoutesTray{
		// getAllProductsStmt: routeMap["getAllProducts"],
		DB: dbInst,
	}
}

func (adminProdRoutes *ProductRoutesTray) CreateProductMultiChain(w http.ResponseWriter, r *http.Request) {
	transaction, err := adminProdRoutes.DB.Begin()
	if err != nil {
		log.Println("Error creating a transation in CreateProduct")
		log.Println(err)
	}

	productRetrieve := &ProductCreate{}

	helpers.ReadJSON(w, r, &productRetrieve)
	fmt.Println("product retrieve at db:",productRetrieve)
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
	tRes, err = transaction.Exec("INSERT INTO tblProductVariation(Product_ID,Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)", prodID, productRetrieve.VariationName, productRetrieve.VariationDescription, productRetrieve.VariationPrice)
	if err != nil {
		fmt.Println("transaction at tblProductVariation has failed")
		fmt.Println(err)
		transaction.Rollback()
		return
	}

	ProdVarID, err := tRes.LastInsertId()
	if err != nil {
		fmt.Println("retrieval of LastInsertID of tblProductVariation has failed")
		fmt.Println(err)
		transaction.Rollback()
		return
	}
	PCR := ProductCreateRetrieve{
		ProductID: prodID,
		VarID:     ProdVarID,
	}
	if productRetrieve.LocationAt == "" {

		err = transaction.Commit()
		if err != nil {
			fmt.Println(err)
			transaction.Rollback()
			return
		}
		helpers.WriteJSON(w, http.StatusAccepted, &PCR)
		return
	}

	tRes, err = transaction.Exec("INSERT INTO tblProductInventoryLocation(Variation_ID, Quantity, Location_AT) VALUES(?,?,?)", ProdVarID, productRetrieve.VariationQuantity, productRetrieve.LocationAt)
	if err != nil {
		fmt.Println("transaction at tblProductInventory has failed")
		fmt.Println(err)
	}
	invID, err := tRes.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}
	PCR.ProdInvLoc = invID
	err = transaction.Commit()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("print of prod created",PCR)
	helpers.WriteJSON(w,http.StatusAccepted,PCR)
}




func (adminProdRoutes *ProductRoutesTray) CreateProductVariation(w http.ResponseWriter, r *http.Request) {

	ProductID := chi.URLParam(r, "ProductID")
	fmt.Println("productID", ProductID)
	variation := VariationCreate{PrimaryImage: ""}
	helpers.ReadJSON(w,r, &variation)
	varitCrt := variCrtd{}
	if variation.PrimaryImage != "" {
		
		varit, err := adminProdRoutes.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)", ProductID,variation.Name, variation.Description, variation.Price)
		if err != nil{
			log.Println("insert into tblProductVariation failed")
			log.Println(err)
			helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed"),400)
			return
		}
		varitCrt.VariationID, err = varit.LastInsertId()
		if err != nil{
			log.Println(err)
			helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed, could not retrieve varitation id"),400)
			return
		}
		helpers.WriteJSON(w, http.StatusCreated,varitCrt)
	}
	prodid, err := strconv.Atoi(ProductID)
	varit, err := adminProdRoutes.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)", prodid,variation.Name, variation.Description, variation.Price)
	if err != nil{
		fmt.Println("insert into tblProductVariation failed")
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed, could not retrieve varitation id"),400)
	}
	varitCrt.VariationID, err = varit.LastInsertId()
	
	helpers.WriteJSON(w, http.StatusCreated,varitCrt)
}




func (adminProdRoutes *ProductRoutesTray) CreateInventoryLocation(w http.ResponseWriter, r *http.Request) {

	prodInvLoc := ProdInvLocCreation{}
	helpers.ReadJSON(w,r,&prodInvLoc)
	
	res ,err:= adminProdRoutes.DB.Exec("INSERT INTO tblProductInventoryLocation(Variation_ID, Quantity, Location_At) VALUES(?,?,?)", prodInvLoc.VarID,prodInvLoc.Quantity,prodInvLoc.Location)
	
	if err != nil{
		fmt.Println("failed to create tblProductInventoryLocation")
		fmt.Println(err)
		helpers.ErrorJSON(w,err,http.StatusForbidden)
		return
	}

	pilID, err := res.LastInsertId()
	
	if err != nil{
		fmt.Println("result of tblProductInventoryLocation failed")
	}
	pilReturn := PILCreated{InvID:pilID, Quantity: prodInvLoc.Quantity, Location: prodInvLoc.Location }
	helpers.WriteJSON(w, http.StatusAccepted, pilReturn)
}


func (adminProdRoutes *ProductRoutesTray) CreatePrimeCategory(w http.ResponseWriter, r *http.Request) {
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}

	fmt.Println("in dblayer: CreatePrimeCategory",category_read)
	result, err := adminProdRoutes.DB.Exec("INSERT INTO tblCategoriesPrime(CategoryName, CategoryDescription) VALUES(?,?)", category_read.CategoryName, category_read.CategoryDescription )
	if err != nil{
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}

	catReturn := CategoryReturn{CategoryId: resultID, CategoryName: category_read.CategoryName, CategoryDescription: category_read.CategoryDescription}
	err = helpers.WriteJSON(w,200,catReturn)
	if err != nil{
		fmt.Println("There was an error trying to send data back",err)
	}

}

func (route *ProductRoutesTray) CreateSubCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}

	
	result, err := route.DB.Exec("INSERT INTO tblCategoriesSub(CategoryName, CategoryDescription) VALUES(?,?)", category_read.CategoryName, category_read.CategoryDescription )
	if err != nil{
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	
	catReturn := CategoryReturn{CategoryId: resultID, CategoryName: category_read.CategoryName, CategoryDescription: category_read.CategoryDescription}
	err = helpers.WriteJSON(w,200,catReturn)
	if err != nil{
		fmt.Println("There was an error trying to send data back",err)
	}
}

func (route *ProductRoutesTray) CreateFinalCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCategoriesFinal(CategoryName, CategoryDescription) VALUES(?,?)", category_read.CategoryName, category_read.CategoryDescription )
	if err != nil{
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	
	catReturn := CategoryReturn{CategoryId: resultID, CategoryName: category_read.CategoryName, CategoryDescription: category_read.CategoryDescription}
	err = helpers.WriteJSON(w,200,catReturn)
	if err != nil{
		fmt.Println("There was an error trying to send data back",err)
	}
}


func (route *ProductRoutesTray) ConnectPrimeToSubCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalSub := CatToCat{}
	err := helpers.ReadJSON(w,r, &FinalSub)
	if err != nil{
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCatPrimeSub(CatPrimeID,  CatSubID) VALUES(?,?)", FinalSub.CatStart, FinalSub.CatEnd)
	if err != nil{
		fmt.Println(err)
	}
	_, err = result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	
	helpers.WriteJSON(w, http.StatusAccepted, FinalSub)
}

func (route *ProductRoutesTray) ConnectSubToFinalCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalSub := CatToCat{}
	err := helpers.ReadJSON(w,r, &FinalSub)
	if err != nil{
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCatSubFinal(CatSubID, CatFinalID) VALUES(?,?)", FinalSub.CatStart, FinalSub.CatEnd)
	if err != nil{
		fmt.Println(err)
	}

	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	helpers.WriteJSON(w, http.StatusAccepted, resultID)
}



func (route *ProductRoutesTray) ConnectFinalToProdCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalProd := CatToProd{}
	err := helpers.ReadJSON(w,r, &FinalProd)
	if err != nil{
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCatFinalProd(CatFinalID, Product_ID) VALUES(?,?)", FinalProd.Cat, FinalProd.Prod)

	if err != nil{
		fmt.Println(err)
	}

	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	helpers.WriteJSON(w, http.StatusAccepted, resultID)
}


func (route *ProductRoutesTray) InsertIntoFinalProd(w http.ResponseWriter, r *http.Request){
	ReadCatR := ReadCat{}
	err := helpers.ReadJSON(w,r,&ReadCatR)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("InsertIntoCategory ReadCatR",ReadCatR)
	FinalProd := "INSERT INTO tblCatFinalProd(CatFinalID, Product_ID) VALUES(?,?)"
	route.DB.Exec(FinalProd, 1,ReadCatR.Category)
}



func (route *ProductRoutesTray) ReturnAllPrimeCategories(w http.ResponseWriter, r *http.Request){
	query := "SELECT CategoryName, CategoryDescription FROM tblCategoriesPrime"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategoryReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategoryReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)

}

func (route *ProductRoutesTray) ReturnAllSubCategories(w http.ResponseWriter, r *http.Request){
	query := "SELECT CategoryName, CategoryDescription FROM tblCategoriesSub"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategoryReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategoryReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)
}

func (route *ProductRoutesTray) ReturnAllFinalCategories(w http.ResponseWriter, r *http.Request){
	query := "SELECT CategoryName, CategoryDescription FROM tblCategoriesFinal"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategoryReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategoryReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)
}





func (route *ProductRoutesTray) EditProduct(w http.ResponseWriter, r *http.Request){
	ProdID := chi.URLParam(r, "ProductID")
	prodEdit := ProductEdit{}
	helpers.ReadJSON(w,r, &prodEdit)
	var buf strings.Builder
	buf.WriteString("UPDATE tblProducts SET")
	var count int = 0
	Varib := []any{}
	if prodEdit.Name != "" {
		if count == 0{
			buf.WriteString(" Product_Name = ?")
			Varib = append(Varib, prodEdit.Name)
			count++
		}
		buf.WriteString(", Product_Name = ?")
		Varib = append(Varib, prodEdit.Name)
	}
	if prodEdit.Description != "" {
		if count == 0{
			buf.WriteString(" Product_Description = ?")
			Varib = append(Varib, prodEdit.Description)
			count++
		}
		buf.WriteString(", Product_Description = ?")
		Varib = append(Varib, prodEdit.Description)
	}
	if count  == 0 {
		helpers.WriteJSON(w,http.StatusAccepted,"failed")
		return
	}

	buf.WriteString(", Modified_Date = ? WHERE Product_ID = ?")
	Varib = append(Varib, time.Now(),ProdID)
	_, err := route.DB.Exec(buf.String(), Varib...)
	if err != nil{
		fmt.Println("err with exec Edit Product Update")
		fmt.Println(err)
	}

	helpers.WriteJSON(w,http.StatusAccepted,&prodEdit)
	
}



func (route *ProductRoutesTray) EditVariation(w http.ResponseWriter, r *http.Request){
	r.Header.Get("Authorization")
	VarID := chi.URLParam(r, "VariationID")
	VaritEdit := VariationEdit{}
	helpers.ReadJSON(w,r, &VaritEdit)
	var buf strings.Builder
	Varib := []any{}
	buf.WriteString("UPDATE tblProductVariation SET")
	var count int = 0
	if VaritEdit.VariationName != "" {
		if count == 0{
			buf.WriteString(" Variation_Name = ?")
			Varib = append(Varib, VaritEdit.VariationName)
			count++
		}
		buf.WriteString(", Variation_Name = ?")
		Varib = append(Varib, VaritEdit.VariationName)
	}
	if VaritEdit.VariationDescription != ""{
		if count == 0{
			buf.WriteString(" Variation_Description = ?")
			Varib = append(Varib, VaritEdit.VariationDescription)
			count++
		}
		buf.WriteString(", Variation_Description = ?")
		Varib = append(Varib, VaritEdit.VariationDescription)
	}
	if VaritEdit.SKU != ""{
		if count == 0 {
			buf.WriteString(" SKU = ?")
			Varib = append(Varib, VaritEdit.SKU)
			count++
		}
		buf.WriteString(", SKU = ?")
		Varib = append(Varib, VaritEdit.SKU)
	}
	if VaritEdit.UPC != ""{
		if count == 0{
			buf.WriteString(" UPC = ?")
			Varib = append(Varib, VaritEdit.UPC)
			count++
		}
		buf.WriteString(", UPC = ?")
		Varib = append(Varib, VaritEdit.UPC)
	}
	if VaritEdit.VariationPrice != 0 {
		if count == 0{
			buf.WriteString(" Variation_Price = ?")
			Varib = append(Varib, VaritEdit.VariationPrice)
			count++
		}
		buf.WriteString(", Variation_Price = ?")
		Varib = append(Varib, VaritEdit.VariationPrice)
	}
	buf.WriteString(" WHERE Variation_ID = ?")
	Varib = append(Varib, VarID)
	_,err := route.DB.Exec(buf.String(),Varib...)
	if err != nil{
		fmt.Println(err)
	}
	helpers.WriteJSON(w, http.StatusAccepted, VaritEdit)
}



func (route *ProductRoutesTray) DeletePrimeCategory(w http.ResponseWriter, r *http.Request){
	CatName := chi.URLParam(r,"CatPrimeName")
	if CatName == ""{
		fmt.Println("No CatPrimeName wasn't pulled")
		return
	}
	_, err := route.DB.Exec("DELETE FROM tblCategoriesPrime WHERE CategoryName = ?", CatName)
	if err != nil{
		fmt.Println("Failed deletion in CatPrimeName")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}

	sendBack := DeletedSendBack{SendBack:false}
	helpers.WriteJSON(w,200,sendBack)
}


func (route *ProductRoutesTray) DeleteSubCategory(w http.ResponseWriter, r *http.Request){
	CatName := chi.URLParam(r,"CatSubName")
	if CatName == ""{
		fmt.Println("No CatSubName wasn't pulled")
		return
	}
	
	_, err := route.DB.Exec("DELETE FROM tblCategoriesSub WHERE CategoryName = ?", CatName)
	if err != nil{
		fmt.Println("Failed deletion in CatSubName")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}
	
	sendBack := DeletedSendBack{SendBack:false}
	helpers.WriteJSON(w,200,sendBack)
}


func (route *ProductRoutesTray) DeleteFinalCategory(w http.ResponseWriter, r *http.Request){
	CatName := chi.URLParam(r,"CatFinalName")
	if CatName == ""{
		fmt.Println("No CatPrimeName wasn't pulled")
		return
	}
	
	_, err := route.DB.Exec("DELETE FROM tblCategoriesFinal WHERE CategoryName = ?", CatName)
	if err != nil{
		fmt.Println("Failed deletion in CatPrimeName")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}

	sendBack := DeletedSendBack{SendBack:false}
	helpers.WriteJSON(w,200,sendBack)
}



func (route *ProductRoutesTray) AddAttribute(w http.ResponseWriter, r *http.Request){
	VarID := chi.URLParam(r,"VariationID")
	if VarID == ""{
		helpers.ErrorJSON(w, errors.New("please input VariationID"),400)
		return
	}
	att := Attribute{}
	
	err := helpers.ReadJSON(w,r,&att)
	if err != nil{
		helpers.ErrorJSON(w, err, 500)
		return
	}
	sql, err := route.DB.Exec("INSERT INTO tblProductAttribute (Variation_ID, AttributeName) VALUES(?,?)",VarID,att.Attribute)
	if err != nil{
		helpers.ErrorJSON(w,err, 400)
		return
	}
	var id int64
	id, err = sql.LastInsertId()
	if err != nil{
		helpers.ErrorJSON(w,errors.New("failed attribute LastInsertID"))
		return
	}
	sendBack := AddedSendBack{IDSendBack: id}
	helpers.WriteJSON(w, 200, sendBack)
}

func (route *ProductRoutesTray) DeleteAttribute(w http.ResponseWriter, r *http.Request){
	VarID := chi.URLParam(r,"VariationID")
	AttName := chi.URLParam(r, "AttributeName")
	if VarID == ""{
		helpers.ErrorJSON(w, errors.New("please input VariationID"),400)
		return
	}

	sql, err := route.DB.Exec("DELETE FROM tblProductAttribute WHERE Variation_ID = ? AND AttributeName = ?", VarID, AttName)
	if err != nil{
		helpers.ErrorJSON(w,err, 400)
		return
	}

	nRows, _ := sql.RowsAffected()
	if nRows < 1{
		helpers.WriteJSON(w, 200, "Not Deleted")
		return
	}
	
	helpers.WriteJSON(w, 200, "Deleted")
}


func (route *ProductRoutesTray) UpdateAttribute(w http.ResponseWriter, r *http.Request){
	VarID := chi.URLParam(r,"VariationID")
	AttName := chi.URLParam(r, "AttributeName")
	if VarID == ""{
		helpers.ErrorJSON(w, errors.New("please input VariationID"),400)
		return
	}
	AttRead := Attribute{}
	helpers.ReadJSON(w,r,&AttRead)
	fmt.Println(AttRead.Attribute,"variatio id and atttribute name", VarID,AttName)
	sql, err := route.DB.Exec("UPDATE tblProductAttribute SET AttributeName = ? WHERE Variation_ID = ? AND AttributeName = ?",AttRead.Attribute ,VarID, AttName)
	if err != nil{
		helpers.ErrorJSON(w,err, 400)
		return
	}

	nRows, _ := sql.RowsAffected()
	if nRows < 1{
		helpers.WriteJSON(w, 200, "No Updated Happened")
		return
	}
	
	helpers.WriteJSON(w, 200, "Attribute Updated")
}



func (route *ProductRoutesTray) GetAllTables(w http.ResponseWriter, r *http.Request){
	sql,err := route.DB.Query("show tables")
	if err != nil{
		fmt.Println("failed to get all tables")
		return
	}
	var table string
	list := []string{}
	for sql.Next(){
		sql.Scan(&table)
		list = append(list, table)
	}
	helpers.WriteJSON(w,200,list)
}


func(route *ProductRoutesTray) UserToAdmin(w http.ResponseWriter, r *http.Request){
	id := chi.URLParam(r,"UserID")
	fmt.Println("UserToAdmin:",id)
	var exists bool
	route.DB.QueryRow("SELECT UserID FROM tblUser WHERE UserID = ?",id).Scan(&exists)
	if exists == false {
		helpers.ErrorJSON(w,errors.New("user doesn't exist") ,400)
		return
	}

	var UserID int64
	err := route.DB.QueryRow("SELECT UserID FROM tblUser WHERE UserID = ?", id).Scan(&UserID)
	if err != nil{
		helpers.ErrorJSON(w,errors.New("issue with scanning user into struct ") ,500)
		return
	}

	sql, err := route.DB.Exec("INSERT INTO tblAdminUsers (UserID, SuperUser) VALUES(?,?)",UserID,false)
	if err != nil{
		helpers.ErrorJSON(w,errors.New("failed insertinginto tblAdminUsers") ,500)
		return
	}
	type returnAdminID struct{
		UserID int64 `json:"AdminUserID"`
	}
	adminID, err := sql.LastInsertId()
	if err != nil{
		helpers.ErrorJSON(w,errors.New("couldn't retrieve id from LastInsertId") ,500)
		return
	}
	rAID := returnAdminID{UserID:adminID}
	helpers.WriteJSON(w,200,rAID)
}
