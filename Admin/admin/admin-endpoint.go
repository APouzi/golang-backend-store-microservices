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

func InstanceAdminRoutes() *AdminRoutes {
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

// func(route *AdminRoutes) CreatePrimeCategory(w http.ResponseWriter, r *http.Request){
// 	var cat Category
// 	helpers.ReadJSON(w,r,&cat)
