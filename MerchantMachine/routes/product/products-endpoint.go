package productendpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/APouzi/MerchantMachinee/routes/helpers"
	"github.com/go-chi/chi/v5"
)

type ProductRoutes struct{
	DBBaseURL string
}

func InstanceProductsRoutes(dbBaseURL string) *ProductRoutes {
	r := &ProductRoutes{
		DBBaseURL: dbBaseURL,
	}
	return r
}


func (route *ProductRoutes) GetAllProductsAndVariationsEndPoint(w http.ResponseWriter, r *http.Request) {
	var ProdJSON *[]ProductRetrieve = &[]ProductRetrieve{}
	// ProdJSON := route.ProductQuery.GetAllProducts(route.DB)
	reqURL := fmt.Sprintf("%s/products", route.DBBaseURL)
	resp, err := http.Get(reqURL)
	if err != nil{
		helpers.ErrorJSON(w, fmt.Errorf("failed db pull: %v", err), http.StatusInternalServerError)
		return
	}
	// fmt.Println(ret.Body)
	body, err := io.ReadAll(resp.Body)
	fmt.Println("Response Body:", string(body))
	if err != nil{
		fmt.Println("Read all Failed", err)
	}
	defer resp.Body.Close()
	
	err = json.Unmarshal(body, ProdJSON)
	if err != nil{	
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	if err := helpers.WriteJSON(w,200,ProdJSON); err != nil {
		fmt.Println("WriteJSON error:", err)
	}
}

func (route *ProductRoutes) GetAllProductsEndPoint(w http.ResponseWriter, r *http.Request) {
	var ProdJSON *[]ProductRetrieve = &[]ProductRetrieve{}
	// ProdJSON := route.ProductQuery.GetAllProducts(route.DB)
	reqURL := fmt.Sprintf("%s/products", route.DBBaseURL)
	resp, err := http.Get(reqURL)
	if err != nil{
		helpers.ErrorJSON(w, fmt.Errorf("failed db pull: %v", err), http.StatusInternalServerError)
		return
	}
	// fmt.Println(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("Read all Failed", err)
	}
	defer resp.Body.Close()
	
	err = json.Unmarshal(body, ProdJSON)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	if err := helpers.WriteJSON(w,200,ProdJSON); err != nil {
		fmt.Println("WriteJSON error:", err)
	}
}




func (route *ProductRoutes) GetOneProductsEndPoint(w http.ResponseWriter, r *http.Request){
	var ProdJSON *ProductRetrieve = &ProductRetrieve{}
	prodID :=  chi.URLParam(r,"ProductID")
	
	reqURL := fmt.Sprintf("%s/products/%s", route.DBBaseURL, prodID)
	resp, err := http.Get(reqURL)
	if err != nil{
		helpers.ErrorJSON(w, fmt.Errorf("failed db pull: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	jd := json.NewDecoder(resp.Body)
	if err := jd.Decode(ProdJSON); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to decode response: %v", err), http.StatusInternalServerError)
		return
	}
	if ProdJSON.ProductID == 0{
		helpers.ErrorJSON(w, errors.New("there was no response"), 404)
		return
	}
	// ProdJSON, err := route.ProductQuery.GetOneProduct(route.DB,query)
	// No additional error checks are required here; decoding errors already handled.
	if err := helpers.WriteJSON(w, http.StatusAccepted,ProdJSON); err != nil {
		fmt.Println("WriteJSON error:", err)
	}

}

func (route *ProductRoutes) SearchProductsEndPoint(w http.ResponseWriter, r *http.Request){
	searchQuery := r.URL.Query().Get("q")
	if searchQuery == "" {
		helpers.ErrorJSON(w, errors.New("search query is required"), http.StatusBadRequest)
		return
	}

	fmt.Println("searchQuery", searchQuery)
	escapedQuery := url.QueryEscape(searchQuery)
	reqURL := fmt.Sprintf("%s/products/search?q=%s", route.DBBaseURL, escapedQuery)

	resp, err := http.Get(reqURL)
	if err != nil{
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to fetch product from DBLayer:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		helpers.ErrorJSON(w, fmt.Errorf("DBLayer service returned status: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var searchResults []ProductWrapper
	err = json.NewDecoder(resp.Body).Decode(&searchResults)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to decode search results: %v", err), http.StatusInternalServerError)
		return
	}

	if err := helpers.WriteJSON(w, http.StatusOK, searchResults); err != nil {
		fmt.Println("WriteJSON error:", err)
	}
}

// GetProductAndVariationsPaginated proxies a pagination request to the DBLayer
func (route *ProductRoutes) GetProductAndVariationsPaginated(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("page_size")
	if page == "" {
		page = "1"
	}
	if pageSize == "" {
		pageSize = "10"
	}
	reqURL := fmt.Sprintf("%s/products/variations/pagination/?page=%s&page_size=%s", route.DBBaseURL, page, pageSize)

	resp, err := http.Get(reqURL)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to call DBLayer: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		helpers.ErrorJSON(w, fmt.Errorf("DBLayer returned status %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	// Decode proxied response body and stream to the client
	var payload interface{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to decode DBLayer response: %v", err), http.StatusInternalServerError)
		return
	}

	if err := helpers.WriteJSON(w, http.StatusOK, payload); err != nil {
		fmt.Println("WriteJSON error:", err)
	}
}

// GetProductAndVariationsByProductID proxies request for a product's variations to DBLayer
func (route *ProductRoutes) GetProductAndVariationsByProductID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	if productID == "" {
		helpers.ErrorJSON(w, fmt.Errorf("productID path param required"), http.StatusBadRequest)
		return
	}
	reqURL := fmt.Sprintf("%s/products/variations/%s", route.DBBaseURL, productID)
	resp, err := http.Get(reqURL)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to call DBLayer: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		helpers.ErrorJSON(w, fmt.Errorf("DBLayer returned status %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}
	var payload []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to decode DBLayer response: %v", err), http.StatusInternalServerError)
		return
	}
	if err := helpers.WriteJSON(w, http.StatusOK, payload); err != nil {
		fmt.Println("WriteJSON error:", err)
	}
}



// func (route *Routes) GetProductCategoryEndPoint(w http.ResponseWriter, r *http.Request){

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
