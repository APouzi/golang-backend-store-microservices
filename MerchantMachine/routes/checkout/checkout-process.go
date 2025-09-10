package checkout

import (
	"bytes"
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


func UpdateInventoryShelfDetailQuantity(shelfID int64, newQuantity int64, w http.ResponseWriter) {
	url := fmt.Sprintf("http://dblayer:8080/inventory/inventory-product-details/%d", shelfID)
	fmt.Println("Updating inventory shelf detail at URL:", url)
	qur := QuantityUpdateResponse{
		Quantity: newQuantity,
	}
	jsonData, err := json.Marshal(qur)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		return
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		helpers.ErrorJSON(w, err, 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		helpers.ErrorJSON(w, fmt.Errorf("failed to update inventory shelf detail: %s", resp.Status), resp.StatusCode)
		return
	}
}	

func GetProductTaxCodeByID(sizeID string, TaxCodeJSON *[]ProductTaxCode, r *http.Request ) {
	url := "http://dblayer:8080/tax-codes-intermediary/" + sizeID
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("failed to fetch tax code:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}


	if err := json.NewDecoder(resp.Body).Decode(TaxCodeJSON); err != nil {
		fmt.Println("failed to decode tax code response:", err)
		return
	}
}