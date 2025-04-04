package userendpoints

type ProductRetrieve struct {
	ProductID   int64  `json:"Product_ID"`
	Name        string `json:"Product_Name"`
	Description string `json:"Product_Description:omitempty"`
}

type ProductWrapper struct {
	Product ProductRetrieve `json:"product"`
}