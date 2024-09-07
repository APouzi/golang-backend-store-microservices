package main

import (
	"database/sql"

	admin "github.com/APouzi/DBLayer/Admin"
	products "github.com/APouzi/DBLayer/Products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)



func RouteDigest(digest *chi.Mux, dbInstance *sql.DB) *chi.Mux{
	// rIndex := indexendpoints.InstanceIndexRoutes(db)

	rProduct := products.GetProductRouteInstance(dbInstance)

	// rUser := userendpoints.InstanceUserRoutes(db)

	rAdmin := admin.GetProductRouteInstance(dbInstance)

	// AuthMiddleWare := authorization.InjectDBRef()

	// rTestRoutes := testroutes.InjectDBRef(db, redis)

	c := cors.New(cors.Options{
        // AllowedOrigins is a list of origins a cross-domain request can be executed from
        // All origins are allowed by default, you don't need to set this.
        AllowedOrigins: []string{"http://localhost:4200"}, //CHANGE LATER
        // AllowOriginFunc is a custom function to validate the origin. It takes the origin
        // as an argument and returns true if allowed or false otherwise. 
        // If AllowOriginFunc is set, AllowedOrigins is ignored.
        // AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },

        // AllowedMethods is a list of methods the client is allowed to use with
        // cross-domain requests. Default is all methods.
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},

        // AllowedHeaders is a list of non simple headers the client is allowed to use with
        // cross-domain requests.
        AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},

        // ExposedHeaders indicates which headers are safe to expose to the API of a CORS
        // API specification
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        // MaxAge indicates how long (in seconds) the results of a preflight request
        // can be cached
        MaxAge: 300,
    })
	digest.Use(c.Handler)

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
	
	digest.Get("/products/{ProductID}",rProduct.GetOneProductEndPoint)
	digest.Get("/db/products/",rProduct.GetAllProductsEndPoint)
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
	digest.Get("/products/variation/{VariationID}", rProduct.GetOneVariationEndPoint)
	digest.Post("/products/inventory-location", rAdmin.CreateInventoryLocation)
	digest.Post("/category/prime", rAdmin.CreatePrimeCategory)
	digest.Post("/category/sub", rAdmin.CreateSubCategory)
	digest.Post("/category/final", rAdmin.CreateFinalCategory)
	digest.Delete("/category/prime/{CatPrimeName}",rAdmin.DeletePrimeCategory)
	digest.Delete("/category/sub/{CatSubName}",rAdmin.DeleteSubCategory)
	digest.Delete("/category/final/{CatFinalName}",rAdmin.DeleteFinalCategory)
	digest.Post("/category/primetosub",rAdmin.ConnectPrimeToSubCategory)
	digest.Post("/category/subtofinal",rAdmin.ConnectSubToFinalCategory)
	digest.Post("/category/finaltoprod",rAdmin.ConnectFinalToProdCategory)
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
	return digest
}