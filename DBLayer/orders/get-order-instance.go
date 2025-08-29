package orders

import "database/sql"

func GetOrderRoutesTrayInstance(databaseInstance *sql.DB) *OrderRoutesTray {
	return &OrderRoutesTray{DB: databaseInstance}
}