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

func GetProductVariationByID(prodID string, ProdJSON *ProductVariation, w http.ResponseWriter) {
	url := "http://dblayer:8080/products/variation/" + prodID
	resp, err := http.Get(url)
	jd := json.NewDecoder(resp.Body)
	jd.Decode(ProdJSON)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to decode response:", err)
	}
	if ProdJSON.VariationID == 0 {
		helpers.ErrorJSON(w, errors.New("there was no response"), 404)
		return
	}
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		fmt.Println("failed to pull product:", err)
	}
}