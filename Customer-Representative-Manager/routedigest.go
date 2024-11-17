package main

import (
	orders "github.com/apouzi/customer-representative-manager/Orders"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func RouteDigest(digest *chi.Mux) *chi.Mux {
	c := cors.New(cors.Options{
        // AllowedOrigins is a list of origins a cross-domain request can be executed from
        // All origins are allowed by default, you don't need to set this.
        AllowedOrigins: []string{"http://localhost:4200"},
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

	digest.Post("/pdf",orders.OrderHandler)

	return digest
}