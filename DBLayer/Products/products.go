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
		helpers.ErrorJSON(w, errors.New("could not get variation"),404)
		return
	}
	if err != nil{
		fmt.Println("scanning error:",err)
	}
	
	helpers.WriteJSON(w,200,&variationJSON)

}


type VariationRet struct{
	Variation_ID sql.NullInt64
	ProductID sql.NullInt64
	Name sql.NullString
	Description sql.NullString 
	Price sql.NullFloat64 
	PrimaryImage sql.NullString 

}


func (prdRoutes *ProductRoutesTray) SearchProductsEndPoint(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query().Get("q")
	if query == ""{
		fmt.Println("no search term")
		helpers.ErrorJSON(w, errors.New("no search term"),404)
		return
	}
	fmt.Println(r.URL.Query().Get("q"))
	fmt.Println("query from product end point", query)
	sqlQuery := `
		SELECT p.Product_ID, p.Product_Name, p.Product_Description,
			   pv.Variation_ID, pv.Variation_Name, pv.Variation_Description, pv.Variation_Price,
			   pil.Inv_ID, pil.Quantity, pil.Location_At
		FROM tblProducts p
		LEFT JOIN tblProductVariation pv ON p.Product_ID = pv.Product_ID
		LEFT JOIN tblProductInventoryLocation pil ON pv.Variation_ID = pil.Variation_ID
		WHERE p.Product_Name LIKE ? OR pv.Variation_Name LIKE ?
	`

	//the question marks in the sql statement are replaced by the values of the query string, specifically the two parameters after the query string.
	rows, err := prdRoutes.DB.Query(sqlQuery, "%"+query+"%", "%"+query+"%")
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("database query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var product_list []map[string]interface{}
	for rows.Next() {
		var product Product
		var variation VariationRet
		var inventory struct {
			Inv_ID     sql.NullInt64
			Quantity  sql.NullInt64
			LocationAt sql.NullString
		}

		

		err := rows.Scan(
			&product.Product_ID, &product.Product_Name, &product.Product_Description,
			&variation.Variation_ID, &variation.Name, &variation.Description, &variation.Price,
			&inventory.Inv_ID, &inventory.Quantity, &inventory.LocationAt,
		)
		if err != nil {
			helpers.ErrorJSON(w, fmt.Errorf("error scanning row: %v", err), http.StatusInternalServerError)
			return
		}

		

		result := map[string]interface{}{
			"product": product,
		}
		if variation.Variation_ID.Valid{
			result["variation"] = variation
		}
		if inventory.Inv_ID.Valid {
			result["inventory"] = inventory
		}
		product_list = append(product_list, result)
	}

	if err = rows.Err(); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("error iterating rows: %v", err), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, product_list)
} 
