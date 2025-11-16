package users

import (
	"database/sql"
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
		http.Error(w, "Invalid userProfileID", http.StatusBadRequest)
		return
	}

	rows, err := routes.dbInstance.Query("SELECT WishlistID, UserProfileID FROM tblUserWishList WHERE UserProfileID = ?", id)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var wishlists []Wishlist = make([]Wishlist, 0)

	for rows.Next() {
		var wishlist Wishlist
		err := rows.Scan(&wishlist.ID, &wishlist.WishListName)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		wishlists = append(wishlists, wishlist)
	}

	helpers.WriteJSON(w, http.StatusOK, wishlists)


}

func (routes *WishlistRoutesTray) GetWishListByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	wishListID := chi.URLParam(r, "wishlistID")
	id, err := strconv.Atoi(wishListID)
	if err != nil {
		helpers.ErrorJSON(w, err,http.StatusBadRequest)
		return
	}

	var wishlist Wishlist
	routes.dbInstance.QueryRow("SELECT WishlistID, WishListName, isDefault FROM tblUserWishList WHERE WishlistID = ?", id).Scan(
		&wishlist.ID,
		&wishlist.WishListName,
		&wishlist.IsDefault,
	)

	rows, err := routes.dbInstance.Query(`
	SELECT p.Product_ID, p.Product_Name, p.Product_Description, p.PRIMARY_IMAGE
	FROM tblWishlistProduct as wp
	JOIN tblProducts as p ON wp.ProductID = p.Product_ID
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
		err := rows.Scan(&item.ProductID, &item.Name, &item.Description, &image)
		if err != nil {
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		if image.Valid {
			item.Image = image.String
		} else {
			item.Image = ""
		}
		wishlist.Items = append(wishlist.Items, item)
	}

	helpers.WriteJSON(w, http.StatusOK, wishlist)

}


func (routes *WishlistRoutesTray) CreateWishlistEndpoint (w http.ResponseWriter, r *http.Request) {
	
	
}

func (routes *WishlistRoutesTray) DeleteWishlistEndpoint (w http.ResponseWriter, r *http.Request) {
//Make sure to not delete default wishlist and make sure that its not allowed on frontend.
}

func (routes *WishlistRoutesTray) AddProductToWishListEndpoint (w http.ResponseWriter, r *http.Request) {
	
}

func (routes *WishlistRoutesTray) RemoveProductFromWishListEndpoint (w http.ResponseWriter, r *http.Request) {
	
}

