package products

import (
	"database/sql"
	"time"
)

type Product struct {
	Product_ID          int
	Product_Name        string
	Product_Description string
	Product_Price       float32
	SKU                 string
	UPC                 string
	PRIMARY_IMAGE       string
	ProductDateAdded    string
	ModifiedDate        string
}

type ProductJSONRetrieve struct {
	Product_ID          int    `json:"Product_ID"`
	Product_Name        string `json:"Product_Name"`
	Product_Description string `json:"Product_Description"`
	PRIMARY_IMAGE       string `json:"PRIMARY_IMAGE,omitempty"`
	ProductDateAdded    string `json:"DateAdded,omitempty"`
	ModifiedDate        string `json:"ModifiedDate,omitempty"`
}

type ProductVariation struct {
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

type Category struct {
	Name string
	Test string `json:"test,omitempty"`
}

type Inventory struct {
}
type rowData struct {
    ProductID          int64
    ProductName        string
    ProductDescription string
    VariationID        *int64
    VariationName      *string
    VariationDesc      *string
    VariationPrice     *float64
    InvID              *int64
    Quantity           *int64
    LocationAt         *string
    SizeID             *int64
    SizeName           *string
    SizeDesc           *string
    SizeVariationID    *int64
    SizeVariationPrice sql.NullFloat64
    SKU                *string
    UPC                *string
    PrimaryImage       *string
}

type ProductSize struct {
    SizeID         *int64       `db:"Size_ID" json:"size_id"`
    SizeName       *string     `db:"Size_Name" json:"size_name"`
    SizeDescription *string    `db:"Size_Description,omitempty" json:"size_description,omitempty"`
    VariationID    *int64        `db:"Variation_ID" json:"variation_id"`
    VariationPrice *float64    `db:"Variation_Price" json:"variation_price"`
    SKU            *string    `db:"SKU,omitempty" json:"sku,omitempty"`
    UPC            *string    `db:"UPC,omitempty" json:"upc,omitempty"`
         *string    `db:"PRIMARY_IMAGE,omitempty" json:"primary_image,omitempty"`
    DateCreated    *time.Time  `db:"Date_Created" json:"date_created"`
    ModifiedDate   *time.Time `db:"Modified_Date,omitempty" json:"modified_date,omitempty"`
}

type ProductPaginated struct {
    ProductID          int64   `json:"product_id"`
    ProductName        string  `json:"product_name"`
    ProductDescription string  `json:"product_description"`
	Variations map[int64]Variation `json:"variations"`
}

type Variation struct {
    VariationID   *int64    `json:"variation_id"`
    Name          *string   `json:"name"`
    Description   *string `json:"description,omitempty"`
	ProductSize map[int64]ProductSize `json:"product_size"`
}

type InventoryPaginated struct {
    InvID      *int64  `json:"inv_id,omitempty"`
    Quantity   *int64  `json:"quantity,omitempty"`
    LocationAt *string `json:"location_at,omitempty"`
}
type ProductVariationPaginated struct {
    Product   ProductPaginated    `json:"product"`
    Variation []Variation  `json:"variation"`
    Inventory InventoryPaginated  `json:"inventory"`
	
}

type PaginatedResponse struct {
    Data       []ProductPaginated `json:"data"`
    Page       int                `json:"page"`
    PageSize   int                `json:"page_size"`
    TotalCount int                `json:"total_count"`
    TotalPages int                `json:"total_pages"`
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
}	TaxCodeID          int    `json:"tax_code_id"`
	TaxCodeName        string `json:"tax_code_name"`
	TaxCodeDescription string `json:"tax_code_description"`
	TaxCode            string `json:"tax_code"`
    Provider           string `json:"provider"`
}

type ProductSizeTaxCode struct {
	SizeID    int `json:"size_id"`
	TaxCodeID int `json:"tax_code_id"`
}