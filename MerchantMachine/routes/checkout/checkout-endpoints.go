package checkout

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/APouzi/MerchantMachinee/routes/helpers"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/webhook"
)


func NewStripeClient() *stripe.Client {
    return stripe.NewClient(stripe.Key)
}

type CheckoutRoutes struct{
	stripe *stripe.Client
	config Config
}

func InstanceCheckoutRoutes(stripe *stripe.Client, config Config) *CheckoutRoutes {
	r := &CheckoutRoutes{
		stripe: stripe,
		config: config,
	}
	return r
}




func(route *CheckoutRoutes) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req := FrontendRequest{}
	helpers.ReadJSON(w,r,&req)


	stripe.Key = route.config.STRIPE_KEY 

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("http://localhost:4200/success"), // URL to redirect to on success
		CancelURL:  stripe.String("http://localhost:4200/canceledorder"),  // URL to redirect to on cancellation
		LineItems:          []*stripe.CheckoutSessionLineItemParams{},
	}

	var ProdJSON *[]ProductResponse = &[]ProductResponse{}
	var ProdSizeJSON *ProductSize = &ProductSize{}
	var InvProdJSON []InventoryProductDetail = []InventoryProductDetail{}



	// Dynamically create line items from the request
	fmt.Println("req!!!!",req)
	for _, item := range req.Items {
		GetProductSizeByID(strconv.FormatInt(item.Size_ID,10), ProdSizeJSON, w)
		GetProductVariationByID(strconv.FormatInt(*ProdSizeJSON.VariationID,10), ProdJSON, w)
		GetProductInventoryDetailByID(strconv.FormatInt(item.Size_ID,10), &InvProdJSON, w)
		var quantityCount int64 = 0
		for _, val := range InvProdJSON {
			quantityCount += val.Quantity
		} 
		
		if quantityCount < item.Quantity || quantityCount == 0{
			// errorMsg := fmt.Sprintf("Not enough stock for %s. Only %d left.", *InvProdJSON.Produc, quantityCount)
			helpers.ErrorJSON(w, errors.New("insufficient inventory"), http.StatusBadRequest)
			return
		}
		
		params.LineItems = append(params.LineItems, &stripe.CheckoutSessionLineItemParams{
			
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
			Currency:   stripe.String("usd"),
			ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
				Name: stripe.String((*ProdJSON)[0].Product.ProductName),
			},
			UnitAmount: stripe.Int64(int64(math.Round(*ProdSizeJSON.VariationPrice * 100))),
			},
			Quantity: stripe.Int64(item.Quantity),
		})
		params.PaymentIntentData = &stripe.CheckoutSessionPaymentIntentDataParams{}
		params.PaymentIntentData.Metadata[fmt.Sprintf("itemsizeqty_%d", item.Size_ID)] = strconv.FormatInt(item.Quantity, 10)
		}

		

	s, err := session.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating checkout session: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the session ID to the frontend
	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w,200,map[string]string{"id": s.ID})
}


func(route *CheckoutRoutes) PaymentConfirmation(w http.ResponseWriter, r *http.Request){
	payload, _ := io.ReadAll(r.Body)
	fmt.Println("hello in payment confirm!")
  event, err := webhook.ConstructEventWithOptions(payload, r.Header.Get("Stripe-Signature"), route.config.STRIPE_WEBHOOK_KEY, webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true,})
  if err != nil {
	fmt.Println(err)
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
fmt.Println("hello in payment confirm! 2")
  switch event.Type {
  case "checkout.session.completed":
    fmt.Println("Hello—received checkout.session.completed")
    var session stripe.CheckoutSession
    json.Unmarshal(event.Data.Raw, &session)
    fmt.Println("This is completed!", &session)	
//   case "checkout.session.async_payment_succeeded":
//     // optional async case
// 	var session stripe.CheckoutSession
// 	fmt.Println("This is failed!", session)
//     err := json.Unmarshal(event.Data.Raw, &session)
//     if err != nil {
//         // handle err
// 	}
//     // session = event.Data.Object.(*stripe.CheckoutSession)
//     // handleCheckout(sess)
  default:
    log.Printf("Unhandled event: %s", event.Type)
  }

  w.WriteHeader(http.StatusOK)
}
