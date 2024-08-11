package productendpoints

import (
	"encoding/json"
	"fmt"
	"io"
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


func (route *ProductRoutes) GetAllProductsEndPoint(w http.ResponseWriter, r *http.Request) {
	var ProdJSON *[]Product = &[]Product{}
	// ProdJSON := route.ProductQuery.GetAllProducts(route.DB)
	ret, err := http.Get("http://dblayer:8080/db/products/")
	if err != nil{
		fmt.Println("failed db pull", err)
	}
	// fmt.Println(ret.Body)
	body, err := io.ReadAll(ret.Body)
	if err != nil{
		fmt.Println("Read all Failed", err)
	}
	defer ret.Body.Close()
	
	err = json.Unmarshal(body, ProdJSON)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	helpers.WriteJSON(w,200,ProdJSON)
}


func (route *ProductRoutes) GetOneProductsEndPoint(w http.ResponseWriter, r *http.Request){
	var ProdJSON *Product = &Product{}
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
	var ProdJSON *Product = &Product{}
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
