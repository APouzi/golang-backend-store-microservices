package main

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	adminendpoints "github.com/Apouzi/Golang-Admin-Service/admin"
	"github.com/Apouzi/Golang-Admin-Service/middleware"
	"github.com/go-chi/chi/v5"
)



func RouteDigest(digest *chi.Mux, firebaseAuth *firebase.App) *chi.Mux{
	// rIndex := indexendpoints.InstanceIndexRoutes(db)

	// rProduct := productendpoints.InstanceProductsRoutes()

	// rUser := userendpoints.InstanceUserRoutes(db)
	fmt.Println("Firebase App before Auth call:", firebaseAuth)
	rAdmin := adminendpoints.InstanceAdminRoutes()
	rAdminCategories := adminendpoints.InstanceAdminCategoriesRoutes()
	
	fireAuth, err := firebaseAuth.Auth(context.Background())
	if err != nil{
		fmt.Println("Error trying to get Firebase Auth client:", err)
		panic(fmt.Sprintf("Failed to initialize Firebase Auth: %v", err))
	}
	
	if fireAuth == nil {
		fmt.Println("Firebase Auth client is nil after initialization")
		panic("Firebase Auth client is nil")
	}
	
	fmt.Println("Firebase Auth client after call:", fireAuth)
	
	// Create and initialize the middleware
	AuthMiddleWare := &middleware.AuthMiddleWare{
		Client: fireAuth,
	}
	
	fmt.Println("AuthMiddleWare initialized:", AuthMiddleWare, "Client:", AuthMiddleWare.Client)

	// rTestRoutes := testroutes.InjectDBRef(db, redis)


	//Index
	// digest.Get("/", rIndex.Index)

	// Testing Routes
	// digest.Get("/products-test-redis",rTestRoutes.GetOneProductRedis)
	// digest.Get("/products-test-sql",rTestRoutes.GetOneProductSQL)
	// digest.Get("/products/test-categories/pullTest", rTestRoutes.PullTestCategory)
	// digest.Post("/products/test-categories", rTestRoutes.CreateTestCategory)


	digest.Group(func(digest chi.Router){
		digest.Use(AuthMiddleWare.ValidateToken)
		digest.Post("/products/", rAdmin.CreateProduct)
		digest.Post("/products/{ProductID}/variation", rAdmin.CreateVariation)
		digest.Post("/products/{ProductID}/variation/{VariationID}/size", rAdmin.CreateProductSize)
		digest.Post("/products/inventory", rAdmin.CreateInventoryLocation)
		digest.Post("/category/prime", rAdminCategories.CreatePrimeCategory)
		digest.Post("/category/sub", rAdminCategories.CreateSubCategory)
		digest.Post("/category/final", rAdminCategories.CreateFinalCategory)
		digest.Delete("/category/prime/{CatPrimeName}",rAdmin.DeletePrimeCategory)
		digest.Delete("/category/sub/{CatSubName}",rAdmin.DeleteSubCategory)
		digest.Delete("/category/final/{CatFinalName}",rAdmin.DeleteFinalCategory)
		digest.Post("/category/primetosub",rAdminCategories.ConnectPrimeToSubCategory)
		digest.Post("/category/subtofinal",rAdminCategories.ConnectSubToFinalCategory)
		digest.Post("/category/finaltoprod",rAdminCategories.ConnectFinalToProdCategory)
		digest.Get("/category/primes", rAdminCategories.ReturnAllPrimeCategories)
		digest.Get("/category/subs", rAdminCategories.ReturnAllSubCategories)
		digest.Get("/category/finals", rAdminCategories.ReturnAllFinalCategories)
		digest.Patch("/products/{ProductID}",rAdmin.EditProduct)
		digest.Patch("/variation/{VariationID}",rAdmin.EditVariation)
		digest.Post("/variation/{VariationID}/attribute",rAdmin.AddAttribute)
		digest.Patch("/variation/{VariationID}/attribute/{AttributeName}",rAdmin.UpdateAttribute)
		digest.Delete("/variation/{VariationID}/attribute/{AttributeName}",rAdmin.DeleteAttribute)
		digest.Post("/admin/{UserID}", rAdmin.UserToAdmin)
		digest.Get("/tables",rAdmin.GetAllTables)
		// digest.Get("/users/profile",rUser.UserProfile)
	})
	// digest.Post("/users/",rUser.Register)
	// digest.Post("/users/login",rUser.Login)

	
	// digest.Post("/superusercreation",rUser.AdminSuperUserCreation)
	
	// digest.Get("/products/{ProductID}",rProduct.GetOneProductsEndPoint)
	// digest.Get("/products/",rProduct.GetAllProductsEndPoint)
	// digest.Get("/products/{CategoryName}",rProduct.GetProductCategoryEndPointFinal)

	// digest.Get("/categories/",r.GetAllCategories)
	
	// digest.Post("/products/test-categories/InsertTest", rAdmin.InsertIntoFinalProd)

	// Admin need to lockdown based on jwt payload and scope
	// digest.Group(func(digest chi.Router){
		// digest.Use(AuthMiddleWare.ValidateToken)
		// digest.Use(AuthMiddleWare.HasAdminScope)
		// digest.Post("/products/", rAdmin.CreateProduct)
	// })
	// digest.Post("/products/", rAdmin.CreateProduct)
	

	
	return digest
}