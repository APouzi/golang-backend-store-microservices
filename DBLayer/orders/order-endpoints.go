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

