package checkout

import "time"

type Config struct{
	STRIPE_KEY string
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


type FrontendRequest struct {
	Items []struct {
		Size_ID int64 `json:"size_id"`
		Quantity   int64  `json:"quantity"`
	} `json:"items"`

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