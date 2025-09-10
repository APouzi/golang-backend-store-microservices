package orders

import (
	"database/sql"
	"encoding/json"
	"time"
)

type JSON = json.RawMessage

type OrderRoutesTray struct {
	DB *sql.DB
}

// =============================
// Addresses & Customers
// =============================

type Address struct {
	AddressID  uint64     `db:"address_id"  json:"address_id"`
	FullName   *string    `db:"full_name"   json:"full_name,omitempty"`
	Line1      string     `db:"line1"       json:"line1"`
	Line2      *string    `db:"line2"       json:"line2,omitempty"`
	City       string     `db:"city"        json:"city"`
	// DB column is "region" (state/province). We expose it as "state" in JSON for your existing API.
	Region     string     `db:"region"      json:"state"`
	PostalCode string     `db:"postal_code" json:"postal_code"`
	Country    string     `db:"country"     json:"country"` // ISO-3166-1 alpha-2
	Phone      *string    `db:"phone"       json:"phone,omitempty"`
	CreatedAt  time.Time  `db:"created_at"  json:"created_at"`
}

type Customer struct {
	CustomerID uint64    `db:"customer_id" json:"customer_id"`
	Email      string    `db:"email"       json:"email"`
	FullName   *string   `db:"full_name"   json:"full_name,omitempty"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
}

type AddressInput struct {
	FullName   *string `json:"full_name,omitempty"`
	Line1      string  `json:"line1"`
	Line2      *string `json:"line2,omitempty"`
	City       string  `json:"city"`
	State      string  `json:"state"` // maps to DB column "region"
	PostalCode string  `json:"postal_code"`
	Country    string  `json:"country"` // ISO-3166-1 alpha-2
	Phone      *string `json:"phone,omitempty"`
}

// =============================
// Orders & Order Items
// =============================

type OrderStatus string

const (
	OrderStatusCreated         OrderStatus = "created"
	OrderStatusAwaitingPayment OrderStatus = "awaiting_payment"
	OrderStatusPaid            OrderStatus = "paid"
	OrderStatusFulfilled       OrderStatus = "fulfilled"
	OrderStatusCancelled       OrderStatus = "cancelled"
	OrderStatusClosed          OrderStatus = "closed"
)

type Order struct {
	OrderID           uint64       `db:"order_id"            json:"order_id"`
	OrderNumber       string       `db:"order_number"        json:"order_number"`
	CustomerID        *uint64      `db:"customer_id"         json:"customer_id,omitempty"`
	Email             string       `db:"email"               json:"email"`
	BillingAddressID  *uint64      `db:"billing_address_id"  json:"billing_address_id,omitempty"`
	BillingAddress    *AddressInput `db:"-"                   json:"billing_address,omitempty"`
	ShippingAddressID *uint64      `db:"shipping_address_id" json:"shipping_address_id,omitempty"`
	ShippingAddress   *AddressInput `db:"-"                   json:"shipping_address,omitempty"`
	Currency       string      `db:"currency"         json:"currency"`          // ISO-4217, e.g. "USD"
	SubtotalCents  int64       `db:"subtotal_cents"   json:"subtotal_cents"`    // minor units
	DiscountCents  int64       `db:"discount_cents"   json:"discount_cents"`
	ShippingCents  int64       `db:"shipping_cents"   json:"shipping_cents"`
	TaxCents       int64       `db:"tax_cents"        json:"tax_cents"`
	TotalCents     int64       `db:"total_cents"      json:"total_cents"`

	Status    OrderStatus `db:"status"     json:"status"`
	PlacedAt  time.Time   `db:"placed_at"  json:"placed_at"`

	Provider        *string `db:"provider"          json:"provider,omitempty"`           // "stripe", "paypal", etc.
	ProviderOrderID *string `db:"provider_order_id" json:"provider_order_id,omitempty"` // e.g., Stripe cs_...
	Metadata        JSON    `db:"metadata"          json:"metadata"`                    // arbitrary JSON

	CreatedAt time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt time.Time `db:"updated_at"  json:"updated_at"`
}

type OrderItem struct {
	OrderItemID       uint64           `db:"order_item_id"       json:"order_item_id"`
	OrderID           uint64           `db:"order_id"            json:"order_id"`
	ProductID         *uint64          `db:"product_id"          json:"product_id,omitempty"`
	SKU               *string          `db:"sku"                 json:"sku,omitempty"`
	Title             string           `db:"title"               json:"title"`


	Qty               int              `db:"qty"                 json:"qty"`
	Currency          string           `db:"currency"            json:"currency"` // ISO-4217

	UnitPriceCents    int64            `db:"unit_price_cents"    json:"unit_price_cents"`
	LineSubtotalCents int64            `db:"line_subtotal_cents" json:"line_subtotal_cents"`
	LineDiscountCents int64            `db:"line_discount_cents" json:"line_discount_cents"`
	LineTaxCents      int64            `db:"line_tax_cents"      json:"line_tax_cents"`
	LineTotalCents    int64            `db:"line_total_cents"    json:"line_total_cents"`
}

// =============================
// Payments & Refunds
// =============================

type PaymentStatus string

const (
	PaymentAuthorized       PaymentStatus = "authorized"
	PaymentCaptured         PaymentStatus = "captured"
	PaymentFailed           PaymentStatus = "failed"
	PaymentRefunded         PaymentStatus = "refunded"
	PaymentPartiallyRefunded PaymentStatus = "partially_refunded"
)

type Payment struct {
	PaymentID          uint64         `db:"payment_id"           json:"payment_id,omitempty"`
	OrderID            uint64         `db:"order_id"             json:"order_id, omitempty"`
	Provider           string         `db:"provider"             json:"provider"`               // "stripe"
	ProviderPaymentID  string         `db:"provider_payment_id"  json:"provider_payment_id"`    // Stripe pi_...
	MethodBrand        *string        `db:"method_brand"         json:"method_brand,omitempty"` // "visa"
	Last4              *string        `db:"last4"                json:"last4,omitempty"`
	Status             PaymentStatus  `db:"status"               json:"status"`
	AmountCents        int64          `db:"amount_cents"         json:"amount_cents"`
	Currency           string         `db:"currency"             json:"currency"`
	RawResponse        JSON           `db:"raw_response"         json:"raw_response,omitempty"`           // provider payload snapshot
	CreatedAt          time.Time      `db:"created_at"           json:"created_at"`
}

type Refund struct {
	RefundID         uint64    `db:"refund_id"          json:"refund_id"`
	PaymentID        uint64    `db:"payment_id"         json:"payment_id"`
	AmountCents      int64     `db:"amount_cents"       json:"amount_cents"`
	Reason           *string   `db:"reason"             json:"reason,omitempty"`
	ProviderRefundID *string   `db:"provider_refund_id" json:"provider_refund_id,omitempty"` // Stripe re_...
	CreatedAt        time.Time `db:"created_at"         json:"created_at"`
}

// =============================
// Promotions & Discounts (optional)
// =============================

type PromotionType string

const (
	PromoFixed   PromotionType = "fixed"
	PromoPercent PromotionType = "percent"
	PromoBogo    PromotionType = "bogo"
)

type Promotion struct {
	PromotionID uint64        `db:"promotion_id" json:"promotion_id"`
	Code        *string       `db:"code"         json:"code,omitempty"`
	Type        PromotionType `db:"type"         json:"type"`
	Value       float64       `db:"value"        json:"value"` // DECIMAL(12,4); use string/rat if you prefer exactness
	MaxUses     *int          `db:"max_uses"     json:"max_uses,omitempty"`
	StartsAt    *time.Time    `db:"starts_at"    json:"starts_at,omitempty"`
	EndsAt      *time.Time    `db:"ends_at"      json:"ends_at,omitempty"`
	Metadata    JSON          `db:"metadata"     json:"metadata"`
}

type OrderDiscount struct {
	OrderDiscountID uint64        `db:"order_discount_id" json:"order_discount_id"`
	OrderID         uint64        `db:"order_id"          json:"order_id"`
	PromotionID     *uint64       `db:"promotion_id"      json:"promotion_id,omitempty"`
	Code            *string       `db:"code"              json:"code,omitempty"`
	AmountCents     int64         `db:"amount_cents"      json:"amount_cents"`
	Allocation      JSON          `db:"allocation"        json:"allocation"` // optional per-line allocation snapshot
}

// =============================
// Shipments & Inventory (optional)
// =============================

type ShipmentStatus string

const (
	ShipmentPending   ShipmentStatus = "pending"
	ShipmentShipped   ShipmentStatus = "shipped"
	ShipmentDelivered ShipmentStatus = "delivered"
	ShipmentReturned  ShipmentStatus = "returned"
)

type Shipment struct {
	ShipmentID     uint64          `db:"shipment_id"    json:"shipment_id"`
	OrderID        uint64          `db:"order_id"       json:"order_id"`
	Carrier        *string         `db:"carrier"        json:"carrier,omitempty"`
	Service        *string         `db:"service"        json:"service,omitempty"`
	TrackingNumber *string         `db:"tracking_number" json:"tracking_number,omitempty"`
	ShippedAt      *time.Time      `db:"shipped_at"     json:"shipped_at,omitempty"`
	Status         ShipmentStatus  `db:"status"         json:"status"`
}

type ShipmentItem struct {
	ShipmentItemID uint64 `db:"shipment_item_id" json:"shipment_item_id"`
	ShipmentID     uint64 `db:"shipment_id"      json:"shipment_id"`
	OrderItemID    uint64 `db:"order_item_id"    json:"order_item_id"`
	Qty            int    `db:"qty"              json:"qty"`
}

type InventoryReason string

const (
	InvReserve   InventoryReason = "reserve"
	InvSale      InventoryReason = "sale"
	InvCancel    InventoryReason = "cancel"
	InvRefund    InventoryReason = "refund"
	InvAdjustment InventoryReason = "adjustment"
)

type InventoryRefType string

const (
	RefOrder  InventoryRefType = "order"
	RefRefund InventoryRefType = "refund"
	RefManual InventoryRefType = "manual"
)

type InventoryMovement struct {
	MovementID   uint64           `db:"movement_id"   json:"movement_id"`
	SKU          string           `db:"sku"           json:"sku"`
	VariationKey *string          `db:"variation_key" json:"variation_key,omitempty"`
	QtyDelta     int              `db:"qty_delta"     json:"qty_delta"`
	Reason       InventoryReason  `db:"reason"        json:"reason"`
	RefType      *InventoryRefType `db:"ref_type"     json:"ref_type,omitempty"`
	RefID        *uint64          `db:"ref_id"        json:"ref_id,omitempty"`
	CreatedAt    time.Time        `db:"created_at"    json:"created_at"`
}

type OrderItemsBulkReturn struct {
	CreatedIds []uint64    `json:"created_ids"`
	OrderItems []OrderItem `json:"order_items"`
}
// =============================
// Order Events (audit log)
// =============================

type OrderEvent struct {
	OrderEventID uint64    `db:"order_event_id" json:"order_event_id"`
	OrderID      uint64    `db:"order_id"       json:"order_id"`
	Type         string    `db:"type"           json:"type"`    // "webhook","status_changed","note",...
	Details      JSON      `db:"details"        json:"details"` // arbitrary JSON (e.g., stripe_event_id)
	CreatedAt    time.Time `db:"created_at"     json:"created_at"`
}
