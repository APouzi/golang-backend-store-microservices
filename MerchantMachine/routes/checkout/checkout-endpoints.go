package checkout

import (
	"bytes"
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
	"github.com/rs/xid"
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
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{},
		},
		AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{
			Enabled: stripe.Bool(true),
		},
		ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
			AllowedCountries: stripe.StringSlice([]string{"US", "CA"}),
		},
		ShippingOptions: []*stripe.CheckoutSessionShippingOptionParams{
			{
				ShippingRate: stripe.String("shr_1S3t2jBAREvrwjtYED1pr2Zh"), 
			},
		},
	}


	ProdJSON := &[]ProductResponse{}
	ProdSizeJSON := &ProductSize{}
	var InvProdJSON []InventoryProductDetail

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			helpers.ErrorJSON(w, fmt.Errorf("invalid quantity for size_id %d", item.Size_ID), http.StatusBadRequest)
			return
		}
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

		taxCodes := &[]ProductTaxCode{}
		taxCode := ""
		GetProductTaxCodeByID(strconv.FormatInt(item.Size_ID, 10), taxCodes, r)
		for _, v := range *taxCodes {
			if v.Provider != "stripe" {
				taxCode = "txcd_32050025"
			}else{
				taxCode = v.Provider
			}
		}

		params.LineItems = append(params.LineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String((*ProdJSON)[0].Product.ProductName),
					TaxCode: stripe.String(taxCode), // Cosmetics Beautifying
					Metadata: map[string]string{
						"product_id": strconv.FormatInt(int64((*ProdJSON)[0].Product.ProductID), 10),
						"size_id":    strconv.FormatInt(int64(item.Size_ID), 10),
					},
				},
				TaxBehavior: stripe.String("exclusive"),
				UnitAmount: stripe.Int64(int64(math.Round(*ProdSizeJSON.VariationPrice * 100))), // in cents
			},
			Quantity: stripe.Int64(item.Quantity),
		})

		key := fmt.Sprintf("itemsizeqty_%d", item.Size_ID)
		params.PaymentIntentData.Metadata[key] = strconv.FormatInt(item.Quantity, 10)
	}
	

	s, err := session.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating checkout session: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"id": s.ID}, )
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
		checkoutSessionID := session.ID
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
		checkoutSessionID := session.ID
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
		listOfProductDetails := []InventoryProductDetail{}
		OrderLineItems := []OrderItem{}
		for sizeID, quantity := range inventoryMap {
			fmt.Printf("  Size ID: %d → Quantity: %d\n", sizeID, quantity)
			http.Get("http:")
			var InvProdJSON []InventoryProductDetail = []InventoryProductDetail{}
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
			UpdateInventoryShelfDetailQuantity(invProdJsonPull.InventoryID, new_quantity, w)
			// Send to the order service
			listOfProductDetails = append(listOfProductDetails, invProdJsonPull)

			unitPrice := *prdSize.VariationPrice * 100
			OrderLineItems = append(OrderLineItems, OrderItem{
				ProductID:   prdSize.SizeID,
				SKU:         prdSize.SKU,
				Title:       *prdSize.SizeName,
				Qty:         int(quantity),
				Currency:    "usd",
				UnitPriceCents: int64(unitPrice),
				LineSubtotalCents: int64(unitPrice * float64(quantity)),
				LineDiscountCents: 0,
				LineTaxCents:      0,
				LineTotalCents:    int64(unitPrice * float64(quantity)),
			})
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

		var paymentIntent *stripe.PaymentIntent = session.PaymentIntent
		if paymentIntent != nil {
			fmt.Printf("Stripe Payment Intent ID: %s\n", paymentIntent.ID)
			// You can now use paymentIntent.ID for further processing, e.g., storing in your DB
		}
		var paymentMethodBrand, paymentMethodID, paymentMethodLast4, customerID, chargeID string 

		if session.Customer != nil {
			customerID = session.Customer.ID
		}
		pi, err := route.fetchPIBits(paymentIntentID)
		if err != nil {
			log.Printf("Error fetching PaymentIntent details: %v", err)
			break
		}
		if pi.LatestCharge != nil {
			chargeID = pi.LatestCharge.ID
		}

		// Payment method details via PaymentMethod
		if pi.PaymentMethod != nil {
			paymentMethodID = pi.PaymentMethod.ID
			if pi.PaymentMethod.Card != nil {
				paymentMethodBrand = string(pi.PaymentMethod.Card.Brand)
				paymentMethodLast4 = pi.PaymentMethod.Card.Last4
			}
		}

		// (Alternative) read last4/brand from the Charge’s payment_method_details:
		if paymentMethodLast4 == "" && pi.LatestCharge != nil &&
		pi.LatestCharge.PaymentMethodDetails != nil &&
		pi.LatestCharge.PaymentMethodDetails.Card != nil {
			paymentMethodLast4 = pi.LatestCharge.PaymentMethodDetails.Card.Last4
			paymentMethodBrand = string(pi.LatestCharge.PaymentMethodDetails.Card.Brand)
		}

		var billingAddress *Address
		if pi.PaymentMethod.BillingDetails != nil && pi.PaymentMethod.BillingDetails.Address != nil {
			billingAddress = &Address{
				Line1:      pi.PaymentMethod.BillingDetails.Address.Line1,
				Line2:      pi.PaymentMethod.BillingDetails.Address.Line2,
				City:       pi.PaymentMethod.BillingDetails.Address.City,
				State:      pi.PaymentMethod.BillingDetails.Address.State,
				PostalCode: pi.PaymentMethod.BillingDetails.Address.PostalCode,
				Country:    pi.PaymentMethod.BillingDetails.Address.Country,
			}
		}

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

		
		// Safely extract shipping and tax cents (guard against nil pointers)
		shippingCents := int64(0)
		if session.ShippingCost != nil {
			shippingCents = session.ShippingCost.AmountSubtotal
		}
		taxCents := int64(0)
		if session.TotalDetails != nil {
			taxCents = session.TotalDetails.AmountTax
			fmt.Println("taxCents:", taxCents)
		}

		// Generate a unique order number: <RANDOM>-<DDMMYYYY>
		orderNumber := GenerateOrderNumber()

		orderSendOff := OrderRecordCreation{
			OrderID: 	   0, // to be filled by order service
			OrderNumber:   orderNumber,
			CustomerID:    nil, // to be filled by order service if applicable
			Email: 	   cust.Email,
			BillingAddressID: nil, // to be filled by order service if applicable
			ShippingAddressID: nil, // to be filled by order service if applicable,
			Currency:       string(session.Currency),
			SubtotalCents:  session.AmountSubtotal, // Assuming no tax/shipping for simplicity
			DiscountCents:  0,
			ShippingCents:  shippingCents,
			TaxCents:       taxCents,
			TotalCents:     session.AmountTotal,
			Status:         "created",
			PlacedAt: 	   time.Unix(session.Created, 0),
			Provider:       stripe.String("stripe"),
			ProviderOrderID: &checkoutSessionID,
		}

		payment_payload := PaymentCreation{
			Provider:          "stripe",
			ProviderPaymentID: paymentIntentID,
			MethodBrand:       &paymentMethodBrand,
			PaymentMethodID:   &paymentMethodID,
			Last4:             &paymentMethodLast4,
			Status:            PaymentCaptured,
			AmountCents:       session.AmountTotal,
			Currency:          string(session.Currency),
			RawResponse:       JSON([]byte{}), // Optionally store raw JSON response from Stripe
			CreatedAt:         time.Now(),
		}

		// Combine order, line items, and payment into a single payload
		type OrderPayload struct {
			Order      OrderRecordCreation `json:"order"`
			LineItems  []OrderItem         `json:"line_items"`
			Payment    PaymentCreation     `json:"payment"`
		}

		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(orderSendOff); err != nil {
			log.Printf("Failed to encode payload: %v", err)
			http.Error(w, "Failed to encode order payload", http.StatusInternalServerError)
			return
		}

		resp, err := http.Post("http://dblayer:8080/summary-order", "application/json", &buf)
		if err != nil {
			log.Printf("Failed to POST order: %v", err)
			http.Error(w, "Failed to send order to service", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		fmt.Println("customer id is:", customerID, "shipping address:", shippingAddress, "billing address:", billingAddress)
		fmt.Printf("Stripe Charge ID: %s\n", chargeID)
		fmt.Println("order send off:", orderSendOff, listOfProductDetails, payment_payload)
	default:
		log.Printf("Unhandled event: %s", event.Type)
  }

  helpers.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}


func(route *CheckoutRoutes) fetchPIBits(piID string) (*stripe.PaymentIntent, error) {
	stripe.Key = route.config.STRIPE_KEY
    params := &stripe.PaymentIntentParams{
        Expand: []*string{
            stripe.String("payment_method"),              // pm_..., includes card.brand/last4
            stripe.String("latest_charge"),               // ch_...
            // Optionally dig into charge’s PM details instead:
            stripe.String("latest_charge.payment_method_details"),
        },
    }
    return paymentintent.Get(piID, params)
}

func GenerateOrderNumber() string {
	id := xid.New().String()[:6] // short unique id (first 6 chars)
	return fmt.Sprintf("%s-%s", id, time.Now().Format("020106")) // ddmmyy
}