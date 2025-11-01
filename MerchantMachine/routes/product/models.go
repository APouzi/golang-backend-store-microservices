package productendpoints

type ProductRetrieve struct {
	ProductID   int64  `json:"Product_ID"`
	Name        string `json:"Product_Name"`
	Description string `json:"Product_Description:omitempty"`
}

type ProductSize struct {
	SizeID          *int64   `json:"size_id,omitempty"`
	SizeName        *string  `json:"size_name,omitempty"`
	SizeDescription *string  `json:"size_description,omitempty"`
	VariationPrice  *float64 `json:"variation_price,omitempty"`
	SKU             *string  `json:"sku,omitempty"`
	UPC             *string  `json:"upc,omitempty"`
}

type Variation struct {
	Name        *string               `json:"name,omitempty"`
	Description *string               `json:"description,omitempty"`
	ProductSize map[int64]ProductSize `json:"product_size,omitempty"`
}

type ProductWrapper struct {
	ProductID          int64               `json:"product_id"`
	ProductName        string              `json:"product_name"`
	ProductDescription string              `json:"product_description"`
	Variations         map[int64]Variation `json:"variations,omitempty"`
}