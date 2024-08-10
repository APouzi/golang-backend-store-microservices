package admin

// Product automatically creates Variation
type ProductCreate struct{
	Name string `json:"Product_Name"`
	Description string `json:"Product_Description"`
	Price float64 `json:"Product_Price"`
	VariationName string `json:"Variation_Name"`
	VariationDescription string `json:"Variation_Description"`
	VariationPrice float32 `json:"Variation_Price"`
	VariationQuantity int  `json:"Variation_Quantity"`
	LocationAt string `json:"Location_At"`
}
type ProductCreateRetrieve struct{
	ProductID int64 `json:"Product_ID"`
	VarID int64 `json:"Variation_ID"`
	ProdInvLoc int64 `json:"Inv_ID,omitempty"`

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