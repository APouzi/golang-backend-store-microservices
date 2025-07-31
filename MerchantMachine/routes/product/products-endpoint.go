package productendpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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


func (route *ProductRoutes) GetAllProductsAndVariationsEndPoint(w http.ResponseWriter, r *http.Request) {
	var ProdJSON *[]ProductRetrieve = &[]ProductRetrieve{}
	// ProdJSON := route.ProductQuery.GetAllProducts(route.DB)
	ret, err := http.Get("http://dblayer:8080/products")
	if err != nil{
		fmt.Println("failed db pull", err)
	}
	// fmt.Println(ret.Body)
	body, err := io.ReadAll(ret.Body)
	fmt.Println("Response Body:", string(body))
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

func (route *ProductRoutes) GetAllProductsEndPoint(w http.ResponseWriter, r *http.Request) {
	var ProdJSON *[]ProductRetrieve = &[]ProductRetrieve{}
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
	var ProdJSON *ProductRetrieve = &ProductRetrieve{}
	prodID :=  chi.URLParam(r,"ProductID")
	
	url := "http://dblayer:8080/products/" + prodID
	resp, err := http.Get(url)
	jd := json.NewDecoder(resp.Body)
	jd.Decode(ProdJSON)
	if err != nil{
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to decode response:", err)
	}
	if ProdJSON.ProductID == 0{
		helpers.ErrorJSON(w, errors.New("there was no response"), 404)
		return
	}
	if err != nil{
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to pull product:", err)
	}
	// ProdJSON, err := route.ProductQuery.GetOneProduct(route.DB,query)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusAccepted,ProdJSON)

}

func (route *ProductRoutes) SearchProductsEndPoint(w http.ResponseWriter, r *http.Request){
	searchQuery := r.URL.Query().Get("q")
	if searchQuery == "" {
		helpers.ErrorJSON(w, errors.New("search query is required"), http.StatusBadRequest)
		return
	}

	fmt.Println("searchQuery", searchQuery)
	url := fmt.Sprintf("http://dblayer:8080/search/?q=%s", searchQuery)

	resp, err := http.Get(url)

	if resp.StatusCode != http.StatusOK {
		helpers.ErrorJSON(w, fmt.Errorf("DBLayer service returned status: %d", err), http.StatusInternalServerError)
		return
	}

	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to search products: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	

	var searchResults []ProductWrapper
	err = json.NewDecoder(resp.Body).Decode(&searchResults)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to decode search results: %v", err), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, searchResults)
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


// func (route *ProductRoutes) GetProductCategoryEndPointFinal(w http.ResponseWriter, r *http.Request){
// 	var ProdJSON *Product = &Product{}
// 	// category := chi.URLParam(r, "CategoryName")

// 	// if err != nil{
// 	// 	fmt.Println("Get Product Category ")
// 	// }
// // TODO needs error handling for none existent categories!
// 	// ProdJSON := route.ProductQuery.GetProductCategoryFinal(route.DB,category)
// 	JSONWrite, err := json.Marshal(ProdJSON)

// 	if err != nil{
// 		fmt.Println(err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(fmt.Sprint("Failed")))
// 	}

// 	w.WriteHeader((http.StatusAccepted))
// 	w.Write(JSONWrite)

// }
