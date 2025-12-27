package users

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/APouzi/DBLayer/helpers"
)

type UserRoutesTray struct {
	dbInstance *sql.DB
}

func UserRoutesTrayInstance(dbInstance *sql.DB) *UserRoutesTray {
	return &UserRoutesTray{
		dbInstance: dbInstance,
	}
}

// CreateUserWithProfileEndpoint inserts into both tblUser and tblUserProfile in a single transaction.
func (routes *UserRoutesTray) CreateUserWithProfileEndpoint(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserWithProfileRequest

	if err := helpers.ReadJSON(w, r, &payload); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if payload.Email == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "email is required"})
		return
	}

	tx, err := routes.dbInstance.Begin()
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to start transaction"})
		return
	}
	defer tx.Rollback()
	fmt.Println("payload:", payload)
	userStmt := `INSERT INTO tblUser (FirstName, LastName, Email) VALUES (?, ?, ?)`
	userResult, err := tx.Exec(userStmt, payload.FirstName, payload.LastName, payload.Email)
	if err != nil {
		helpers.WriteJSON(w, http.StatusConflict, map[string]string{"error": "failed to create user (email may already exist)"})
		return
	}

	userID, err := userResult.LastInsertId()
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to retrieve user ID"})
		return
	}

	profileStmt := `INSERT INTO tblUserProfile (UserID, PhoneNumberCell, PhoneNumberHome) VALUES (?, ?, ?)`
	profileResult, err := tx.Exec(profileStmt, userID, payload.PhoneNumberCell, payload.PhoneNumberHome)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user profile"})
		return
	}

	profileID, err := profileResult.LastInsertId()
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to retrieve profile ID"})
		return
	}

	if err := tx.Commit(); err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to commit transaction"})
		return
	}

	response := map[string]interface{}{
		"user_id":    userID,
		"profile_id": profileID,
	}

	helpers.WriteJSON(w, http.StatusCreated, response)
}

func (routes *UserRoutesTray) UpdateUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {

}

func (routes *UserRoutesTray) DeleteUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {

}

func (routes *UserRoutesTray) GetUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {

}



