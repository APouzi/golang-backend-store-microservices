package routes

import (
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/APouzi/MerchantMachinee/routes/checkout"
	"github.com/APouzi/MerchantMachinee/routes/customer"
	productendpoints "github.com/APouzi/MerchantMachinee/routes/product"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v82"
)


func RouteDigest(digest *chi.Mux, firebaseAuth *firebase.App, stripeClient *stripe.Client, config checkout.Config, redisClient *redis.Client) *chi.Mux{
	// rIndex := indexendpoints.InstanceIndexRoutes(db)

	dbURL := os.Getenv("DBLAYER_URL")
	if dbURL == "" {
		dbURL = "http://dblayer:8080"
	}
	rProduct := productendpoints.InstanceProductsRoutes(dbURL)

	rCheckout := checkout.InstanceCheckoutRoutes(stripeClient, config)
	rCustomer := customer.InstanceCustomerRoutes(firebaseAuth)
	
	// rUser := userendpoints.InstanceUserRoutes(db)
	
	
	AuthMiddleWare := authorization.InjectSystemRefrences(firebaseAuth, redisClient)
	
	// digest.Use(AuthMiddleWare.CheckUserRegistration)
	// rTestRoutes := testroutes.InjectDBRef(db, redis)

	// c := cors.New(cors.Options{
    //     // AllowedOrigins is a list of origins a cross-domain request can be executed from
    //     // All origins are allowed by default, you don't need to set this.
    //     AllowedOrigins: []string{"http://localhost:4200", "http://127.0.0.1:4200"}, //CHANGE LATER
    //     // AllowOriginFunc is a custom function to validate the origin. It takes the origin
    //     // as an argument and returns true if allowed or false otherwise. 
    //     // If AllowOriginFunc is set, AllowedOrigins is ignored.
    //     // AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },

    //     // AllowedMethods is a list of methods the client is allowed to use with
    //     // cross-domain requests. Default is all methods.
    //     AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},

    //     // AllowedHeaders is a list of non simple headers the client is allowed to use with
    //     // cross-domain requests.
    //     AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},

    //     // ExposedHeaders indicates which headers are safe to expose to the API of a CORS
    //     // API specification
    //     ExposedHeaders:   []string{"Link"},
    //     AllowCredentials: false,
    //     // MaxAge indicates how long (in seconds) the results of a preflight request
    //     // can be cached
    //     MaxAge: 300,
    // })
	// digest.Use(digest.Middlewares().Handler)

	digest.Options("/*", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNoContent)
})

	//Index
	// digest.Get("/", rIndex.Index)

	// Testing Routes
	// digest.Get("/products-test-redis",rTestRoutes.GetOneProductRedis)
	// digest.Get("/products-test-sql",rTestRoutes.GetOneProductSQL)
	// digest.Get("/products/test-categories/pullTest", rTestRoutes.PullTestCategory)
	// digest.Post("/products/test-categories", rTestRoutes.CreateTestCategory)


	// digest.Group(func(digest chi.Router){
	// 	// digest.Use(AuthMiddleWare.ValidateToken)
		
	// 	// digest.Get("/users/profile",rUser.UserProfile)
	// })
	// digest.Post("/users/",rUser.Register)
	// digest.Post("/users/login",rUser.Login)

	
	// digest.Post("/superusercreation",rUser.AdminSuperUserCreation)
	digest.Get("/products/{ProductID}",rProduct.GetOneProductsEndPoint)
	digest.Get("/products",rProduct.GetAllProductsAndVariationsEndPoint)
	digest.Get("/products/variations/pagination/", rProduct.GetProductAndVariationsPaginated)
	digest.Get("/products/variations/{productID}", rProduct.GetProductAndVariationsByProductID)
	digest.Get("/search/",rProduct.SearchProductsEndPoint)
	
	digest.Post("/stripe/webhook/payment-confirmation", rCheckout.PaymentConfirmation)
	// digest.Get("/products/{CategoryName}",rProduct.GetProductCategoryEndPointFinal)

	// digest.Get("/categories/",r.GetAllCategories)
	
	// digest.Post("/products/test-categories/InsertTest", rAdmin.InsertIntoFinalProd)

	// Admin need to lockdown based on jwt payload and scope
	
	// digest.Get("/users/{userProfileID:[0-9]+}/wishlists/{wishlistID:[0-9]+}", rWishlist.GetWishListByIDEndpoint)
	digest.Group(func(digest chi.Router){
	digest.Use(AuthMiddleWare.CheckUserRegistration)
	// 	// digest.Use(AuthMiddleWare.HasAdminScope)
	// 	// digest.Post("/products/", rAdmin.CreateProduct)
	digest.Post("/checkout",rCheckout.CreateCheckoutSession)
	digest.Post("/register-login-oauth",rCustomer.RegisterCustomer)
	digest.Get("/customer/profile",rCustomer.GetCustomerProfile)
	digest.Patch("/customer/profile",rCustomer.UpdateCustomerProfile)
	digest.Delete("/customer/profile",rCustomer.DeleteCustomerProfile)
	digest.Get("/users/{userProfileID:[0-9]+}/wishlists/{wishlistID:[0-9]+}", rCustomer.GetCustomerWishList)
	digest.Get("/users/{userProfileID:[0-9]+}/wishlists/all", rCustomer.GetAllCustomerWishLists)
	digest.Post("/wishlists/{wishlistID:[0-9]+}/products", rCustomer.AddProductToWishListEndpoint)
	digest.Post("/users/{userProfileID:[0-9]+}/wishlists/default/products", rCustomer.AddProductToDefaultWishListEndpoint)
	digest.Delete("/wishlists/{wishlistID:[0-9]+}/products", rCustomer.RemoveProductFromWishListEndpoint)
	// digest.Post("/products/", rAdmin.CreateProduct)
	// digest.Post("/products/{ProductID}/variation", rAdmin.CreateVariation)
	// digest.Post("/products/inventory", rAdmin.CreateInventoryLocation)
	// digest.Post("/category/prime", rAdmin.CreatePrimeCategory)
	// digest.Post("/category/sub", rAdmin.CreateSubCategory)
	// digest.Post("/category/final", rAdmin.CreateFinalCategory)
	// digest.Delete("/category/prime/{CatPrimeName}",rAdmin.DeletePrimeCategory)
	// digest.Delete("/category/sub/{CatSubName}",rAdmin.DeleteSubCategory)
	// digest.Delete("/category/final/{CatFinalName}",rAdmin.DeleteFinalCategory)
	// digest.Post("/category/primetosub",rAdmin.ConnectPrimeToSubCategory)
	// digest.Post("/category/subtofinal",rAdmin.ConnectSubToFinalCategory)
	// digest.Post("/category/finaltoprod",rAdmin.ConnectFinalToProdCategory)
	// digest.Get("/category/primes", rAdmin.ReturnAllPrimeCategories)
	// digest.Get("/category/subs", rAdmin.ReturnAllSubCategories)
	// digest.Get("/category/finals", rAdmin.ReturnAllFinalCategories)
	// digest.Patch("/products/{ProductID}",rAdmin.EditProduct)
	// digest.Patch("/variation/{VariationID}",rAdmin.EditVariation)
	// digest.Post("/variation/{VariationID}/attribute",rAdmin.AddAttribute)
	// digest.Patch("/variation/{VariationID}/attribute/{AttributeName}",rAdmin.UpdateAttribute)
	// digest.Delete("/variation/{VariationID}/attribute/{AttributeName}",rAdmin.DeleteAttribute)
	// digest.Post("/admin/{UserID}", rAdmin.UserToAdmin)

	// digest.Get("/tables",rAdmin.GetAllTables)
	return digest
}