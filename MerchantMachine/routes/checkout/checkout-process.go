package checkout

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/APouzi/MerchantMachinee/routes/helpers"
	"github.com/stripe/stripe-go/v82"
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



func ProcessInventoryAndOrderItems( inventoryMap map[int64]int64, w http.ResponseWriter, GetProductInventoryDetailByID func(string, *[]InventoryProductDetail, http.ResponseWriter), GetProductSizeByID func(string, *ProductSize, http.ResponseWriter), UpdateInventoryShelfDetailQuantity func(int64, int64, http.ResponseWriter), ) ([]InventoryProductDetail, []OrderItem) {
	var listOfProductDetails []InventoryProductDetail
	var OrderLineItems []OrderItem

	for sizeID, quantity := range inventoryMap {
		fmt.Printf("  Size ID: %d → Quantity: %d\n", sizeID, quantity)
		var InvProdJSON []InventoryProductDetail
		strsizeID := strconv.FormatInt(sizeID, 10)
		GetProductInventoryDetailByID(strsizeID, &InvProdJSON, w)
		if len(InvProdJSON) == 0 {
			fmt.Printf("No inventory found for SizeID %s\n", strsizeID)
			continue
		}
		prdSize := &ProductSize{}
		GetProductSizeByID(strsizeID, prdSize, w)
		invProdJsonPull := InvProdJSON[0]
		new_quantity := invProdJsonPull.Quantity - quantity
		fmt.Println("old quantity is:", invProdJsonPull.Quantity, "new quantity of the of the is:", new_quantity)
		UpdateInventoryShelfDetailQuantity(invProdJsonPull.InventoryID, new_quantity, w)
		listOfProductDetails = append(listOfProductDetails, invProdJsonPull)

		unitPrice := *prdSize.VariationPrice * 100
		OrderLineItems = append(OrderLineItems, OrderItem{
			ProductID:         prdSize.SizeID,
			SKU:               prdSize.SKU,
			Title:             *prdSize.SizeName,
			Qty:               int(quantity),
			Currency:          "usd",
			UnitPriceCents:    int64(unitPrice),
			LineSubtotalCents: int64(unitPrice * float64(quantity)),
			LineDiscountCents: 0,
			LineTaxCents:      0,
			LineTotalCents:    int64(unitPrice * float64(quantity)),
		})
	}

	return listOfProductDetails, OrderLineItems
}




func ExtractInventoryMapFromSessionMetadata(metadata map[string]string) map[int64]int64 {
	inventoryMap := make(map[int64]int64)
	for key, value := range metadata {
		if strings.HasPrefix(key, "itemsizeqty_") {
			idPart := strings.TrimPrefix(key, "itemsizeqty_")
			if id, err := strconv.ParseInt(idPart, 10, 64); err == nil {
				if qty, err := strconv.ParseInt(value, 10, 64); err == nil {
					inventoryMap[id] = qty
				} else {
					fmt.Printf("⚠️ Skipping invalid quantity for %s: %s\n", key, value)
				}
			} else {
				fmt.Printf("⚠️ Skipping invalid item ID for %s\n", key)
			}
		}
	}
	return inventoryMap
}



func ExtractCustomerInfoFromSession(session *stripe.CheckoutSession) CustomerInfo {
	var cust CustomerInfo
	if session.CustomerDetails != nil {
		cust.Email = session.CustomerDetails.Email
		cust.Name = session.CustomerDetails.Name

		if session.CustomerDetails.Address != nil {
			cust.BillingAddress = &Address{
				Line1:      session.CustomerDetails.Address.Line1,
				Line2:      session.CustomerDetails.Address.Line2,
				City:       session.CustomerDetails.Address.City,
				State:      session.CustomerDetails.Address.State,
				PostalCode: session.CustomerDetails.Address.PostalCode,
				Country:    session.CustomerDetails.Address.Country,
			}
		}
	}
	return cust
}



func ExtractPaymentMethodDetails(pi *stripe.PaymentIntent, paymentMethodLast4, paymentMethodBrand string) (string, string) {
	if paymentMethodLast4 == "" && pi.LatestCharge != nil && pi.LatestCharge.PaymentMethodDetails != nil && pi.LatestCharge.PaymentMethodDetails.Card != nil {
		paymentMethodLast4 = pi.LatestCharge.PaymentMethodDetails.Card.Last4
		paymentMethodBrand = string(pi.LatestCharge.PaymentMethodDetails.Card.Brand)
	}
	return paymentMethodLast4, paymentMethodBrand
}

func ExtractBillingAddress(pi *stripe.PaymentIntent) *Address {
	if pi.PaymentMethod != nil && pi.PaymentMethod.BillingDetails != nil && pi.PaymentMethod.BillingDetails.Address != nil {
		return &Address{
			Line1:      pi.PaymentMethod.BillingDetails.Address.Line1,
			Line2:      pi.PaymentMethod.BillingDetails.Address.Line2,
			City:       pi.PaymentMethod.BillingDetails.Address.City,
			State:      pi.PaymentMethod.BillingDetails.Address.State,
			PostalCode: pi.PaymentMethod.BillingDetails.Address.PostalCode,
			Country:    pi.PaymentMethod.BillingDetails.Address.Country,
		}
	}
	return nil
}

func ExtractShippingAddress(pi *stripe.PaymentIntent) *Address {
	var shippingAddress *Address
	if pi.Shipping != nil && pi.Shipping.Address != nil {
		shippingAddress = &Address{
			Line1:      pi.Shipping.Address.Line1,
			Line2:      pi.Shipping.Address.Line2,
			City:       pi.Shipping.Address.City,
			State:      pi.Shipping.Address.State,
			PostalCode: pi.Shipping.Address.PostalCode,
			Country:    pi.Shipping.Address.Country,
		}
	}
	return shippingAddress
}