package products

import (
	"fmt"
	"net/http"

	"github.com/APouzi/DBLayer/helpers"
)

type Product struct {
	Product_ID          int
	Product_Name        string
	Product_Description string
	Product_Price       float32
	SKU                 string
	UPC                 string
	PRIMARY_IMAGE       string
	ProductDateAdded    string
	ModifiedDate        string
}

type ProductRoutes struct{

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