package adminendpoints

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Apouzi/Golang-Admin-Service/helpers"
	"github.com/go-chi/chi/v5"
)

type AdminRoutes struct{
	DB *sql.DB
}

func InstanceAdminRoutes(  ) *AdminRoutes {
	r := &AdminRoutes{
	}
	return r
}



func(route *AdminRoutes)  AdminTableScopeCheck(adminTable string, tableName string ,adminID any, w http.ResponseWriter) bool{
	// strQ := "SELECT AdminID FROM" + adminTable + "WHERE Tablename = " + string(adminID) + " AND AdminID = " + adminID
	var exists bool
	var strBuild strings.Builder
	strBuild.WriteString("SELECT AdminID FROM ")
	strBuild.WriteString(adminTable)
	strBuild.WriteString(" WHERE TableName = ? AND AdminID = ?")
	route.DB.QueryRow(strBuild.String(), tableName, adminID).Scan(&exists)
	
	if exists == false{
		fmt.Println("Failed Query AdminTableScopeCheck endpoint")
		return false
	}

	return true
}



// Needs to get SKU, UPC, Primary Image to get created. Primary Image needs to be a google/AWS bucket
func(route *AdminRoutes) CreateProduct(w http.ResponseWriter, r *http.Request){

	fmt.Println("hit!")

	productFromRequest := &ProductCreate{}

	err := helpers.ReadJSON(w, r, &productFromRequest)
	if err != nil{
		fmt.Println("read json error:",err)
	}
	fmt.Println("create product at admin:", productFromRequest)
	sendOff, err := json.Marshal(productFromRequest)
	if err != nil{
		fmt.Println("There was an error marshalling data:",err)
	}
	createdProductResult, err := http.Post("http://dblayer:8080/db/products/","application/json",bytes.NewReader(sendOff))
	if err != nil{
		fmt.Println(err)
	}
	prodcreate:= &ProductCreateRetrieve{}
	
	decode := json.NewDecoder(createdProductResult.Body)
	decode.Decode(prodcreate)
	fmt.Println("attempted result is",prodcreate)
	helpers.WriteJSON(w,http.StatusAccepted,prodcreate)

}


func (route *AdminRoutes) CreateVariation(w http.ResponseWriter, r *http.Request){
	ProductID := chi.URLParam(r, "ProductID")
	variation := VariationCreate{}
	helpers.ReadJSON(w,r, &variation)

// Check if product exists, if not, then return false
	pil := ProductRetrieve{}
	url := "http://dblayer:8080/products/" + ProductID
	resp, err := http.Get(url)
	if err != nil{
		helpers.WriteJSON(w,500,"Error getting data from database")
	}

	if resp.StatusCode == 404{
		helpers.ErrorJSON(w,errors.New("could not find the coresponding product to create variation"), 404)
		return
	}
	

	jDecode := json.NewDecoder(resp.Body)
	if err = jDecode.Decode(&pil); err != nil || pil.ProductID == 0{
		fmt.Println("There is an error decoding!", err)
	}

	url = "http://dblayer:8080/products/" + ProductID + "/variation"
	varbytes, err:= json.Marshal(variation)
	if err != nil{
		helpers.ErrorJSON(w,errors.New("could not martial inputted variation"), 404)
	}
	varReader := bytes.NewReader(varbytes)
	resp, err = http.Post(url, "application/json",varReader)
	if err != nil{
		helpers.WriteJSON(w,500,"Error getting data from database")
	}
	if resp.StatusCode == 404{
		helpers.ErrorJSON(w,errors.New("could not find the coresponding product to create variation"), 404)
	}
	verify := variCrtd{}
	jDecode = json.NewDecoder(resp.Body)
	if err = jDecode.Decode(&verify); err != nil || pil.ProductID == 0{
		fmt.Println("There is an error decoding!", err)
	}

	helpers.WriteJSON(w,200,verify)
	 
	
}



func(route *AdminRoutes) CreateInventoryLocation(w http.ResponseWriter, r *http.Request){
	// Test for Variantion existness
	prodInvLoc := &ProdInvLocCreation{}
	helpers.ReadJSON(w,r,&prodInvLoc)
	
	variation_id := strconv.Itoa(int(prodInvLoc.VarID))
	
	url := "http://dblayer:8080/products/variation/" + variation_id
	fmt.Println("resulting url:", url)
	resp, err := http.Get(url)
	responsevariation := VariationRetrieve{}
	if err != nil || resp.StatusCode == http.StatusForbidden{
		helpers.ErrorJSON(w, errors.New("attempt to retrieve variation failed, could not retrieve varitation id"),resp.StatusCode)
		return
	}

	jDecode := json.NewDecoder(resp.Body)
	if err = jDecode.Decode(&responsevariation); err != nil{
		fmt.Println("There is an error decoding!", err)
	}
	
	prodinv, err := json.Marshal(prodInvLoc)
	if err != nil{
		fmt.Println("there was an error marshalling data")
		helpers.ErrorJSON(w, errors.New("there was an error marshalling data"),resp.StatusCode)
		return
	}

	prodinvreader := bytes.NewReader(prodinv)
	resp, err = http.Post("http://dblayer:8080/products/inventory-location","application/json",prodinvreader)
	if err != nil{
		helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed"),resp.StatusCode)
		return
	}
	respDecode := json.NewDecoder(resp.Body)
	pilReturn := &PILCreated{ }
	err = respDecode.Decode(pilReturn)
	if err != nil{
		helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed"),resp.StatusCode)
		return
	}
	
	helpers.WriteJSON(w, http.StatusAccepted, pilReturn)
}



func (route *AdminRoutes) CreatePrimeCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}

	url := "http://dblayer:8080/category/prime"
	catBytes, err := json.Marshal(category_read)
	if err != nil{
		fmt.Println(err)
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	catret := &CategoryReturn{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}
	
	helpers.WriteJSON(w,  200, catret)
}

func (route *AdminRoutes) CreateSubCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}


	url := "http://dblayer:8080/category/sub"
	catBytes, err := json.Marshal(category_read)
	if err != nil{
		fmt.Println(err)
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	catret := &CategoryReturn{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}
	
	helpers.WriteJSON(w, http.StatusCreated, catret)
}

func (route *AdminRoutes) CreateFinalCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}
	url := "http://dblayer:8080/category/final"
	catBytes, err := json.Marshal(category_read)
	if err != nil{
		fmt.Println(err)
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	catret := &CategoryReturn{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}
	
	helpers.WriteJSON(w, http.StatusCreated, catret)
}



func (route *AdminRoutes) ConnectPrimeToSubCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalSub := CatToCat{}
	err := helpers.ReadJSON(w,r, &FinalSub)
	if err != nil{
		fmt.Println(err)
	}


	url := "http://dblayer:8080/category/primetosub"
	catBytes, err := json.Marshal(FinalSub)
	if err != nil{
		fmt.Println(err)
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	catret := &CatToCat{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}
	
	helpers.WriteJSON(w, http.StatusCreated, catret)
}

func (route *AdminRoutes) ConnectSubToFinalCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalSub := CatToCat{}
	err := helpers.ReadJSON(w,r, &FinalSub)
	if err != nil{
		fmt.Println(err)
	}
	url := "http://dblayer:8080/category/subtofinal"
	catBytes, err := json.Marshal(FinalSub)
	if err != nil{
		fmt.Println(err)
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	catret := &CatToCat{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}

	helpers.WriteJSON(w, http.StatusAccepted, catret)
}


func (route *AdminRoutes) ConnectFinalToProdCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalProd := CatToProd{}
	err := helpers.ReadJSON(w,r, &FinalProd)
	if err != nil{
		fmt.Println(err)
	}
	url := "http://dblayer:8080/category/finaltoprod"
	catBytes, err := json.Marshal(FinalProd)
	if err != nil{
		fmt.Println(err)
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	catret := &CatToProd{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}

	helpers.WriteJSON(w, http.StatusAccepted, FinalProd)
}


func (route *AdminRoutes) InsertIntoFinalProd(w http.ResponseWriter, r *http.Request){
	ReadCatR := ReadCat{}
	err := helpers.ReadJSON(w,r,&ReadCatR)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("InsertIntoCategory ReadCatR",ReadCatR)
	FinalProd := "INSERT INTO tblCatFinalProd(CatFinalID, Product_ID) VALUES(?,?)"
	route.DB.Exec(FinalProd, 1,ReadCatR.Category)
}



func (route *AdminRoutes) ReturnAllPrimeCategories(w http.ResponseWriter, r *http.Request){
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

func (route *AdminRoutes) ReturnAllSubCategories(w http.ResponseWriter, r *http.Request){
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

func (route *AdminRoutes) ReturnAllFinalCategories(w http.ResponseWriter, r *http.Request){
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





func (route *AdminRoutes) EditProduct(w http.ResponseWriter, r *http.Request){
	ProdID := chi.URLParam(r, "ProductID")
	prodEdit := &ProductEdit{}
	helpers.ReadJSON(w,r, prodEdit)
	
	fmt.Println(prodEdit)
	url := "http://dblayer:8080/products/"+ProdID
	fmt.Println("url:", url)
	prodBytes, err := json.Marshal(prodEdit)
	if err != nil{
		fmt.Println(err)
	}
	// prodDecode:= bytes.NewReader(prodBytes)
	req, err := http.NewRequest("PATCH",url,bytes.NewReader(prodBytes))
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(prodEdit)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}


	helpers.WriteJSON(w,http.StatusAccepted,&prodEdit)
	
}



func (route *AdminRoutes) EditVariation(w http.ResponseWriter, r *http.Request){
	r.Header.Get("Authorization")
	VarID := chi.URLParam(r, "VariationID")
	VaritEdit := &VariationEdit{}
	helpers.ReadJSON(w,r, VaritEdit)
	url := "http://dblayer:8080/variation/" + VarID

	varitBytes, err := json.Marshal(VaritEdit)
	if err != nil{
		fmt.Println(err)
	}
	// prodDecode:= bytes.NewReader(prodBytes)
	req, err := http.NewRequest("PATCH",url,bytes.NewReader(varitBytes))
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	responseDecode := json.NewDecoder(resp.Body)
	varitReturn := &VariationEdit{}
	err = responseDecode.Decode(varitReturn)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}


	helpers.WriteJSON(w, http.StatusAccepted, VaritEdit)
}



func (route *AdminRoutes) DeletePrimeCategory(w http.ResponseWriter, r *http.Request){
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


func (route *AdminRoutes) DeleteSubCategory(w http.ResponseWriter, r *http.Request){
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


func (route *AdminRoutes) DeleteFinalCategory(w http.ResponseWriter, r *http.Request){
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



func (route *AdminRoutes) AddAttribute(w http.ResponseWriter, r *http.Request){
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

	url := "http://dblayer:8080/variation/" + VarID + "/attribute"

	attributeBytes, err := json.Marshal(att)
	if err != nil{
		fmt.Println(err)
	}
	// prodDecode:= bytes.NewReader(prodBytes)
	req, err := http.NewRequest("POST",url,bytes.NewReader(attributeBytes))
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	responseDecode := json.NewDecoder(resp.Body)
	attRet := &AddedSendBack{}
	err = responseDecode.Decode(attRet)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}
	
	helpers.WriteJSON(w, 200, attRet)
}

func (route *AdminRoutes) DeleteAttribute(w http.ResponseWriter, r *http.Request){
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


func (route *AdminRoutes) UpdateAttribute(w http.ResponseWriter, r *http.Request){
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



func (route *AdminRoutes) GetAllTables(w http.ResponseWriter, r *http.Request){
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


func(route *AdminRoutes) UserToAdmin(w http.ResponseWriter, r *http.Request){
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

