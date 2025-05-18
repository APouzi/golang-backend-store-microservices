package adminendpoints

import "database/sql"

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
	Name        string `json:"Product_Name"`
	Description string `json:"Product_Description"`
}

type VariationEdit struct {
	VariationID          int64   `json:"Variation_ID"`
	VariationProductID   int64   `json:"Product_ID"`
	VariationName        string  `json:"Variation_Name"`
	VariationDescription string  `json:"Variation_Description"`
	VariationPrice       float32 `json:"Variation_Price"`
	SKU                  string  `json:"SKU"`
	UPC                  string  `json:"UPC"`
	PrimaryImage         string  `json:"Primary_Image,omitempty"`
	VariationQuantity    int     `json:"Variation_Quantity"`
	LocationAt           string  `json:"Location_At"`
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
	CategoryName string `json:"CategoryName"`
	CategoryDescription string `json:"CategoryDescription"`
}

type CategoryReturn struct{
	CategoryId int64 `json:"Category_ID"`
	CategoryName string `json:"CategoryName"`
	CategoryDescription string `json:"CategoryDescription"`
}

type CategoriesList struct{
	collection []CategoryReturn
}