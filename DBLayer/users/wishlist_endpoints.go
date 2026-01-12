package users

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/APouzi/DBLayer/helpers"
	"github.com/go-chi/chi/v5"
)

type WishlistRoutesTray struct {
	dbInstance *sql.DB
}

func WishlistRoutesTrayInstance(dbInstance *sql.DB) *WishlistRoutesTray {
	return &WishlistRoutesTray{
		dbInstance: dbInstance,
	}
}


func (routes *WishlistRoutesTray) GetAllWishListsEndPoint(w http.ResponseWriter, r *http.Request) {
	userProfileID := chi.URLParam(r, "userProfileID")
	
	id, err := strconv.Atoi(userProfileID)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// First get all wishlists for the user
	wishlistRows, err := routes.dbInstance.Query("SELECT WishlistID, WishlistName, isDefault FROM tblUserWishList WHERE UserID = ?", id)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer wishlistRows.Close()

	var wishlists []Wishlist = make([]Wishlist, 0)

	// Scan all wishlists first
	for wishlistRows.Next() {
		var wishlist Wishlist
		err := wishlistRows.Scan(&wishlist.ID, &wishlist.WishListName, &wishlist.IsDefault)
		if err != nil {
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		wishlist.Items = make([]Products, 0)
		wishlists = append(wishlists, wishlist)
	}

	// Create the map after the slice is fully populated to avoid pointer invalidation
	wishlistMap := make(map[string]*Wishlist)
	for i := range wishlists {
		wishlistMap[wishlists[i].ID] = &wishlists[i]
	}

	// Now get all products for all wishlists through size joins
	productRows, err := routes.dbInstance.Query(`
		SELECT wp.WishlistID, p.Product_ID, p.Product_Name, p.Product_Description, p.PRIMARY_IMAGE, 
		       pv.Variation_Name, ps.Size_Name, ps.Variation_Price
		FROM tblWishlistProduct wp
		JOIN tblProductSize ps ON wp.SizeID = ps.Size_ID
		JOIN tblProductVariation pv ON ps.Variation_ID = pv.Variation_ID
		JOIN tblProducts p ON pv.Product_ID = p.Product_ID
		WHERE wp.WishlistID IN (SELECT WishlistID FROM tblUserWishList WHERE UserID = ?)
	`, id)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer productRows.Close()

	// Group products by wishlist
	for productRows.Next() {
		var wishlistID int
		var product Products
		var image sql.NullString
		var variationName string
		var sizeName string
		var price float64
		
		err := productRows.Scan(&wishlistID, &product.ProductID, &product.Name, &product.Description, &image, &variationName, &sizeName, &price)
		if err != nil {
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		
		if image.Valid {
			product.Image = image.String
		} else {
			product.Image = ""
		}
		product.VariationName = variationName
		product.Size = sizeName
		product.Price = &price

		// Add product to the corresponding wishlist
		wishlistKey := strconv.Itoa(wishlistID)
		if wishlist, exists := wishlistMap[wishlistKey]; exists {
			wishlist.Items = append(wishlist.Items, product)
		}
	}

	helpers.WriteJSON(w, http.StatusOK, wishlists)
}

func (routes *WishlistRoutesTray) GetWishListByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	wishListID := chi.URLParam(r, "wishlistID")
	userProfileID := chi.URLParam(r, "userProfileID")
	
	id, err := strconv.Atoi(wishListID)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userProfileID)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	var wishlist Wishlist
	// Validate that the wishlist belongs to the user and get wishlist details
	err = routes.dbInstance.QueryRow("SELECT WishlistID, WishListName, isDefault FROM tblUserWishList WHERE WishlistID = ? AND UserID = ?", id, userID).Scan(
		&wishlist.ID,
		&wishlist.WishListName,
		&wishlist.IsDefault,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.ErrorJSON(w, fmt.Errorf("wishlist not found"), http.StatusNotFound)
			return
		}
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	rows, err := routes.dbInstance.Query(`
	SELECT p.Product_ID, p.Product_Name, p.Product_Description, p.PRIMARY_IMAGE, 
	       pv.Variation_Name, ps.Size_Name, ps.Variation_Price
	FROM tblWishlistProduct as wp
	JOIN tblProductSize ps ON wp.SizeID = ps.Size_ID
	JOIN tblProductVariation pv ON ps.Variation_ID = pv.Variation_ID
	JOIN tblProducts as p ON pv.Product_ID = p.Product_ID
	WHERE wp.WishlistID = ?`, id)
	if err != nil {
		helpers.ErrorJSON(w, err,http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	wishlist.Items = make([]Products, 0)

	for rows.Next() {
		var item Products
		var image sql.NullString
		var variationName string
		var sizeName string
		var price float64
		err := rows.Scan(&item.ProductID, &item.Name, &item.Description, &image, &variationName, &sizeName, &price)
		if err != nil {
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		if image.Valid {
			item.Image = image.String
		} else {
			item.Image = ""
		}
		item.VariationName = variationName
		item.Size = sizeName
		item.Price = &price
		wishlist.Items = append(wishlist.Items, item)
	}

	helpers.WriteJSON(w, http.StatusOK, wishlist)

}


func (routes *WishlistRoutesTray) CreateWishlistEndpoint(w http.ResponseWriter, r *http.Request) {
	userProfileID := chi.URLParam(r, "userProfileID")
	fmt.Println("made it into DB wishlist creation")
	id, err := strconv.Atoi(userProfileID)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid userProfileID format"), http.StatusBadRequest)
		return
	}

	// Parse request payload for wishlist details
	var payload struct {
		WishlistName string `json:"wishlist_name"`
		IsDefault    *bool  `json:"is_default,omitempty"`
	}
	err = helpers.ReadJSON(w, r, &payload)
	if err != nil {
		fmt.Println("Error reading JSON payload:", err)
		helpers.ErrorJSON(w, fmt.Errorf("invalid request payload"), http.StatusBadRequest)
		return
	}

	// Set default values
	wishlistName := payload.WishlistName
	if wishlistName == "" {
		wishlistName = "My Wishlist"
	}
	isDefault := false
	if payload.IsDefault != nil {
		isDefault = *payload.IsDefault
	}

	fmt.Println("made it into DB wishlist creation 2")
	res, err := routes.dbInstance.Exec("INSERT INTO tblUserWishList (UserID, WishlistName, isDefault) VALUES (?, ?, ?)", id, wishlistName, isDefault)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to create wishlist"), http.StatusInternalServerError)
		return
	}
	var wishlistId int64
	wishlistId, err = res.LastInsertId()
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to retrieve wishlist ID"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, map[string]interface{}{"created": true, "wishlist_id": wishlistId, "wishlist_name": wishlistName})
}

func (routes *WishlistRoutesTray) DeleteWishlistEndpoint(w http.ResponseWriter, r *http.Request) {
	wishListID := chi.URLParam(r, "wishlistID")
	
	// Convert wishListID to integer
	wishlistIDInt, err := strconv.Atoi(wishListID)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid wishlist ID format"), http.StatusBadRequest)
		return
	}



	// Check if wishlist exists and belongs to the user, and get isDefault status
	var isDefault bool
	err = routes.dbInstance.QueryRow("SELECT isDefault FROM tblUserWishList WHERE WishlistID = ?", wishlistIDInt).Scan(&isDefault)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.ErrorJSON(w, fmt.Errorf("wishlist not found"), http.StatusNotFound)
			return
		}
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Prevent deletion of default wishlist
	if isDefault {
		helpers.ErrorJSON(w, fmt.Errorf("cannot delete default wishlist"), http.StatusBadRequest)
		return
	}

	// Delete all products from the wishlist first
	_, err = routes.dbInstance.Exec("DELETE FROM tblWishlistProduct WHERE WishlistID = ?", wishlistIDInt)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to remove products from wishlist"), http.StatusInternalServerError)
		return
	}

	// Delete the wishlist
	result, err := routes.dbInstance.Exec("DELETE FROM tblUserWishList WHERE WishlistID = ?", wishlistIDInt)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to delete wishlist"), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("failed to verify deletion"), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		helpers.ErrorJSON(w, fmt.Errorf("wishlist not found"), http.StatusNotFound)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]any{"deleted": true})
}

func (routes *WishlistRoutesTray) AddProductToWishListEndpoint(w http.ResponseWriter, r *http.Request) {
	wishListID := chi.URLParam(r, "wishlistID")
	fmt.Println("Received request to add product size to wishlist:", wishListID)

	// Convert wishListID to integer
	wishlistIDInt, err := strconv.Atoi(wishListID)
	if err != nil {
		fmt.Println("Error converting wishlist ID to integer:", err)
		helpers.ErrorJSON(w, fmt.Errorf("invalid wishlist ID format"), http.StatusBadRequest)
		return
	}

	// Parse size ID from request body
	var payload struct {
		SizeID int `json:"size_id"`
	}
	err = helpers.ReadJSON(w, r, &payload)
	if err != nil {
		fmt.Println("Error reading JSON payload:", err)
		helpers.ErrorJSON(w, fmt.Errorf("invalid request payload"), http.StatusBadRequest)
		return
	}
	fmt.Println("Parsed size ID:", payload.SizeID)

	// Validate that the wishlist exists
	var exists bool
	err = routes.dbInstance.QueryRow("SELECT EXISTS(SELECT 1 FROM tblUserWishList WHERE WishlistID = ?)", wishlistIDInt).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking if wishlist exists:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to validate wishlist"), http.StatusInternalServerError)
		return
	}
	if !exists {
		fmt.Printf("Wishlist with ID %d does not exist\n", wishlistIDInt)
		helpers.ErrorJSON(w, fmt.Errorf("wishlist not found"), http.StatusNotFound)
		return
	}

	// Check if product size is already in the wishlist to avoid duplicates
	var productExists bool
	err = routes.dbInstance.QueryRow("SELECT EXISTS(SELECT 1 FROM tblWishlistProduct WHERE WishlistID = ? AND SizeID = ?)", wishlistIDInt, payload.SizeID).Scan(&productExists)
	if err != nil {
		fmt.Println("Error checking if product size is already in wishlist:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to validate product size"), http.StatusInternalServerError)
		return
	}
	if productExists {
		fmt.Printf("Product size %d is already in wishlist %d\n", payload.SizeID, wishlistIDInt)
		helpers.ErrorJSON(w, fmt.Errorf("product size already in wishlist"), http.StatusConflict)
		return
	}

	// Insert product size into wishlist
	_, err = routes.dbInstance.Exec("INSERT INTO tblWishlistProduct (WishlistID, SizeID) VALUES (?, ?)", wishlistIDInt, payload.SizeID)
	if err != nil {
		fmt.Println("Error inserting product size into wishlist:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to add product size to wishlist"), http.StatusInternalServerError)
		return
	}

	fmt.Println("Product size added to wishlist successfully")
	helpers.WriteJSON(w, http.StatusCreated, map[string]any{"added": true})
}

func (routes *WishlistRoutesTray) AddProductToDefaultWishListEndpoint(w http.ResponseWriter, r *http.Request) {
	userProfileID := chi.URLParam(r, "userProfileID")
	fmt.Println("Received request to add product size to default wishlist for user:", userProfileID)

	// Convert userProfileID to integer
	userID, err := strconv.Atoi(userProfileID)
	if err != nil {
		fmt.Println("Error converting user profile ID to integer:", err)
		helpers.ErrorJSON(w, fmt.Errorf("invalid user profile ID format"), http.StatusBadRequest)
		return
	}

	// Parse size ID from request body
	var payload struct {
		SizeID int `json:"size_id"`
	}
	err = helpers.ReadJSON(w, r, &payload)
	if err != nil {
		fmt.Println("Error reading JSON payload:", err)
		helpers.ErrorJSON(w, fmt.Errorf("invalid request payload"), http.StatusBadRequest)
		return
	}
	fmt.Println("Parsed size ID:", payload.SizeID)

	// Find the user's default wishlist
	var wishlistID int
	err = routes.dbInstance.QueryRow("SELECT WishlistID FROM tblUserWishList WHERE UserID = ? AND IsDefault = TRUE", userID).Scan(&wishlistID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No default wishlist found for user %d\n", userID)
			helpers.ErrorJSON(w, fmt.Errorf("default wishlist not found"), http.StatusNotFound)
			return
		}
		fmt.Println("Error finding default wishlist:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to find default wishlist"), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Found default wishlist ID: %d\n", wishlistID)

	// Check if product size is already in the wishlist to avoid duplicates
	var productExists bool
	err = routes.dbInstance.QueryRow("SELECT EXISTS(SELECT 1 FROM tblWishlistProduct WHERE WishlistID = ? AND SizeID = ?)", wishlistID, payload.SizeID).Scan(&productExists)
	if err != nil {
		fmt.Println("Error checking if product size is already in wishlist:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to validate product size"), http.StatusInternalServerError)
		return
	}
	if productExists {
		fmt.Printf("Product size %d is already in default wishlist %d\n", payload.SizeID, wishlistID)
		helpers.ErrorJSON(w, fmt.Errorf("product size already in wishlist"), http.StatusConflict)
		return
	}

	// Insert product size into default wishlist
	_, err = routes.dbInstance.Exec("INSERT INTO tblWishlistProduct (WishlistID, SizeID) VALUES (?, ?)", wishlistID, payload.SizeID)
	if err != nil {
		fmt.Println("Error inserting product size into default wishlist:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to add product size to default wishlist"), http.StatusInternalServerError)
		return
	}

	fmt.Println("Product size added to default wishlist successfully")
	helpers.WriteJSON(w, http.StatusCreated, map[string]any{"added": true, "wishlist_id": wishlistID})
}

func (routes *WishlistRoutesTray) RemoveProductFromWishListEndpoint(w http.ResponseWriter, r *http.Request) {
	wishListID := chi.URLParam(r, "wishlistID")
	fmt.Println("Received request to remove product size from wishlist:", wishListID)

	// Convert wishListID to integer
	wishlistIDInt, err := strconv.Atoi(wishListID)
	if err != nil {
		fmt.Println("Error converting wishlist ID to integer:", err)
		helpers.ErrorJSON(w, fmt.Errorf("invalid wishlist ID format"), http.StatusBadRequest)
		return
	}

	// Parse size ID from request body
	var payload struct {
		SizeID int `json:"size_id"`
	}
	err = helpers.ReadJSON(w, r, &payload)
	if err != nil {
		fmt.Println("Error reading JSON payload:", err)
		helpers.ErrorJSON(w, fmt.Errorf("invalid request payload"), http.StatusBadRequest)
		return
	}
	fmt.Println("Parsed size ID:", payload.SizeID)

	// Validate that the wishlist exists
	var exists bool
	err = routes.dbInstance.QueryRow("SELECT EXISTS(SELECT 1 FROM tblUserWishList WHERE WishlistID = ?)", wishlistIDInt).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking if wishlist exists:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to validate wishlist"), http.StatusInternalServerError)
		return
	}
	if !exists {
		fmt.Printf("Wishlist with ID %d does not exist\n", wishlistIDInt)
		helpers.ErrorJSON(w, fmt.Errorf("wishlist not found"), http.StatusNotFound)
		return
	}

	// Check if product size exists in the wishlist
	var productExists bool
	err = routes.dbInstance.QueryRow("SELECT EXISTS(SELECT 1 FROM tblWishlistProduct WHERE WishlistID = ? AND SizeID = ?)", wishlistIDInt, payload.SizeID).Scan(&productExists)
	if err != nil {
		fmt.Println("Error checking if product size is in wishlist:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to validate product size"), http.StatusInternalServerError)
		return
	}
	if !productExists {
		fmt.Printf("Product size %d is not in wishlist %d\n", payload.SizeID, wishlistIDInt)
		helpers.ErrorJSON(w, fmt.Errorf("product size not found in wishlist"), http.StatusNotFound)
		return
	}

	// Remove product size from wishlist
	result, err := routes.dbInstance.Exec("DELETE FROM tblWishlistProduct WHERE WishlistID = ? AND SizeID = ?", wishlistIDInt, payload.SizeID)
	if err != nil {
		fmt.Println("Error removing product size from wishlist:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to remove product size from wishlist"), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected:", err)
		helpers.ErrorJSON(w, fmt.Errorf("failed to verify removal"), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		fmt.Printf("No rows affected when removing product size %d from wishlist %d\n", payload.SizeID, wishlistIDInt)
		helpers.ErrorJSON(w, fmt.Errorf("product size not found in wishlist"), http.StatusNotFound)
		return
	}

	fmt.Println("Product size removed from wishlist successfully")
	helpers.WriteJSON(w, http.StatusOK, map[string]any{"removed": true})
}

