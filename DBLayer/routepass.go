package main

import (
	"database/sql"

	admin "github.com/APouzi/DBLayer/admin"
	inventory "github.com/APouzi/DBLayer/inventory"
	"github.com/APouzi/DBLayer/orders"
	products "github.com/APouzi/DBLayer/products"
	"github.com/APouzi/DBLayer/users"
	"github.com/go-chi/chi/v5"
)



func RouteDigest(digest *chi.Mux, dbInstance *sql.DB) *chi.Mux{
	// rIndex := indexendpoints.InstanceIndexRoutes(db)

	rProduct := products.GetProductRouteInstance(dbInstance)

	rCatagories := admin.GetCategoriesRouteInstance(dbInstance)

	rInventory := inventory.GetInventoryRoutesTrayInstance(dbInstance)

	// rUser := userendpoints.InstanceUserRoutes(db)

	rAdmin := admin.GetProductRouteInstance(dbInstance)

	rOrders := orders.GetOrderRoutesTrayInstance(dbInstance)

	rWishlist := users.WishlistRoutesTrayInstance(dbInstance)

	// AuthMiddleWare := authorization.InjectDBRef()

	// rTestRoutes := testroutes.InjectDBRef(db, redis)


	//Index
	// digest.Get("/", rIndex.Index)

	// Testing Routes
	// digest.Get("/products-test-redis",rTestRoutes.GetOneProductRedis)
	// digest.Get("/products-test-sql",rTestRoutes.GetOneProductSQL)
	// digest.Get("/products/test-categories/pullTest", rTestRoutes.PullTestCategory)
	// digest.Post("/products/test-categories", rTestRoutes.CreateTestCategory)


	// digest.Group(func(digest chi.Router){
	// 	digest.Use(AuthMiddleWare.ValidateToken)
	// 	// digest.Get("/users/profile",rUser.UserProfile)
	// })
	// digest.Post("/users/",rUser.Register)
	// digest.Post("/users/login",rUser.Login)

	
	// digest.Post("/superusercreation",rUser.AdminSuperUserCreation)
	
	digest.Get("/variations/{VariationID}", rProduct.GetOneVariationEndPoint) //This needs to change to just 
	digest.Get("/products/variations/",rProduct.GetOneProductVariationByParamEndPoint)
	digest.Get("/products/search", rProduct.SearchProductsEndPoint)
	digest.Get("/products/variations/{productID}",rProduct.GetProductAndVariationsByProductID)
	digest.Get("/products/variations/pagination/",rProduct.GetProductAndVariationsPaginated)
	digest.Get("/products",rProduct.GetAllProductsEndPoint)
	digest.Get("/products/{ProductID}",rProduct.GetOneProductEndPoint)
	digest.Get("/product-size/{SizeID}",rProduct.GetOneProductSizeEndPoint)
	// digest.Get("/products/{CategoryName}",r.GetProductCategoryEndPointFinal)

	// digest.Get("/categories/",r.GetAllCategories)
	
	// digest.Post("/products/test-categories/InsertTest", rAdmin.InsertIntoFinalProd)

	// // Admin need to lockdown based on jwt payload and scope
	// digest.Group(func(digest chi.Router){
	// 	digest.Use(AuthMiddleWare.ValidateToken)
	// 	digest.Use(AuthMiddleWare.HasAdminScope)
		digest.Post("/db/products/", rAdmin.CreateProductMultiChain)
	// })
	digest.Post("/products/{ProductID}/variation", rAdmin.CreateProductVariation)
	digest.Post("/products/inventory-location", rAdmin.CreateInventoryLocation)
	digest.Post("/products/{ProductID}/variation/{VariationID}/size", rAdmin.CreateProductSize)
	digest.Delete("/category/prime/{CatPrimeName}",rCatagories.DeletePrimeCategory)
	digest.Delete("/category/sub/{CatSubName}",rCatagories.DeleteSubCategory)
	digest.Delete("/category/final/{CatFinalName}",rCatagories.DeleteFinalCategory)
	digest.Get("/category/final/",rProduct.GetAllProductsInFinalCategoryViewEndPoint)
	digest.Get("/category/prime/",rProduct.GetAllProductsInPrimeCategoryViewEndPoint)
	digest.Get("/category/sub/",rProduct.GetAllProductsInSubCategoryViewEndPoint)
	digest.Post("/category/prime", rCatagories.CreatePrimeCategory)
	digest.Patch("/category/prime",rCatagories.EditPrimeCategory)
	digest.Patch("/category/sub",rCatagories.EditSubCategory)
	digest.Patch("/category/final",rCatagories.EditFinalCategory)
	digest.Post("/category/sub", rCatagories.CreateSubCategory)
	digest.Post("/category/final", rCatagories.CreateFinalCategory)
	digest.Post("/category/primetosub",rCatagories.ConnectPrimeToSubCategory)
	digest.Post("/category/subtofinal",rCatagories.ConnectSubToFinalCategory)
	digest.Post("/category/finaltoprod",rCatagories.ConnectFinalToProdCategory)
	digest.Get("/category/tree", rProduct.ReturnCategoryTree)
	digest.Get("/category/prime", rProduct.ReturnAllPrimeCategories)
	digest.Get("/category/sub", rProduct.ReturnAllSubCategories)
	digest.Get("/category/final", rProduct.ReturnAllFinalCategories)
	// digest.Get("/category/primes", rAdmin.ReturnAllPrimeCategories)
	// digest.Get("/category/subs", rAdmin.ReturnAllSubCategories)
	// digest.Get("/category/finals", rAdmin.ReturnAllFinalCategories)
	digest.Patch("/products/{ProductID}",rAdmin.EditProduct)
	digest.Patch("/variation/{VariationID}",rAdmin.EditVariation)
	digest.Post("/variation/{VariationID}/attribute",rAdmin.AddAttribute)	
	// digest.Patch("/variation/{VariationID}/attribute/{AttributeName}",rAdmin.UpdateAttribute)
	// digest.Delete("/variation/{VariationID}/attribute/{AttributeName}",rAdmin.DeleteAttribute)
	// digest.Post("/admin/{UserID}", rAdmin.UserToAdmin)

	// digest.Get("/tables",rAdmin.GetAllTables)


	//=====================
	//Inventory Routes
	//=====================
	digest.Get("/inventory/locations",rInventory.GetAllLocations)
	digest.Get("/inventory/locations/",rInventory.GetLocationByParam)
	digest.Get("/inventory/locations/{location-id}",rInventory.GetLocationByID)
	digest.Get("/inventory/inventory-product-details",rInventory.GetAllInventoryProductDetails)
	digest.Get("/inventory/inventory-product-details/",rInventory.GetInventoryProductDetailFromParameter)
	digest.Get("/inventory/inventory-product-details/{inventory-id}",rInventory.GetAllInventoryProductDetailsByID)
	digest.Get("/inventory/inventory-shelf-details",rInventory.GetAllInventoryShelfDetail)
	digest.Get("/inventory/inventory-shelf-details/",rInventory.GetInventoryShelfDetailsByParameter)
	digest.Get("/inventory/inventory-shelf-details/{inventory-shelf-id}",rInventory.GetInventoryShelfDetailByInventoryShelfID)
	digest.Get("/inventory/inventory-location-transfers",rInventory.GetAllLocationTransfers)
	digest.Get("/inventory/inventory-location-transfers/",rInventory.GetLocationTransfersByParam)
	digest.Get("/inventory/inventory-location-transfers/{transfers-id}",rInventory.GetInventoryLocationTransfersById)
	digest.Post("/inventory/inventory-product-details",rInventory.CreateInventoryProductDetail)
	digest.Patch("/inventory/inventory-product-details/{inventory-id}",rInventory.UpdateInventoryProductDetail)
	digest.Patch("/inventory/inventory-shelf-details/{inventory-shelf-id}",rInventory.UpdateInventoryShelfDetail)




	//=====================
	//Order Routes
	//=====================

	digest.Post("/summary-order",rOrders.CreateOrderRecord)
	digest.Post("/order-item-details",rOrders.CreateOrderItemRecord)
	digest.Post("/order-item-details/bulk",rOrders.CreateOrderItemRecordsBulk)
	digest.Post("/payment",rOrders.CreatePayment)
	digest.Post("/refund",rOrders.CreateRefund)
	digest.Get("/tax-codes-intermediary/{SizeID}", rProduct.GetAllProductTaxCodeEndPointFromProductSizeIntermediary)
	digest.Get("/tax-codes/{TaxCodeID}",rProduct.GetOneProductTaxCodeEndpoint)
	digest.Get("/tax-codes",rProduct.GetAllProductTaxCodeEndPoint)


	// =====================
	// Wishlist Routes
	// =====================
	digest.Get("/users/{userProfileID:[0-9]+}/wishlists", rWishlist.GetAllWishListsEndPoint)
	digest.Get("/users/{userProfileID:[0-9]+}/wishlists/{wishlistID:[0-9]+}", rWishlist.GetWishListByIDEndpoint)
	return digest
}