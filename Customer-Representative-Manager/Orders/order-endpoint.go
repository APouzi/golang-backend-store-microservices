package orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/apouzi/customer-representative-manager/helpers"
	"github.com/jung-kurt/gofpdf"
)

type Product struct {
	Product_ID   string  `json:"product_id"`
	Product_Name string  `json:"product_name"`
	Quantity     string  `json:"quantity"`
	Price        float64 `json:"price"`
	Discount     string  `json: "discount"`
}

type OrderSummary struct {
	// Left Column Details
	ShippingAddress string    `json:"shipping_address"`
	CustomerName    string    `json:"customer_name"`
	PaymentMethod   string    `json:"payment_method"`
	OrderName       string    `json:"order_name"`
	OrderDateTime   string    `json:"order_date_time"`
	ProductList     []Product `json:"product_list"`

	// Right Column Details
	TotalAmount    float64 `json:"total_amount"`
	Taxes          float64 `json:"taxes"`
	ShippingMethod string  `json:"shipping_method"`
	ShippingCost   float64 `json:"shipping_cost"`
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hey!")

	order_sum := &OrderSummary{ProductList: []Product{}}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(order_sum)
	fmt.Println(order_sum)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	addHeader(pdf)
	drawBorder(pdf, 30)
	addOrderTable(pdf, order_sum)
	addOrderSummary(pdf)
	addFooter(pdf)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `inline; filename="order_summary.pdf"`)
	if err := pdf.Output(w); err != nil {
		http.Error(w, "failed to create pdf", http.StatusInternalServerError)
		return
	}

}


type CreateOrderUsingModels struct {
	Customer        *Customer `json:"customer"`                    // use Email (+ optional FullName)
	BillingAddress  *Address  `json:"billing_address,omitempty"`   // uses json:"state" -> Address.Region
	ShippingAddress *Address  `json:"shipping_address,omitempty"`
	Order           Order     `json:"order"`                       // partial: number, email, currency, discounts, shipping, tax, provider...
	Items           []OrderItem `json:"items"`                     // qty, unit_price_cents, sku/title/variation, (optional) line-level tax/discount
}

// Response we’ll return (same types, with computed fields filled)
type NormalizedOrder struct {
	Customer        *Customer   `json:"customer,omitempty"`
	BillingAddress  *Address    `json:"billing_address,omitempty"`
	ShippingAddress *Address    `json:"shipping_address,omitempty"`
	Order           Order       `json:"order"`
	Items           []OrderItem `json:"items"`
}

func CreateOrderSummaryRecord(w http.ResponseWriter, r *http.Request) {
	var requisition CreateOrderUsingModels
	err := helpers.ReadJSON(w,r,&requisition)
	if err != nil {
		helpers.ErrorJSON(w,err)
		return
	}
	if err := validateReq(requisition); err != nil {
		helpers.ErrorJSON(w,err)
		return
	}

	now := time.Now().UTC()

	// Normalize email/currency casing, timestamps, and defaults
	requisition.Order.Email = strings.TrimSpace(requisition.Order.Email)
	requisition.Order.Currency = strings.ToUpper(strings.TrimSpace(requisition.Order.Currency))
	if requisition.Order.Status == "" {
		requisition.Order.Status = OrderStatusCreated
	}
	if requisition.Order.PlacedAt.IsZero() {
		requisition.Order.PlacedAt = now
	}
	requisition.Order.CreatedAt = now
	requisition.Order.UpdatedAt = now

	// Normalize addresses (Region from json:"state" already populates Address.Region)
	if requisition.BillingAddress != nil {
		requisition.BillingAddress.Country = strings.ToUpper(requisition.BillingAddress.Country)
		requisition.BillingAddress.CreatedAt = now
	}
	if requisition.ShippingAddress != nil {
		requisition.ShippingAddress.Country = strings.ToUpper(requisition.ShippingAddress.Country)
		requisition.ShippingAddress.CreatedAt = now
	}

	// Build/compute items
	items := make([]OrderItem, 0, len(requisition.Items))
	var subtotal int64
	weights := make([]int64, len(requisition.Items)) // for proration if needed

	for i, it := range requisition.Items {
		// Ensure currency on lines (inherit from order if empty)
		if strings.TrimSpace(it.Currency) == "" {
			it.Currency = requisition.Order.Currency
		} else {
			it.Currency = strings.ToUpper(it.Currency)
		}
		if it.Currency != requisition.Order.Currency {
			http.Error(w, fmt.Sprintf("items[%d].currency must match order.currency", i), http.StatusBadRequest)
			return
		}
		if it.Qty <= 0 {
			http.Error(w, fmt.Sprintf("items[%d].qty must be > 0", i), http.StatusBadRequest)
			return
		}
		if it.UnitPriceCents < 0 {
			http.Error(w, fmt.Sprintf("items[%d].unit_price_cents must be >= 0", i), http.StatusBadRequest)
			return
		}

		it.LineSubtotalCents = int64(it.Qty) * it.UnitPriceCents
		// Keep any caller-provided line discounts/tax; if zero, we may prorate below.
		if it.LineDiscountCents < 0 || it.LineTaxCents < 0 {
			http.Error(w, fmt.Sprintf("items[%d] line amounts must be >= 0", i), http.StatusBadRequest)
			return
		}

		subtotal += it.LineSubtotalCents
		weights[i] = it.LineSubtotalCents
		items = append(items, it)
	}

	// If caller didn’t specify line-level values, prorate order-level discount/tax
	var sumLineDisc, sumLineTax int64
	for _, it := range items {
		sumLineDisc += it.LineDiscountCents
		sumLineTax += it.LineTaxCents
	}

	if sumLineDisc == 0 && requisition.Order.DiscountCents > 0 {
		alloc := prorate(requisition.Order.DiscountCents, weights)
		for i := range items {
			items[i].LineDiscountCents = alloc[i]
		}
	}
	if sumLineTax == 0 && requisition.Order.TaxCents > 0 {
		alloc := prorate(requisition.Order.TaxCents, weights)
		for i := range items {
			items[i].LineTaxCents = alloc[i]
		}
	}

	// Finish line totals
	for i := range items {
		items[i].LineTotalCents = items[i].LineSubtotalCents - items[i].LineDiscountCents + items[i].LineTaxCents
		// Ensure line currency still matches
		items[i].Currency = requisition.Order.Currency
	}

	// Fill order header amounts
	requisition.Order.SubtotalCents = subtotal
	requisition.Order.TotalCents = subtotal - requisition.Order.DiscountCents + requisition.Order.ShippingCents + requisition.Order.TaxCents

	// Return normalized result (ready for your insert logic)
	resp := NormalizedOrder{
		Customer:        requisition.Customer,
		BillingAddress:  requisition.BillingAddress,
		ShippingAddress: requisition.ShippingAddress,
		Order:           requisition.Order,
		Items:           items,
	}

	helpers.WriteJSON(w,http.StatusOK,resp)   
}



func validateReq(req CreateOrderUsingModels) error {
	if req.Customer == nil || strings.TrimSpace(req.Customer.Email) == "" {
		return errors.New("customer.email is required")
	}
	if strings.TrimSpace(req.Order.OrderNumber) == "" {
		return errors.New("order.order_number is required")
	}
	if strings.TrimSpace(req.Order.Email) == "" {
		return errors.New("order.email is required")
	}
	if strings.TrimSpace(req.Order.Currency) == "" {
		return errors.New("order.currency is required (ISO-4217)")
	}
	if len(req.Items) == 0 {
		return errors.New("items must be non-empty")
	}
	if req.Order.DiscountCents < 0 || req.Order.ShippingCents < 0 || req.Order.TaxCents < 0 {
		return errors.New("order discount_cents/shipping_cents/tax_cents must be >= 0")
	}
	return nil
}

// prorate splits 'amount' across weights so sum(out) == amount (round-half-up).
// This is useful for allocating order-level discounts/taxes across line items
// in a way that matches Stripe and other payment APIs, which require line-level
// amounts to sum exactly to the total. For example, if you have a $5 discount
// and 3 items with weights [100, 200, 700], prorate will distribute the $5
// proportionally (rounded) so the sum is $5, avoiding rounding errors.
func prorate(amount int64, weights []int64) []int64 {
	out := make([]int64, len(weights))
	if amount == 0 || len(weights) == 0 {
		return out
	}
	var total int64
	for _, w := range weights {
		total += w
	}
	if total == 0 {
		out[0] = amount
		return out
	}
	var allocated int64
	maxIdx, maxWeight := 0, int64(0)
	for i, weight := range weights {
		share := 0.0
		if weight > 0 {
			share = float64(amount) * (float64(weight) / float64(total))
		}
		part := int64(math.Round(share))
		out[i] = part
		allocated += part
		if weight > maxWeight {
			maxWeight = weight
			maxIdx = i
		}
	}
	diff := amount - allocated
	if diff != 0 {
		out[maxIdx] += diff
	}
	return out
}
