package checkout

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/APouzi/MerchantMachinee/routes/helpers"
)

func GetProductByID(prodID string, ProdJSON *ProductJSONRetrieve, w http.ResponseWriter) {
	url := "http://dblayer:8080/products/" + prodID
	resp, err := http.Get(url)
	jd := json.NewDecoder(resp.Body)
	jd.Decode(ProdJSON)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to decode response:", err)
	}
	if ProdJSON.Product_ID == 0 {
		helpers.ErrorJSON(w, errors.New("there was no response"), 404)
		return
	}
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to pull product:", err)
	}
}

func GetProductVariationByID(prodID string, ProdJSON *[]ProductResponse, w http.ResponseWriter) {
	url := "http://dblayer:8080/products/variations/" + prodID
	resp, err := http.Get(url)
	jd := json.NewDecoder(resp.Body)
	jd.Decode(ProdJSON)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to decode response:", err)
	}
	if len(*ProdJSON) == 0 {
		helpers.ErrorJSON(w, errors.New("there was no response"), 404)
		return
	}
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to pull product:", err)
	}
}

func GetProductSizeByID(productSizeID string, ProdJSON *ProductSize, w http.ResponseWriter) {
	url := "http://dblayer:8080/product-size/" + productSizeID
	resp, err := http.Get(url)
	jd := json.NewDecoder(resp.Body)
	jd.Decode(ProdJSON)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to decode response:", err)
	}
	// if len(*ProdJSON) == 0 {
	// 	helpers.ErrorJSON(w, errors.New("there was no response"), 404)
	// 	return
	// }
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to pull product:", err)
	}
}

func GetProductInventoryDetailByID(prodID string, InvProdJSON *[]InventoryProductDetail, w http.ResponseWriter) {
	url := "http://dblayer:8080/inventory/inventory-product-details/?product_size_id=" + prodID
	resp, err := http.Get(url)
	jd := json.NewDecoder(resp.Body)
	jd.Decode(InvProdJSON)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to decode response:", err)
	}
	if len(*InvProdJSON) == 0 {
		helpers.ErrorJSON(w, errors.New("there was no response"), 404)
		return
	}
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to pull product:", err)
	}
}