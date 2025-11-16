package users

import (
	"database/sql"
	"net/http"
)

type UserRoutesTray struct {
	dbInstance *sql.DB
}

func UserRoutesTrayInstance(dbInstance *sql.DB) *UserRoutesTray {
	return &UserRoutesTray{
		dbInstance: dbInstance,
	}
}

func (routes *UserRoutesTray) CreateUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {

}

func (routes *UserRoutesTray) UpdateUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {

}

func (routes *UserRoutesTray) DeleteUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {

}

func (routes *UserRoutesTray) GetUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {

}



