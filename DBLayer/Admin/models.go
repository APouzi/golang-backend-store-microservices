package admin

import "time"

// Product automatically creates Variation
type ProductCreate struct{
	Name string `json:"Product_Name"`
	Description string `json:"Product_Description"`
	Variations []VariationCreateChain `json:"Variations"`
}


type VariationCreateChain struct{
	Name string `json:"Variation_Name"`
	Description string `json:"Variation_Description"`
	Price float32 `json:"Variation_Price"`
	PrimaryImage string `json:"Primary_Image,omitempty"`
	VariationQuantity int  `json:"Variation_Quantity"`
	LocationAt string `json:"Location_At"`
}


type ProductCreateRetrieve struct{
	ProductID int64 `json:"Product_ID"`
	VarID int64 `json:"Variation_ID"`
	ProdInvLoc int64 `json:"Inv_ID,omitempty"`

}

type ProductRetrieve struct{
	ProductID int64 `json:"Product_ID"`
	Name string `json:"Product_Name"`
	Description string `json:"Product_Description"`
}


type VariationRetrieve struct{
	Variation_ID int64 `json:"Variation_ID"`
	ProductID int64 `json:"Product_ID"`
	Name string `json:"Variation_Name"`
	Description string `json:"Variation_Description"`
	Price float32 `json:"Variation_Price"`
	PrimaryImage string `json:"PRIMARY_IMAGE,omitempty"`

}

type VariationCreate struct{
	ProductID int64 `json:"Product_ID"`
	Name string `json:"Variation_Name"`
	Description string `json:"Variation_Description"`
	Price float32 `json:"Variation_Price"`
	PrimaryImage string `json:"Primary_Image,omitempty"`
	VariationQuantity int  `json:"Variation_Quantity"`
	LocationAt string `json:"Location_At"`
}

type ProdExist struct{
	ProductExists bool `json:"Product_Exists"`
	Message string `json:"Message"`
}

type variCrtd struct{
	VariationID int64 `json:"Product_ID"`
	LocationExists bool `json:"Location_Exists"`
}

type ProdInvLocCreation struct{
	VarID int64 `json:"Variation_ID"`
	Quantity int `json:"Quantity"`
	Location string `json:"Location"`
}

type PILCreated struct{
	InvID int64 `json:"Inv_ID"`
	Quantity int `json:"Quantity"`
	Location string `json:"Location"`
}

type CategoryInsert struct{
	CategoryName *string `db:"CategoryName" json:"category_name"`
	CategoryDescription *string `db:"CategoryDescription" json:"category_description"`
}

type CategoryReturn struct{
	CategoryId int64 ` json:"category_id"`
	CategoryName *string ` json:"category_name"`
	CategoryDescription *string ` json:"category_description"`
}

type CategoryEdit struct{
	CategoryId int64 `db:"Category_ID" json:"category_id"`
	CategoryName *string `db:"CategoryName" json:"category_name"`
	CategoryDescription *string `db:"CategoryDescription" json:"category_description"`
}

type CatToCat struct {
	CatStart int `json:"CategoryStart" db:"category_start"`
	CatEnd   int `json:"CategoryEnd" db:"category_end"`
}

type CatToProd struct {
	Cat  int `json:"Category" db:"category"`
	Prod int `json:"Product" db:"product"`
}

type ReadCat struct {
	Category int `json:"category" db:"category"`
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
	PrimaryImage         string  `db:"PRIMARY_IMAGE,omitempty" json:"primary_image,omitempty"`
	VariationQuantity    int     `db:"Quantity" json:"variation_quantity"`
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

type CategoriesList struct{
	collection []CategoryReturn
}

type ProductSize struct {
    SizeID         *int64       `db:"Size_ID" json:"size_id"`
    SizeName       *string     `db:"Size_Name" json:"size_name"`
    SizeDescription *string    `db:"Size_Description,omitempty" json:"size_description,omitempty"`
    VariationID    *int64        `db:"Variation_ID" json:"variation_id"`
    VariationPrice *float64    `db:"Variation_Price" json:"variation_price"`
    SKU            *string    `db:"SKU,omitempty" json:"sku,omitempty"`
	Price		  *float64    `db:"Price,omitempty" json:"price,omitempty"`
    UPC            *string    `db:"UPC,omitempty" json:"upc,omitempty"`
    PrimaryImage   *string    `db:"PRIMARY_IMAGE,omitempty" json:"primary_image,omitempty"`
    DateCreated    *time.Time  `db:"Date_Created" json:"date_created"`
    ModifiedDate   *time.Time `db:"Modified_Date,omitempty" json:"modified_date,omitempty"`
}