package inventory

type Location struct {
	LocationID    int     `json:"location_id" db:"location_id"`
	Name          string  `json:"name" db:"name"`
	Description   *string `json:"description,omitempty" db:"description"`
	Latitude      float64 `json:"latitude" db:"latitude"`
	Longitude     float64 `json:"longitude" db:"longitude"`
	StreetAddress *string `json:"street_address,omitempty" db:"street_address"`
}

type InventoryProductDetail struct {
	InventoryID int     `json:"inventory_id" db:"inventory_id"`
	Quantity    int     `json:"quantity" db:"quantity"`
	ProductID   int     `json:"product_id" db:"product_id"`
	LocationID  int     `json:"location_id" db:"location_id"`
	Description *string `json:"description,omitempty" db:"description"`
}
