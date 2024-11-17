package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type Config struct{

}

func (app *Config) GetAllowedOrigins() []string{
	exeDir, err := filepath.Abs("./")
    if err != nil {
        log.Fatal(err)
    }

	err = godotenv.Load(filepath.Join(exeDir, ".env"))
    if err != nil {
        log.Fatal("Error loading .env file\n","exeDir: ", exeDir)
    }
	allowedHostString := os.Getenv("ALLOWED_HOSTS")
	var AllowedOriginsFromEnv []string
	err = json.Unmarshal([]byte(allowedHostString),&AllowedOriginsFromEnv)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("GetAllowedOrigins completed")
	fmt.Println("AllowedOriginsFromEnv:",AllowedOriginsFromEnv)
	return AllowedOriginsFromEnv
}

func main() {
	const webport = 8002

	// flags to initailize this

	app := Config{
	}

	fmt.Printf("Starting Store Backend on port %d \n", webport)
	// fbDB, err :=fireBaseInit()

	// if err != nil {
	// 	log.Panic(err)
	// }
	serve := &http.Server{
		Addr:    fmt.Sprintf(":%d", webport),
		Handler: app.StartRouter(),
	}
	fmt.Println("\nFirebase App:","\n")

	err := serve.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) StartRouter() http.Handler { // Change the receiver to (*Config)
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   app.GetAllowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           301,
	}))
	// Test if this is working, maybe for microservice
	mux.Use(middleware.Heartbeat("/ping"))

	//Pass the mux to routes to use.
	RouteDigest(mux)
	return mux
}