package users

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"

	"github.com/APouzi/DBLayer/helpers"
	"github.com/go-sql-driver/mysql"
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
		helpers.ErrorJSON(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	if payload.Email == "" {
		helpers.ErrorJSON(w, errors.New("email is required"), http.StatusBadRequest)
		return
	}

	// Validate Email
	if _, err := mail.ParseAddress(payload.Email); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid email format"), http.StatusBadRequest)
		return
	}

	// Validate Phone Numbers (E.164 format: +1234567890)
	e164Regex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	if payload.PhoneNumberMobileE164 != "" && !e164Regex.MatchString(payload.PhoneNumberMobileE164) {
		helpers.ErrorJSON(w, fmt.Errorf("invalid mobile phone number format (E.164 required)"), http.StatusBadRequest)
		return
	}
	if payload.PhoneNumberHomeE164 != "" && !e164Regex.MatchString(payload.PhoneNumberHomeE164) {
		helpers.ErrorJSON(w, fmt.Errorf("invalid home phone number format (E.164 required)"), http.StatusBadRequest)
		return
	}

	tx, err := routes.dbInstance.Begin()
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to start transaction"), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	fmt.Println("payload:", payload)

	// Use INSERT ... ON DUPLICATE KEY UPDATE to handle both creation and existing users atomically
	// The UserID=LAST_INSERT_ID(UserID) trick ensures we get the correct ID back
	userStmt := `INSERT INTO tblUser (FirstName, LastName, Email) VALUES (?, ?, ?) 
				 ON DUPLICATE KEY UPDATE UserID=LAST_INSERT_ID(UserID), FirstName=VALUES(FirstName), LastName=VALUES(LastName)`
	
	userResult, err := tx.Exec(userStmt, payload.FirstName, payload.LastName, payload.Email)
	if err != nil {
		status, msg := mapMySQLError(err, "failed to create/update user")
		helpers.ErrorJSON(w, errors.New(msg), status)
		return
	}

	userID, err := userResult.LastInsertId()
	if err != nil {
		fmt.Printf("Failed to retrieve LastInsertId: %v\n", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to retrieve user ID"), http.StatusInternalServerError)
		return
	}
	fmt.Printf("User upserted into tblUser with ID: %d\n", userID)

	// Handle empty strings for phone numbers to ensure they are treated as NULL in DB
	var mobilePhone interface{} = payload.PhoneNumberMobileE164
	if payload.PhoneNumberMobileE164 == "" {
		mobilePhone = nil
	}
	var homePhone interface{} = payload.PhoneNumberHomeE164
	if payload.PhoneNumberHomeE164 == "" {
		homePhone = nil
	}

	// Upsert Profile (Insert or Update if exists)
	profileStmt := `INSERT INTO tblUserProfile (UserID, PhoneNumberMobileE164, PhoneNumberHomeE164, PrimaryShippingAddressID, PrimaryBillingAddressID, PreferredLocale, PreferredTimeZone) 
					VALUES (?, ?, ?, ?, ?, ?, ?)
					ON DUPLICATE KEY UPDATE
					PhoneNumberMobileE164 = VALUES(PhoneNumberMobileE164),
					PhoneNumberHomeE164 = VALUES(PhoneNumberHomeE164),
					PreferredLocale = VALUES(PreferredLocale),
					PreferredTimeZone = VALUES(PreferredTimeZone)`
	
	_, err = tx.Exec(profileStmt, userID, mobilePhone, homePhone, payload.PrimaryShippingAddressID, payload.PrimaryBillingAddressID, payload.PreferredLocale, payload.PreferredTimeZone)
	if err != nil {
		status, msg := mapMySQLError(err, "failed to create/update user profile")
		helpers.ErrorJSON(w, errors.New(msg), status)
		return
	}
	fmt.Println("User profile upserted into tblUserProfile")

	if err := tx.Commit(); err != nil {
		fmt.Printf("Failed to commit transaction: %v\n", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to commit transaction"), http.StatusInternalServerError)
		return
	}
	fmt.Println("Transaction committed successfully")

	response := CreateUserResponse{
		UserID: userID,
	}

	helpers.WriteJSON(w, http.StatusCreated, response)
}

func (routes *UserRoutesTray) UpdateUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserWithProfileRequest

	if err := helpers.ReadJSON(w, r, &payload); err != nil {
		helpers.ErrorJSON(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	if payload.Email == "" {
		helpers.ErrorJSON(w, errors.New("email is required"), http.StatusBadRequest)
		return
	}

	// Validate Email
	if _, err := mail.ParseAddress(payload.Email); err != nil {
		helpers.ErrorJSON(w, errors.New("invalid email format"), http.StatusBadRequest)
		return
	}

	// Validate Phone Numbers (E.164 format: +1234567890)
	e164Regex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	if payload.PhoneNumberMobileE164 != "" && !e164Regex.MatchString(payload.PhoneNumberMobileE164) {
		helpers.ErrorJSON(w, errors.New("invalid mobile phone number format (E.164 required)"), http.StatusBadRequest)
		return
	}
	if payload.PhoneNumberHomeE164 != "" && !e164Regex.MatchString(payload.PhoneNumberHomeE164) {
		helpers.ErrorJSON(w, errors.New("invalid home phone number format (E.164 required)"), http.StatusBadRequest)
		return
	}

	tx, err := routes.dbInstance.Begin()
	if err != nil {
		helpers.ErrorJSON(w, errors.New("failed to start transaction"), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get UserID
	var userID int64
	err = tx.QueryRow("SELECT UserID FROM tblUser WHERE Email = ?", payload.Email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.ErrorJSON(w, errors.New("user not found"), http.StatusNotFound)
			return
		}
		helpers.ErrorJSON(w, errors.New("failed to retrieve user"), http.StatusInternalServerError)
		return
	}

	// Update tblUser
	updateUserStmt := `UPDATE tblUser SET FirstName = ?, LastName = ? WHERE UserID = ?`
	_, err = tx.Exec(updateUserStmt, payload.FirstName, payload.LastName, userID)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("failed to update user basic info"), http.StatusInternalServerError)
		return
	}

	// Handle empty strings for phone numbers
	var mobilePhone interface{} = payload.PhoneNumberMobileE164
	if payload.PhoneNumberMobileE164 == "" {
		mobilePhone = nil
	}
	var homePhone interface{} = payload.PhoneNumberHomeE164
	if payload.PhoneNumberHomeE164 == "" {
		homePhone = nil
	}

	// Upsert tblUserProfile
	upsertProfileStmt := `INSERT INTO tblUserProfile (UserID, PhoneNumberMobileE164, PhoneNumberHomeE164) 
						  VALUES (?, ?, ?) 
						  ON DUPLICATE KEY UPDATE 
						  PhoneNumberMobileE164 = VALUES(PhoneNumberMobileE164), 
						  PhoneNumberHomeE164 = VALUES(PhoneNumberHomeE164)`
	_, err = tx.Exec(upsertProfileStmt, userID, mobilePhone, homePhone)
	if err != nil {
		status, msg := mapMySQLError(err, "failed to update user profile details")
		helpers.ErrorJSON(w, errors.New(msg), status)
		return
	}

	if err := tx.Commit(); err != nil {
		helpers.ErrorJSON(w, errors.New("failed to commit transaction"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, GenericMessageResponse{Message: "user profile updated successfully"})
}

func (routes *UserRoutesTray) DeleteUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {
	rawEmail := r.URL.Query().Get("email")
	if rawEmail == "" {
		helpers.ErrorJSON(w, fmt.Errorf("email is required"), http.StatusBadRequest)
		return
	}

	email, err := url.QueryUnescape(rawEmail)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid email"), http.StatusBadRequest)
		return
	}

	parsed, err := mail.ParseAddress(email)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid email format"), http.StatusBadRequest)
		return
	}
	email = parsed.Address

	var userID int64
	err = routes.dbInstance.QueryRow("SELECT UserID FROM tblUser WHERE Email = ?", email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.ErrorJSON(w, fmt.Errorf("user not found"), http.StatusNotFound)
			return
		}
		helpers.ErrorJSON(w, fmt.Errorf("failed to lookup user"), http.StatusInternalServerError)
		return
	}

	tx, err := routes.dbInstance.Begin()
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to start transaction"), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM tblUserProfile WHERE UserID = ?", userID); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to delete user profile"), http.StatusInternalServerError)
		return
	}

	result, err := tx.Exec("DELETE FROM tblUser WHERE UserID = ?", userID)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to delete user"), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to retrieve affected rows"), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		helpers.ErrorJSON(w, fmt.Errorf("user not found"), http.StatusNotFound)
		return
	}

	if err := tx.Commit(); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to commit transaction"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, GenericMessageResponse{Message: "user and profile deleted successfully"})
}

func (routes *UserRoutesTray) GetUserProfileEndpoint(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		var req UserEmailRequest
		if err := helpers.ReadJSON(w, r, &req); err == nil {
			email = req.Email
		}
	}
	if email == "" {
		helpers.ErrorJSON(w, fmt.Errorf("email is required"), http.StatusBadRequest)
		return
	}
	query := `SELECT up.UserID, u.FirstName, u.LastName, u.Email, up.PhoneNumberMobileE164, up.PhoneNumberHomeE164, up.PrimaryShippingAddressID, up.PrimaryBillingAddressID, up.PreferredLocale, up.PreferredTimeZone
			  FROM tblUserProfile up
			  JOIN tblUser u ON up.UserID = u.UserID
			  WHERE u.Email = ?`

	row := routes.dbInstance.QueryRow(query, email)
	var profile FullUserProfile

	err := row.Scan(
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Email,
		&profile.PhoneNumberMobileE164,
		&profile.PhoneNumberHomeE164,
		&profile.PrimaryShippingAddressID,
		&profile.PrimaryBillingAddressID,
		&profile.PreferredLocale,
		&profile.PreferredTimeZone,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			helpers.ErrorJSON(w, fmt.Errorf("user profile not found"), http.StatusNotFound)
			return
		}
		helpers.ErrorJSON(w, fmt.Errorf("failed to retrieve user profile"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, profile)
}

// mapMySQLError converts common MySQL errors into client-facing HTTP responses.
func mapMySQLError(err error, fallback string) (int, string) {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062:
			return http.StatusConflict, "email already exists"
		case 1452:
			return http.StatusBadRequest, "primary_shipping_address_id or primary_billing_address_id does not reference an existing address"
		case 3819:
			return http.StatusBadRequest, "phone numbers must satisfy E.164 format"
		case 1048:
			return http.StatusBadRequest, "a required field is missing"
		case 1406:
			return http.StatusBadRequest, "one of the fields exceeds the allowed length"
		}
	}

	return http.StatusInternalServerError, fallback
}




