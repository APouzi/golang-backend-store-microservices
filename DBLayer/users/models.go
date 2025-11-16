package users

import "time"

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

