package adminendpoints

import (
	"database/sql"
	"time"
)

type ProductCreate struct {
	Name                 string  `json:"Product_Name"`
	Description          string  `json:"Product_Description"`
	Variations			 []VariationCreateWithProduct `json: "Variations"`
}

type VariationCreateWithProduct struct {
	Name              string  `json:"Variation_Name"`
	Description       string  `json:"Variation_Description"`
	Price             float32 `json:"Variation_Price"`
	PrimaryImage      string  `json:"Primary_Image,omitempty"`
	VariationQuantity int     `json:"Variation_Quantity"`
	LocationAt        string  `json:"Location_At"`
}


type ProductCreateRetrieve struct {
	ProductID  int64 `json:"Product_ID"`
	VarID      int64 `json:"Variation_ID"`
	ProdInvLoc int64 `json:"Inv_ID,omitempty"`
}

type ProductRetrieve struct {
	ProductID   int64  `json:"Product_ID"`
	Name        string `json:"Product_Name"`
	Description string `json:"Product_Description"`
}

type VariationRetrieve struct {
	Variation_ID int64          `json:"Variation_ID"`
	ProductID    int64          `json:"Product_ID"`
	Name         string         `json:"Variation_Name"`
	Description  string         `json:"Variation_Description"`
	Price        float32        `json:"Variation_Price"`
	PrimaryImage sql.NullString `json:"PRIMARY_IMAGE,omitempty"`
}

type VariationCreate struct {
	ProductID         int64   `json:"Product_ID"`
	Name              string  `json:"Variation_Name"`
	Description       string  `json:"Variation_Description"`
	Price             float32 `json:"Variation_Price"`
	PrimaryImage      string  `json:"Primary_Image,omitempty"`
	VariationQuantity int     `json:"Variation_Quantity"`
	LocationAt        string  `json:"Location_At"`
}

type ProdExist struct {
	ProductExists bool   `json:"Product_Exists"`
	Message       string `json:"Message"`
}

type variCrtd struct {
	VariationID    int64 `json:"Product_ID"`
	LocationExists bool  `json:"Location_Exists"`
}

type ProdInvLocCreation struct {
	VarID    int64  `json:"Variation_ID"`
	Quantity int    `json:"Quantity"`
	Location string `json:"Location"`
}
type PILCreated struct {
	InvID    int64  `json:"Inv_ID"`
	Quantity int    `json:"Quantity"`
	Location string `json:"Location"`
}

type CatToCat struct {
	CatStart int `json:"CategoryStart"`
	CatEnd   int `json:"CategoryEnd"`
}

type CatToProd struct {
	Cat  int `json:"Category"`
	Prod int `json:"Product"`
}

type ReadCat struct {
	Category int `json:"category"`
}

type ProductEdit struct {
	Name         *string `db:"Product_Name" json:"product_name,omitempty"`
	Description  *string `db:"Product_Description" json:"product_description,omitempty"`
	PrimaryImage *string `db:"PRIMARY_IMAGE" json:"primary_image,omitempty"`
}

type VariationEdit struct {
	VariationID          int64   `db:"Variation_ID" json:"variation_id"`
	VariationProductID   int64   `db:"Product_ID" json:"product_id"`
	VariationName        string  `db:"Variation_Name" json:"variation_name"`
	VariationDescription string  `db:"Variation_Description" json:"variation_description"`
	VariationPrice       float32 `db:"Variation_Price" json:"variation_price"`
	SKU                  string  `db:"SKU" json:"sku"`
	UPC                  string  `db:"UPC" json:"upc"`
	PrimaryImage         string  `db:"Primary_Image,omitempty" json:"primary_image,omitempty"`
	VariationQuantity    int     `db:"Variation_Quantity" json:"variation_quantity"`
	LocationAt           string  `db:"Location_At" json:"location_at"`
}

type DeletedSendBack struct {
	SendBack bool `json:"Deleted"`
}
type AddedSendBack struct {
	IDSendBack int64 `json:"AddedID"`
}
type Attribute struct {
	Attribute string `json:"attribute"`
}

type CategoryInsert struct{
	CategoryName string `json:"category_name"`
	CategoryDescription string `json:"category_description"`
}

type CategoryReturn struct{
	CategoryId int64 `json:"category_id"`
	CategoryName *string `json:"category_name"`
	CategoryDescription *string `json:"category_description"`
}

type CategoryTree struct{
  	Categories map[string]PrimeCategoryTree `json:"categories"`
}

type PrimeCategoryTree struct{
    PrimeCategoryID int `json:"prime_category_id"`
    PrimeCategoryName string `json:"prime_category_name"` 
    Categories map[string]SubCategoryTree `json:"categories"`
}

type SubCategoryTree struct{
    SubCategoryID int `json:"sub_category_id"`
    SubCategoryName string `json:"sub_category_name"`
    Categories map[string]FinalCategoryTree `json:"categories"`
}

type FinalCategoryTree struct{
    FinalCategoryID int `json:"final_category_id"`
    FinalCategoryName string `json:"final_name"`
}


type ProductSize struct {
    SizeID         *int64       `db:"Size_ID" json:"size_id"`
    SizeName       *string     `db:"Size_Name" json:"size_name"`
    SizeDescription *string    `db:"Size_Description,omitempty" json:"size_description,omitempty"`
    VariationID    *int64        `db:"Variation_ID" json:"variation_id"`
    VariationPrice *float64    `db:"Variation_Price" json:"variation_price"`
    SKU            *string    `db:"SKU,omitempty" json:"sku,omitempty"`
    UPC            *string    `db:"UPC,omitempty" json:"upc,omitempty"`
	Price 		*float64    `db:"Price,omitempty" json:"price,omitempty"`
    PrimaryImage   *string    `db:"PRIMARY_IMAGE,omitempty" json:"primary_image,omitempty"`
    DateCreated    *time.Time  `db:"Date_Created" json:"date_created"`
    ModifiedDate   *time.Time `db:"Modified_Date,omitempty" json:"modified_date,omitempty"`
}