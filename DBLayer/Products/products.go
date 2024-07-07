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

func GetProductRouteInstance() *ProductRoutes{
	return &ProductRoutes{}
}


func (route *ProductRoutes) GetAllProductsEndPoint(w http.ResponseWriter, r *http.Request) {
	ProdJSON := &Product{Product_ID: 55556}
	err := helpers.WriteJSON(w,200,ProdJSON)
	if err != nil {
		fmt.Println("GetAllProductsEndPoint failed",err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprint("Failed")))
	}
}