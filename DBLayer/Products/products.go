package products

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/APouzi/DBLayer/helpers"
	"github.com/go-chi/chi/v5"
)

type ProductRoutesTray struct{
	DB *sql.DB
	getAllProductsStmt *sql.Stmt
	getOneProductStmt *sql.Stmt
	getOneVariationProductStment *sql.Stmt
}

func GetProductRouteInstance(dbInst *sql.DB) *ProductRoutesTray{
	
	routeMap := prepareProductRoutes(dbInst)
	return &ProductRoutesTray{
		getAllProductsStmt: routeMap["getAllProducts"],
		getOneProductStmt: routeMap["getOneProducts"],
		getOneVariationProductStment: routeMap["getOneVariationProducts"],
		DB: dbInst,
	}
}

func prepareProductRoutes(dbInst *sql.DB) map[string]*sql.Stmt{
	sqlStmentsMap := make(map[string]*sql.Stmt)

	getAllPrdStment, err := dbInst.Prepare("SELECT Product_ID, Product_Name, Product_Description FROM tblProducts")
	if err != nil{
		log.Fatal(err)
	}
	GetOneProductStmt, err := dbInst.Prepare("SELECT Product_ID, Product_Name, Product_Description, PRIMARY_IMAGE, Date_Created, Modified_Date FROM tblProducts WHERE Product_ID = ?")
	if err != nil{
		log.Fatal(err)
	}
	GetOneVariationStmt, err := dbInst.Prepare("SELECT Variation_ID, Product_ID, Variation_Name, Variation_Description, Variation_Price, PRIMARY_IMAGE FROM tblProductVariation WHERE Variation_ID = ?")
	if err != nil{
		log.Fatal(err)
	}
	if err != nil{
		fmt.Println("failed to create sql statements")
	}
	
	sqlStmentsMap["getAllProducts"] = getAllPrdStment
	sqlStmentsMap["getOneProducts"] = GetOneProductStmt
	sqlStmentsMap["getOneVariationProducts"] = GetOneVariationStmt

	return sqlStmentsMap
}


func (prdRoutes *ProductRoutesTray) GetAllProductsEndPoint(w http.ResponseWriter, r *http.Request) {

	rows, _ := prdRoutes.getAllProductsStmt.Query()
	ListProducts := []ProductJSONRetrieve{}
	prodJSON := ProductJSONRetrieve{}
	defer rows.Close()
	for rows.Next(){
		
		err := rows.Scan(
			&prodJSON.Product_ID,
			&prodJSON.Product_Name,
			&prodJSON.Product_Description,
		)

		if err != nil{
			fmt.Println("Scanning Error:",err)
		}
		ListProducts = append(ListProducts, prodJSON)
	}

	helpers.WriteJSON(w,200,ListProducts)

}

type ProductJSON struct {
	Product_ID          int     `json:"Product_ID"`
	Product_Name        string  `json:"Product_Name"`
	Product_Description string  `json:"Product_Description"`
	PRIMARY_IMAGE       string  `json:"PRIMARY_IMAGE,omitempty"`
	ProductDateAdded   string  `json:"DateAdded"`
	ModifiedDate       string `json:"ModifiedDate"`
}

func (prdRoutes *ProductRoutesTray) GetOneProductEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit getoneproductendpoint")
	productID, err :=  strconv.Atoi(chi.URLParam(r,"ProductID"))
	if err != nil{
		fmt.Println("String to Int failed:", err)
	}
	rows :=prdRoutes.getOneProductStmt.QueryRow(productID)
	prodJSON := ProductJSON{}
	
	err = rows.Scan(
		&prodJSON.Product_ID, 
		&prodJSON.Product_Name, 
		&prodJSON.Product_Description,  
		&prodJSON.PRIMARY_IMAGE,
		&prodJSON.ProductDateAdded,
		&prodJSON.ModifiedDate,
	)
	if err == sql.ErrNoRows{
		helpers.WriteJSON(w,404,"Cant find product son")
	}
	if err != nil{
		fmt.Println("scanning error:",err)
	}
	
	helpers.WriteJSON(w,200,prodJSON)

}

type VariationRetrieve struct{
	Variation_ID int64 `json:"Variation_ID"`
	ProductID int64 `json:"Product_ID"`
	Name string `json:"Variation_Name"`
	Description string `json:"Variation_Description"`
	Price float32 `json:"Variation_Price"`
	PrimaryImage sql.NullString `json:"PRIMARY_IMAGE,omitempty"`

}

type FailedDBQuery struct {
	Msg string `json:"message"`
}

func (prdRoutes *ProductRoutesTray) GetOneVariationEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit getoneproductendpoint")
	productVariationID, err :=  strconv.Atoi(chi.URLParam(r,"VariationID"))
	if err != nil{
		fmt.Println("String to Int failed:", err)
	}
	fmt.Println("hit getoneproductendpoint", productVariationID)
	// row := route.DB.QueryRow("SELECT Variation_ID FROM tblProductVariation WHERE Variation_ID = ?",pil.VarID)
	rows := prdRoutes.getOneVariationProductStment.QueryRow(productVariationID)
	variationJSON := &VariationRetrieve{}
	
	err = rows.Scan(
		&variationJSON.Variation_ID,
		&variationJSON.ProductID, 
		&variationJSON.Name, 
		&variationJSON.Description,  
		&variationJSON.Price,
		&variationJSON.PrimaryImage,
	)
	if err == sql.ErrNoRows{
		helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed, could not retrieve varitation id"),404)
		return
	}
	if err != nil{
		fmt.Println("scanning error:",err)
	}
	
	helpers.WriteJSON(w,200,&variationJSON)

}

