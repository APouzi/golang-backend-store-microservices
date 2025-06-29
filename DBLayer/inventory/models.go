package inventory

import (
	"database/sql"
	"time"
)

type InventoryRoutesTray struct {
	DB *sql.DB
}

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


type InventoryShelfDetail struct {
    InventoryShelfID  int    `db:"inventory_shelf_id" json:"inventory_shelf_id"`
    InventoryID       int    `db:"inventory_id" json:"inventory_id"`
    QuantityAtShelf   int    `db:"quantity_at_shelf" json:"quantity_at_shelf"`
    ProductID         int    `db:"product_id" json:"product_id"`
    Shelf             string `db:"shelf" json:"shelf"`
}


type InventoryLocationTransfer struct {
    TransfersID          int       `db:"transfers_id" json:"transfers_id"`
    SourceLocationID     int       `db:"source_location_id" json:"source_location_id"`
    DestinationLocationID int      `db:"destination_location_id" json:"destination_location_id"`
    ProductID            int       `db:"product_id" json:"product_id"`
    Quantity             int       `db:"quantity" json:"quantity"`
    TransferDate         time.Time `db:"transfer_date" json:"transfer_date"`
    Description          *string   `db:"description" json:"description,omitempty"`
    Status               *string   `db:"status" json:"status,omitempty"`
}


