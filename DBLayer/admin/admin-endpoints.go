package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

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




// func (adminRoutes *ProductRoutesTray) CreateInventory(w http.ResponseWriter, r *http.Request){
// 	prodInvLoc := ProdInvLocCreation{}
// 	helpers.ReadJSON(w,r,&prodInvLoc)
// }


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
