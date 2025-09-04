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
	"time"

	"github.com/APouzi/MerchantMachinee/routes/helpers"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/paymentintent"
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




func (route *CheckoutRoutes) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating checkout session...")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FrontendRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid request body: %w", err), http.StatusBadRequest)
		return
	}
	if len(req.Items) == 0 {
		helpers.ErrorJSON(w, errors.New("no items"), http.StatusBadRequest)
		return
	}

	stripe.Key = route.config.STRIPE_KEY

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String("http://localhost:4200/success"),
		CancelURL:          stripe.String("http://localhost:4200/canceledorder"),
		LineItems:          []*stripe.CheckoutSessionLineItemParams{},
		// Initialize PaymentIntentData once, with a non-nil Metadata map
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{},
		},
		AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{
			Enabled: stripe.Bool(true),
		},
	}

	fmt.Println("req!!!!", req)

	// Reusable holders
	ProdJSON := &[]ProductResponse{}
	ProdSizeJSON := &ProductSize{}
	var InvProdJSON []InventoryProductDetail

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			helpers.ErrorJSON(w, fmt.Errorf("invalid quantity for size_id %d", item.Size_ID), http.StatusBadRequest)
			return
		}

		// fetch size, product, and inventory
		GetProductSizeByID(strconv.FormatInt(item.Size_ID, 10), ProdSizeJSON, w)
		if ProdSizeJSON.VariationID == nil {
			helpers.ErrorJSON(w, fmt.Errorf("size %d has no variation", item.Size_ID), http.StatusBadRequest)
			return
		}

		GetProductVariationByID(strconv.FormatInt(*ProdSizeJSON.VariationID, 10), ProdJSON, w)
		if len(*ProdJSON) == 0 {
			helpers.ErrorJSON(w, fmt.Errorf("variation %d not found", *ProdSizeJSON.VariationID), http.StatusBadRequest)
			return
		}

		InvProdJSON = InvProdJSON[:0] // reuse slice
		GetProductInventoryDetailByID(strconv.FormatInt(item.Size_ID, 10), &InvProdJSON, w)

		var quantityCount int64
		for _, v := range InvProdJSON {
			quantityCount += v.Quantity
		}
		if quantityCount <= 0 || quantityCount < item.Quantity {
			helpers.ErrorJSON(w, errors.New("insufficient inventory"), http.StatusBadRequest)
			return
		}

		// Add line item
		params.LineItems = append(params.LineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String((*ProdJSON)[0].Product.ProductName),
				},
				UnitAmount: stripe.Int64(int64(math.Round(*ProdSizeJSON.VariationPrice * 100))), // in cents
			},
			Quantity: stripe.Int64(item.Quantity),
		})

		// Safe now: Metadata is initialized
		key := fmt.Sprintf("itemsizeqty_%d", item.Size_ID)
		params.PaymentIntentData.Metadata[key] = strconv.FormatInt(item.Quantity, 10)
	}
	fmt.Printf("Stripe Checkout Session Params: %+v\n", params)
	

	s, err := session.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating checkout session: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"id": s.ID})
}



func(route *CheckoutRoutes) PaymentConfirmation(w http.ResponseWriter, r *http.Request){
	payload, _ := io.ReadAll(r.Body)
  event, err := webhook.ConstructEventWithOptions(payload, r.Header.Get("Stripe-Signature"), route.config.STRIPE_WEBHOOK_KEY, webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true,})
  if err != nil {
	fmt.Println(err)
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  switch event.Type {
  case "checkout.session.completed":
    var session stripe.CheckoutSession
    json.Unmarshal(event.Data.Raw, &session)
  if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		inventoryMap := make(map[int64]int64) // map[SizeID]Quantity

		for key, value := range session.Metadata {
			if strings.HasPrefix(key, "itemsizeqty_"){
				// Extract the item ID from the key: item_<id>_qty
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
		// productsOrdered := []OrderLineItem{}
		for sizeID, quantity := range inventoryMap {
			fmt.Printf("  Size ID: %d → Quantity: %d\n", sizeID, quantity)
			http.Get("http:")
			var InvProdJSON []InventoryProductDetail = []InventoryProductDetail{}
			strsizeID := strconv.FormatInt(sizeID, 10)
			if err != nil{
				fmt.Println("oh no!")
				return
			}
			GetProductInventoryDetailByID(strsizeID, &InvProdJSON, w)
			if len(InvProdJSON) == 0 {
				fmt.Printf("No inventory found for SizeID %s\n", strsizeID)
				continue
			}
			new_quantity := InvProdJSON[0].Quantity - quantity
			fmt.Println("old quantity is:", InvProdJSON[0].Quantity, "new quantity of the of the is:", new_quantity)
			UpdateInventoryShelfDetailQuantity(InvProdJSON[0].InventoryID, new_quantity, w)
			// Send to the order service
		
		
			// productsOrdered = append(productsOrdered, OrderLineItem{
			// 	Size_ID:  sizeID,
			// 	Quantity: quantity,
			// 	UnitPrice: InvProdJSON[0].ProductID,
			// })

		}
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

		// ---- Extract Order fields ----
		paymentIntentID := ""
		if session.PaymentIntent != nil {
			// If not expanded, Stripe still fills the ID field.
			paymentIntentID = session.PaymentIntent.ID
		}

		order := Order{
			StripeSessionID:       session.ID,
			PaymentIntentID: paymentIntentID,
			Status:   string(session.PaymentStatus),
			Currency:        string(session.Currency), // Currency is a stripe.Currency type (alias of string)
			TotalAmount:     session.AmountTotal,
			CreatedAt:       time.Unix(session.Created, 0),
			// CustomerName:        cust,
			LineItems:       []OrderLineItem{},
		}


		fmt.Println("order for:",order)
  default:
    log.Printf("Unhandled event: %s", event.Type)
  }

  w.WriteHeader(http.StatusOK)
}
