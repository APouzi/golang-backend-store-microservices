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
	var payload User

	if err := helpers.ReadJSON(w, r, &payload); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	updateStmt := `UPDATE tblUser SET FirstName = ?, LastName = ? WHERE Email = ?`
	result, err := routes.dbInstance.Exec(updateStmt, payload.FirstName, payload.LastName, payload.Email)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update user profile"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to retrieve affected rows"})
		return
	}

	if rowsAffected == 0 {
		helpers.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "user profile updated successfully"})


}

func (routes *UserRoutesTray) DeleteUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {

}

func (routes *UserRoutesTray) GetUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		var req struct {
			Email string `json:"email"`
		}
		if err := helpers.ReadJSON(w, r, &req); err == nil {
			email = req.Email
		}
	}
	if email == "" {
		helpers.ErrorJSON(w, fmt.Errorf("email is required"), http.StatusBadRequest)
		return
	}
	query := `SELECT up.UserProfileID, up.UserID, u.FirstName, u.LastName, u.Email, up.PhoneNumberCell, up.PhoneNumberHome
			  FROM tblUserProfile up
			  JOIN tblUser u ON up.UserID = u.UserID
			  WHERE u.Email = ?`

	row := routes.dbInstance.QueryRow(query, email)
	var profile UserProfile
	var user User
	var payload map[string]interface{} = make(map[string]interface{})
	err := row.Scan(&profile.UserProfileID, &profile.UserID, &user.FirstName, &user.LastName, &user.Email, &profile.PhoneNumberCell, &profile.PhoneNumberHome)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.ErrorJSON(w, fmt.Errorf("user profile not found"), http.StatusNotFound)
			return
		}
		helpers.ErrorJSON(w, fmt.Errorf("failed to retrieve user profile"), http.StatusInternalServerError)
		return
	}
	payload["user_id"] = profile.UserID
	payload["user_profile_id"] = profile.UserProfileID
	payload["first_name"] = user.FirstName
	payload["last_name"] = user.LastName
	payload["email"] = user.Email
	payload["phone_number_cell"] = profile.PhoneNumberCell
	payload["phone_number_home"] = profile.PhoneNumberHome
	helpers.WriteJSON(w, http.StatusOK, payload)
}



