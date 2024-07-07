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


}