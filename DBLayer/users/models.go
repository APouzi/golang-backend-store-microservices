package users

import "time"

// CreateUserWithProfileRequest represents the payload needed to create a user and their profile in one call.
type CreateUserWithProfileRequest struct {
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	PhoneNumberCell  string `json:"phone_number_cell"`
	PhoneNumberHome  string `json:"phone_number_home"`
}

// UserProfile represents the persisted user profile record.
type UserProfile struct {
	UserProfileID    int    `json:"user_profile_id"`
	UserID           int    `json:"user_id"`
	PhoneNumberCell  string `json:"phone_number_cell"`
	PhoneNumberHome  string `json:"phone_number_home"`
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
	ProductID   string   `json:"product_Id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       *float64 `json:"price,omitempty"`
	Image       string   `json:"image"`
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

