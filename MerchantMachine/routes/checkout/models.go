package checkout

import "time"
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

type FrontendRequest struct {
	Items []struct {
		Product_ID int64 `json:"Product_ID"`
		Quantity   int64  `json:"quantity"`
	} `json:"items"`

}