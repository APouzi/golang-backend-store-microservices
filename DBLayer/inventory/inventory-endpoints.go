package inventory

import (
	"database/sql"
)


func GetInventoryRoutesTrayInstance (databaseInstance *sql.DB) *InventoryRoutesTray{
	return &InventoryRoutesTray{DB: databaseInstance}
}
