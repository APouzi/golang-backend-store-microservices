package users

import "time"

// CreateUserWithProfileRequest represents the payload needed to create a user and their profile in one call.
type CreateUserWithProfileRequest struct {
	Email                    string  `json:"email"`
	FirstName                string  `json:"first_name"`
	LastName                 string  `json:"last_name"`
	PhoneNumberMobileE164    string  `json:"phone_number_mobile_e164"`
	PhoneNumberHomeE164      string  `json:"phone_number_home_e164"`
	PrimaryShippingAddressID *int    `json:"primary_shipping_address_id,omitempty"`
	PrimaryBillingAddressID  *int    `json:"primary_billing_address_id,omitempty"`
	PreferredLocale          *string `json:"preferred_locale,omitempty"`
	PreferredTimeZone        *string `json:"preferred_timezone,omitempty"`
}

// CreateUserResponse represents the response after successfully creating a user.
type CreateUserResponse struct {
	UserID int64 `json:"user_id"`
}

// UserEmailRequest represents a request containing just an email (e.g. for lookups).
type UserEmailRequest struct {
	Email string `json:"email"`
}

// GenericMessageResponse represents a simple success message.
type GenericMessageResponse struct {
	Message string `json:"message"`
}

// UserProfile represents the persisted user profile record.
type UserProfile struct {
	UserID                int    `json:"user_id"`
	PhoneNumberMobileE164 string `json:"phone_number_mobile_e164"`
	PhoneNumberHomeE164   string `json:"phone_number_home_e164"`
}

// FullUserProfile represents the complete user profile information including joined user data.
type FullUserProfile struct {
	UserID                   int     `json:"user_id"`
	FirstName                string  `json:"first_name"`
	LastName                 string  `json:"last_name"`
	Email                    string  `json:"email"`
	PhoneNumberMobileE164    *string `json:"phone_number_mobile_e164"`
	PhoneNumberHomeE164      *string `json:"phone_number_home_e164"`
	PrimaryShippingAddressID *int    `json:"primary_shipping_address_id"`
	PrimaryBillingAddressID  *int    `json:"primary_billing_address_id"`
	PreferredLocale          *string `json:"preferred_locale"`
	PreferredTimeZone        *string `json:"preferred_timezone"`
}

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Products struct {
	ProductID     string   `json:"product_Id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	VariationName string   `json:"variation_name"`
	Price         *float64 `json:"price,omitempty"`
	Size          string   `json:"size"`
	Image         string   `json:"image"`
}


type Wishlist struct {
	ID          string         `json:"id"`
	WishListName        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	Items       []Products `json:"items"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
	IsDefault   *bool          `json:"isDefault,omitempty"`
}

type CreateWishlistRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IsDefault   *bool   `json:"isDefault,omitempty"`
}

type UpdateWishlistRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsDefault   *bool   `json:"isDefault,omitempty"`
}

type AddWishlistItemRequest struct {
	ProductID   string                 `json:"productId"`
	VariationID *string                `json:"variationId,omitempty"`
	SizeID      string                 `json:"sizeId"`
	Quantity    *int                   `json:"quantity,omitempty"`
	Notes       *string                `json:"notes,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Price       *float64               `json:"price,omitempty"`
	Image       *string                `json:"image,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

