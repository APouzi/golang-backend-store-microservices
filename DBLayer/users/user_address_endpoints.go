package users

import "database/sql"

type UserAddressRoutesTray struct {
	dbInstance *sql.DB
}

func UserAddressRoutesTrayInstance(dbInstance *sql.DB) *UserAddressRoutesTray {
	return &UserAddressRoutesTray{
		dbInstance: dbInstance,
	}
}

func (routes *UserAddressRoutesTray) GetUserAddressesEndpoint() {}

func (routes *UserAddressRoutesTray) CreateUserAddressEndpoint() {}

func (routes *UserAddressRoutesTray) UpdateUserAddressEndpoint() {}

func (routes *UserAddressRoutesTray) DeleteUserAddressEndpoint() {}

