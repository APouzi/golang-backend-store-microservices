package products

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/APouzi/DBLayer/helpers"
)

type ProductRoutes struct{
	getAllProductsStmt *sql.Stmt
}

func GetProductRouteInstance(dbInst *sql.DB) *ProductRoutes{
	
	routeMap := prepareProductRoutes(dbInst)
	return &ProductRoutes{
		getAllProductsStmt: routeMap["getAllProducts"],
	}
}

func prepareProductRoutes(dbInst *sql.DB) map[string]*sql.Stmt{
	sqlStmentsMap := make(map[string]*sql.Stmt)

	getAllPrdStment, err := dbInst.Prepare("SELECT Product_ID, Product_Name, Product_Description FROM tblProducts")
	if err != nil{
		fmt.Println("failed to create sql statements")
	}
	
	sqlStmentsMap["getAllProducts"] = getAllPrdStment


	return sqlStmentsMap
}


func (prdRoutes *ProductRoutes) GetAllProductsEndPoint(w http.ResponseWriter, r *http.Request) {

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

func (prdRoutes *ProductRoutes) InsertIntoTblProducts(w http.ResponseWriter, r *http.Request){
	// tRes, err := transaction.Exec("INSERT INTO tblProducts(Product_Name, Product_Description) VALUES(?,?)", productRetrieve.Name,productRetrieve.Description)
	// if err != nil{
	// 	fmt.Println("transaction at tblProduct has failed")
	// 	fmt.Println(err)
	// 	transaction.Rollback()
	// }
}


func (prdRoutes *ProductRoutes) InsertIntoTblProductVariation(w http.ResponseWriter, r *http.Request){
	// tRes, err = transaction.Exec("INSERT INTO tblProductVariation(Product_ID,Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)",prodID, productRetrieve.VariationName, productRetrieve.VariationDescription, productRetrieve.VariationPrice)
	// if err != nil{
	// 	fmt.Println("transaction at tblProductVariation has failed")
	// 	fmt.Println(err)
	// 	transaction.Rollback()
	// 	return
	// }
}


func (prdRoutes *ProductRoutes) InsertIntoTblProductInventoryLocation(w http.ResponseWriter, r *http.Request){
	// tRes, err = transaction.Exec("INSERT INTO tblProductInventoryLocation(Variation_ID, Quantity, Location_AT) VALUES(?,?,?)",  ProdVarID,productRetrieve.VariationQuantity, productRetrieve.LocationAt)
	// if err != nil {
	// 	fmt.Println("transaction at tblProductInventory has failed")
	// 	fmt.Println(err)
	// }

}