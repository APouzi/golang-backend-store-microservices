package productendpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/APouzi/MerchantMachinee/routes/helpers"
	"github.com/go-chi/chi/v5"
)

type ProductRoutes struct{

}

func InstanceProductsRoutes( ) *ProductRoutes {
	r := &ProductRoutes{

	}
	return r
}

//Delete THIS AFTER FIX!!!!
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
var ProdJSON Product

func (route *ProductRoutes) GetAllProductsEndPoint(w http.ResponseWriter, r *http.Request) {
	// ProdJSON := route.ProductQuery.GetAllProducts(route.DB)
	ProdJSON = Product{Product_ID:55556}
	JSONWrite,err := json.Marshal(ProdJSON)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Failed")))
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(JSONWrite)
	
}


func (route *ProductRoutes) GetOneProductsEndPoint(w http.ResponseWriter, r *http.Request){
	_, err :=  strconv.Atoi(chi.URLParam(r,"ProductID"))
	if err != nil{
		fmt.Println("String to Int failed:", err)
	}
	// ProdJSON, err := route.ProductQuery.GetOneProduct(route.DB,query)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusAccepted,ProdJSON)

}

// func (route *Routes) GetProductCategoryEndPoint(w http.ResponseWriter, r *http.Request){
// 	category, err := strconv.Atoi(chi.URLParam(r, "CategoryName"))
	
// 	if err != nil{
// 		fmt.Println("Get Product Category ")
// 	}

// 	ProdJSON := route.ProductQuery.GetProductCategoryFinal(route.DB,category)
// 	JSONWrite, err := json.Marshal(ProdJSON)

// 	if err != nil{
// 		fmt.Println(err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(fmt.Sprint("Failed")))
// 	}

// 	w.WriteHeader((http.StatusAccepted))
// 	w.Write(JSONWrite)

// }


func (route *ProductRoutes) GetProductCategoryEndPointFinal(w http.ResponseWriter, r *http.Request){
	// category := chi.URLParam(r, "CategoryName")

	// if err != nil{
	// 	fmt.Println("Get Product Category ")
	// }
// TODO needs error handling for none existent categories!
	// ProdJSON := route.ProductQuery.GetProductCategoryFinal(route.DB,category)
	JSONWrite, err := json.Marshal(ProdJSON)

	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Failed")))
	}

	w.WriteHeader((http.StatusAccepted))
	w.Write(JSONWrite)

}
