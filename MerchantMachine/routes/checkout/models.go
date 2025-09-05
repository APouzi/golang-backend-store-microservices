package checkout

import (
	"encoding/json"
	"time"

	"github.com/stripe/stripe-go/v82"
)

type Config struct{
	STRIPE_KEY string
    STRIPE_WEBHOOK_KEY string
}


type ProductJSONRetrieve struct {
	Product_ID          int    `json:"Product_ID"`
	Product_Name        string `json:"Product_Name"`
	Product_Description string `json:"Product_Description"`
	PRIMARY_IMAGE       string `json:"PRIMARY_IMAGE,omitempty"`
	ProductDateAdded    string `json:"DateAdded,omitempty"`
	ModifiedDate        string `json:"ModifiedDate,omitempty"`
}

type ProductResponse struct {
    Inventory struct {
        InvID      int    `json:"Inv_ID"`
        Quantity   int    `json:"Quantity"`
        LocationAt string `json:"LocationAt"`
    } `json:"inventory"`
    Product struct {
        ProductID          int     `json:"Product_ID"`
        ProductName        string  `json:"Product_Name"`
        ProductDescription string  `json:"Product_Description"`
        ProductPrice       float64 `json:"Product_Price"`
        SKU                string  `json:"SKU"`
        UPC                string  `json:"UPC"`
        PrimaryImage       string  `json:"PRIMARY_IMAGE"`
        ProductDateAdded   string  `json:"ProductDateAdded"`
        ModifiedDate       string  `json:"ModifiedDate"`
    } `json:"product"`
    Variation struct {
        VariationID int     `json:"variation_id"`
        Name        string  `json:"name"`
        Description string  `json:"description"`
        Price       float64 `json:"price"`
    } `json:"variation"`
}


type ProductVariation []struct {
	VariationID          int        `db:"Variation_ID" json:"variation_id"`
	ProductID            int        `db:"Product_ID" json:"product_id"`
	VariationName        string     `db:"Variation_Name" json:"variation_name"`
	VariationDescription string     `db:"Variation_Description,omitempty" json:"variation_description,omitempty"`
	VariationPrice       float64    `db:"Variation_Price" json:"variation_price"`
	SKU                  *string    `db:"SKU,omitempty" json:"sku,omitempty"`
	UPC                  *string    `db:"UPC,omitempty" json:"upc,omitempty"`
	PrimaryImage         *string    `db:"PRIMARY_IMAGE,omitempty" json:"primary_image,omitempty"`
	DateCreated          time.Time  `db:"Date_Created" json:"date_created"`
	ModifiedDate         *time.Time `db:"Modified_Date,omitempty" json:"modified_date,omitempty"`
}


type InventoryProductDetail struct {
    InventoryID int64  `json:"inventory_id" db:"inventory_id"`    // Primary Key
    Quantity    int64  `json:"quantity" db:"quantity"`            // NOT NULL
    ProductID   int64  `json:"product_id" db:"product_id"`        // Foreign Key to tblProductVariation
    LocationID  int64  `json:"location_id" db:"location_id"`      // Foreign Key to tblLocation
    Description string `json:"description,omitempty" db:"description"` // TEXT (nullable)
}

type OrderLineItem struct {
    ProductID  int64
    ProductSizeID int64
    VariationID int64
    Quantity   int64
    UnitPrice  int64
}

type QuantityUpdateResponse struct {
    Quantity    int64 `json:"quantity"`
}


type FrontendRequest struct {
	Items []struct {
		Size_ID int64 `json:"size_id"`
		Quantity   int64  `json:"quantity"`
	} `json:"items"`

}

type ProductSize struct {
    SizeID         *uint64       `db:"Size_ID" json:"size_id"`
    SizeName       *string     `db:"Size_Name" json:"size_name"`
    SizeDescription *string    `db:"Size_Description,omitempty" json:"size_description,omitempty"`
    VariationID    *int64        `db:"Variation_ID" json:"variation_id"`
    VariationPrice *float64    `db:"Variation_Price" json:"variation_price"`
    SKU            *string    `db:"SKU,omitempty" json:"sku,omitempty"`
    UPC            *string    `db:"UPC,omitempty" json:"upc,omitempty"`
    DateCreated    *time.Time  `db:"Date_Created" json:"date_created"`
    ModifiedDate   *time.Time `db:"Modified_Date,omitempty" json:"modified_date,omitempty"`
}




type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type CustomerInfo struct {
	Email          string   `json:"email"`
	Name           string   `json:"name"`
	BillingAddress *Address `json:"billing_address,omitempty"`
}

type OrderInfo struct {
	SessionID       string                                 `json:"session_id"`
	PaymentIntentID string                                 `json:"payment_intent_id"`
	PaymentStatus   stripe.CheckoutSessionPaymentStatus     `json:"payment_status"`
	Currency        string                                 `json:"currency"`
	AmountTotal     int64                                  `json:"amount_total"`
	CreatedAt       time.Time                              `json:"created_at"`
	Customer        CustomerInfo                           `json:"customer"`
}

type Order struct {
    ID              string    // internal order ID
    StripeSessionID string
    PaymentIntentID string
    CustomerEmail   string
    CustomerName    string
    BillingAddress  string
    ShippingAddress string
    TotalAmount     int64
    Currency        string
    Status          string
    CreatedAt       time.Time
    LineItems       []OrderLineItem
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

type OrderItemInput struct {
	ProductID       *uint64               `json:"product_id,omitempty"`
	SKU             *string               `json:"sku,omitempty"`
	Title           string                `json:"title"`
	Variation       map[string]any        `json:"variation,omitempty"`
	Qty             int                   `json:"qty"`
	UnitPriceCents  int64                 `json:"unit_price_cents"`
	// Currency is taken from the order; include here only if you truly need per-line currency.
}

type JSON = json.RawMessage
type OrderStatus string

const (
	OrderStatusCreated         OrderStatus = "created"
	OrderStatusAwaitingPayment OrderStatus = "awaiting_payment"
	OrderStatusPaid            OrderStatus = "paid"
	OrderStatusFulfilled       OrderStatus = "fulfilled"
	OrderStatusCancelled       OrderStatus = "cancelled"
	OrderStatusClosed          OrderStatus = "closed"
)

type OrderRecordCreation struct {
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

type PaymentStatus string

const (
	PaymentAuthorized       PaymentStatus = "authorized"
	PaymentCaptured         PaymentStatus = "captured"
	PaymentFailed           PaymentStatus = "failed"
	PaymentRefunded         PaymentStatus = "refunded"
	PaymentPartiallyRefunded PaymentStatus = "partially_refunded"
)

type PaymentCreation struct {
	PaymentID          uint64         `db:"payment_id"           json:"payment_id"`
	OrderID            uint64         `db:"order_id"             json:"order_id"`
	Provider           string         `db:"provider"             json:"provider"`               // "stripe"
	ProviderPaymentID  string         `db:"provider_payment_id"  json:"provider_payment_id"`    // Stripe pi_...
	MethodBrand        *string        `db:"method_brand"         json:"method_brand,omitempty"` // "visa"
	PaymentMethodID  *string        `db:"payment_method_id"    json:"payment_method_id,omitempty"`

	Last4              *string        `db:"last4"                json:"last4,omitempty"`
	Status             PaymentStatus  `db:"status"               json:"status"`
	AmountCents        int64          `db:"amount_cents"         json:"amount_cents"`
	Currency           string         `db:"currency"             json:"currency"`
	RawResponse        JSON           `db:"raw_response"         json:"raw_response"`           // provider payload snapshot
	CreatedAt          time.Time      `db:"created_at"           json:"created_at"`
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
	LineSubtotalCents int64            `db:"line_subtotal_cents" json:"line_subtotal_cents"` //from db
	LineDiscountCents int64            `db:"line_discount_cents" json:"line_discount_cents"`
	LineTaxCents      int64            `db:"line_tax_cents"      json:"line_tax_cents"`
	LineTotalCents    int64            `db:"line_total_cents"    json:"line_total_cents"` //subtotal - discount + tax
}

type ProductTaxCode struct {
	TaxCodeID          int    `json:"tax_code_id"`
	TaxCodeName        string `json:"tax_code_name"`
	TaxCodeDescription string `json:"tax_code_description"`
	TaxCode            string `json:"tax_code"`
    Provider           string `json:"provider"`
}


type ProductSizeTaxCode struct {
	SizeID    int `json:"size_id"`
	TaxCodeID int `json:"tax_code_id"`
}