package checkout

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/APouzi/MerchantMachinee/routes/helpers"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)



var productDatabase = map[int64]struct {
	Name  string
	Price int64 // Price in cents
}{
	4: {"T-Shirt", 2000},
	5: {"Mug", 1500},
}

// FrontendRequest represents the structure of the JSON sent from the Angular app


func NewStripeClient() *stripe.Client {
    return stripe.NewClient(stripe.Key)
}

type CheckoutRoutes struct{
	stripe *stripe.Client
}

func InstanceCheckoutRoutes(stripe *stripe.Client) *CheckoutRoutes {
	r := &CheckoutRoutes{
		stripe: stripe,
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
	

	stripe.Key = ""

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("http://localhost:4200/success"), // URL to redirect to on success
		CancelURL:  stripe.String("http://localhost:4200/cancel"),  // URL to redirect to on cancellation
		LineItems:          []*stripe.CheckoutSessionLineItemParams{},
	}

	var ProdJSON *ProductVariation = &ProductVariation{}



	// Dynamically create line items from the request
	fmt.Println("req!!!!",req)
	for _, item := range req.Items {
		// GetProductByID(strconv.FormatInt(item.Product_ID,10), ProdJSON, w)
		GetProductVariationByID(strconv.FormatInt(item.Product_ID,10), ProdJSON, w)
		params.LineItems = append(params.LineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
			Currency:   stripe.String("usd"),
			ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
				Name: stripe.String(ProdJSON.VariationName),
			},
			//If we want to charge 19.99 then we have to do 1999 for stripe
			UnitAmount: stripe.Int64(int64(math.Round(ProdJSON.VariationPrice * 100))),
			},
			Quantity: stripe.Int64(item.Quantity),
		})
		}

	s, err := session.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating checkout session: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the session ID to the frontend
	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w,200,map[string]string{"id": s.ID})
	// json.NewEncoder(w).Encode(map[string]string{"id": s.ID})
}
