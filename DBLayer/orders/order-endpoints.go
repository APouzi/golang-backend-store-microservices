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
