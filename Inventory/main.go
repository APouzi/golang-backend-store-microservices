package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)    // Adds a unique request ID to each request
	r.Use(middleware.RealIP)       // Sets RemoteAddr to the client's real IP
	r.Use(middleware.Recoverer)    // Recovers from panics and writes a 500 if there was one

	// Basic GET endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	
	routerDigest(r)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
