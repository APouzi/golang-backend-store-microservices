package orders

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/APouzi/DBLayer/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
)

// POST /orders
func (rt *OrderRoutesTray) CreateOrderRecord(w http.ResponseWriter, r *http.Request) {
	var ord Order
	err := helpers.ReadJSON(w, r, &ord)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Normalize + defaults (db also has defaults; we keep this predictable)
	ord.Email = strings.TrimSpace(ord.Email)
	ord.OrderNumber = strings.TrimSpace(ord.OrderNumber)
	if ord.Currency == "" {
		http.Error(w, "currency is required (ISO-4217, e.g., USD)", http.StatusBadRequest)
		return
	}
	ord.Currency = strings.ToUpper(strings.TrimSpace(ord.Currency))
	if ord.Status == "" {
		ord.Status = OrderStatusCreated
	}
	if ord.PlacedAt.IsZero() {
		ord.PlacedAt = time.Now().UTC()
	}
	// metadata column is NOT NULL; ensure "{}" at minimum
	if len(ord.Metadata) == 0 {
		ord.Metadata = json.RawMessage(`{}`)
	}

	// Basic validation
	if ord.OrderNumber == "" {
		http.Error(w, "order_number is required", http.StatusBadRequest)
		return
	}
	if ord.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	if ord.SubtotalCents < 0 || ord.DiscountCents < 0 || ord.ShippingCents < 0 || ord.TaxCents < 0 || ord.TotalCents < 0 {
		http.Error(w, "amounts must be >= 0", http.StatusBadRequest)
		return
	}

	// Hand off to DB layer
	id, err := InsertOrder(r.Context(), rt.DB, &ord)
	if err != nil {
		// Handle duplicate order_number nicely
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			http.Error(w, "order_number already exists", http.StatusConflict)
			return
		}
		log.Println("insert orders failed:", err)
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}
	ord.OrderID = id

	helpers.WriteJSON(w, http.StatusCreated, ord)
}



func InsertOrder(ctx context.Context, db *sql.DB, ord *Order) (uint64, error) {
	// Ensure placed_at is set even if caller forgot; DB has default too.
	if ord.PlacedAt.IsZero() {
		ord.PlacedAt = time.Now().UTC()
	}
	// Ensure metadata is non-null (column is NOT NULL)
	if len(ord.Metadata) == 0 {
		ord.Metadata = json.RawMessage(`{}`)
	}

	const q = `
INSERT INTO orders
(order_number, customer_id, email, billing_address_id, shipping_address_id,
 currency, subtotal_cents, discount_cents, shipping_cents, tax_cents, total_cents,
 status, placed_at, provider, provider_order_id, metadata)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	res, err := db.ExecContext(ctx, q,
		ord.OrderNumber,
		nullableUint64(ord.CustomerID),
		ord.Email,
		nullableUint64(ord.BillingAddressID),
		nullableUint64(ord.ShippingAddressID),
		ord.Currency,
		ord.SubtotalCents,
		ord.DiscountCents,
		ord.ShippingCents,
		ord.TaxCents,
		ord.TotalCents,
		ord.Status,
		ord.PlacedAt,
		nullableString(ord.Provider),
		nullableString(ord.ProviderOrderID),
		[]byte(ord.Metadata), // JSON NOT NULL
	)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Println("Inserted order id:", lastID)
	return uint64(lastID), nil
}


func nullableUint64(p *uint64) any {
	if p == nil {
		return nil
	}
	return *p
}
func nullableString(p *string) any {
	if p == nil {
		return nil
	}
	return *p
}

func (h *OrderRoutesTray) CreateOrderItemRecord(w http.ResponseWriter, r *http.Request) {



	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var in OrderItem
	err := helpers.ReadJSON(w, r, &in)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Minimal guards (kept simple; DB has FK + CHECK)
	if in.OrderID == 0 || in.Title == "" || in.Qty <= 0 || len(in.Currency) != 3 {
		http.Error(w, "missing/invalid required fields", http.StatusBadRequest)
		return
	}


	const q = `
	INSERT INTO order_items
	  (order_id, product_id, sku, title,
	   qty, currency, unit_price_cents, line_subtotal_cents,
	   line_discount_cents, line_tax_cents, line_total_cents)
	VALUES (?,?,?,?,?,?,?,?,?,?,?)`

	res, err := h.DB.Exec(
		q,
		in.OrderID,
		in.ProductID,
		in.SKU,
		in.Title,
		in.Qty,
		in.Currency,
		in.UnitPriceCents,
		in.LineSubtotalCents,
		in.LineDiscountCents,
		in.LineTaxCents,
		in.LineTotalCents,
	)
	if err != nil {
		// Surface FK/constraint issues plainly for now
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "failed to fetch insert id", http.StatusInternalServerError)
		return
	}
	in.OrderItemID = uint64(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(in)
}


func (h *OrderRoutesTray) CreateOrderItemRecordsBulk(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var orderItems []OrderItem
	err := helpers.ReadJSON(w, r, &orderItems)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	
	OrderItemsBulkReturn := OrderItemsBulkReturn{
		CreatedIds: []uint64{},
	}
	

	// Minimal guards (kept simple; DB has FK + CHECK)
	for _, orderItem := range orderItems {
		if orderItem.OrderID == 0 || orderItem.Title == "" || orderItem.Qty <= 0 || len(orderItem.Currency) != 3 {
			http.Error(w, "missing/invalid required fields", http.StatusBadRequest)
			return
		}

		const q = `
		INSERT INTO order_items
		(order_id, product_id, sku, title,
		qty, currency, unit_price_cents, line_subtotal_cents,
		line_discount_cents, line_tax_cents, line_total_cents)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)`

		res, err := h.DB.Exec(
			q,
			orderItem.OrderID,
			orderItem.ProductID,
			orderItem.SKU,
			orderItem.Title,
			orderItem.Qty,
			orderItem.Currency,
			orderItem.UnitPriceCents,
			orderItem.LineSubtotalCents,
			orderItem.LineDiscountCents,
			orderItem.LineTaxCents,
			orderItem.LineTotalCents,
		)
		if err != nil {
			// Surface FK/constraint issues plainly for now
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "failed to fetch insert id", http.StatusInternalServerError)
			return
		}
		orderItem.OrderItemID = uint64(id)

		OrderItemsBulkReturn.CreatedIds = append(OrderItemsBulkReturn.CreatedIds, orderItem.OrderItemID)
		OrderItemsBulkReturn.OrderItems = append(OrderItemsBulkReturn.OrderItems, orderItem)
		
	}
	


	helpers.WriteJSON(w, http.StatusCreated, OrderItemsBulkReturn)
}



func (h *OrderRoutesTray) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var p Payment
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	const q = `
	INSERT INTO payments
	  (order_id, provider, provider_payment_id, method_brand, last4,
	   status, amount_cents, currency, raw_response)
	VALUES (?,?,?,?,?,?,?,?,?)`

	res, err := h.DB.Exec(
		q,
		p.OrderID, p.Provider, p.ProviderPaymentID, p.MethodBrand, p.Last4,
		p.Status, p.AmountCents, p.Currency, p.RawResponse,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, _ := res.LastInsertId()   
	p.PaymentID = uint64(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(p)
}

func (h *OrderRoutesTray) GetPayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "paymentID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	const q = `
	SELECT payment_id, order_id, provider, provider_payment_id, method_brand, last4,
	       status, amount_cents, currency, raw_response, created_at
	FROM payments WHERE payment_id = ?`

	var p Payment
	if err := h.DB.QueryRow(q, id).Scan(
		&p.PaymentID, &p.OrderID, &p.Provider, &p.ProviderPaymentID,
		&p.MethodBrand, &p.Last4, &p.Status, &p.AmountCents,
		&p.Currency, &p.RawResponse, &p.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(p)
}


func (h *OrderRoutesTray) CreateRefund(w http.ResponseWriter, r *http.Request) {
	var ref Refund
	if err := json.NewDecoder(r.Body).Decode(&ref); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	const q = `
	INSERT INTO refunds
	  (payment_id, amount_cents, reason, provider_refund_id)
	VALUES (?,?,?,?)`

	res, err := h.DB.Exec(
		q,
		ref.PaymentID, ref.AmountCents, ref.Reason, ref.ProviderRefundID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, _ := res.LastInsertId()
	ref.RefundID = uint64(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ref)
}

func (h *OrderRoutesTray) GetRefund(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "refundID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	const q = `
	SELECT refund_id, payment_id, amount_cents, reason, provider_refund_id, created_at
	FROM refunds WHERE refund_id = ?`

	var ref Refund
	if err := h.DB.QueryRow(q, id).Scan(
		&ref.RefundID, &ref.PaymentID, &ref.AmountCents,
		&ref.Reason, &ref.ProviderRefundID, &ref.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ref)
}


