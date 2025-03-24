package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)
type Config struct{

}

func init() {
    exeDir, err := filepath.Abs("./")
    if err != nil {
        log.Fatal(err)
    }

    // Load the .env file from the directory
    err = godotenv.Load(filepath.Join(exeDir, ".env"))
    if err != nil {
        log.Fatal("Error loading .env file\n","exeDir: ", exeDir)
    }
}


func fireBaseInit() (*firebase.App, error){
	fmt.Println("Firebase Initialization started")
	opt := option.WithCredentialsFile("yangwang-9510b-firebase-adminsdk-ilyhb-eba85e8cbf.json")
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	
	return firebaseApp, nil


}

func main() {
	const webport = 8001

	// flags to initailize this
	var initializeDB, initailizeView string

	flag.StringVar(&initializeDB, "initdb", "", "Initalize Database")
	flag.StringVar(&initailizeView, "initView", "", "Intialize Views")
	flag.Parse()

	// if TestInitCreateThenDelete(app.DB) == false{
	// 	log.Fatal("Connection Test had failed")
	// }
	app := Config{
	}
	fmt.Printf("Starting Store Backend on port %d \n", webport)
	fbDB, err :=fireBaseInit()

	if err != nil {
		log.Panic(err)
	}


	serve := &http.Server{
		Addr:    fmt.Sprintf(":%d", webport),
		Handler: app.StartRouter(fbDB),
	}
	

	

	fmt.Println("\nFirebase App:",fbDB,"\n")

	err = serve.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

	
	// fmt.Println("test", reflect.TypeOf(router))

}

func (app *Config) StartRouter(firebase *firebase.App) http.Handler { // Change the receiver to (*Config)
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
	RouteDigest(mux,firebase)
	return mux
}

func (app *Config) GetAllowedOrigins() []string{
	allowedHostString := os.Getenv("ALLOWED_HOSTS")
	var AllowedOriginsFromEnv []string
	err := json.Unmarshal([]byte(allowedHostString),&AllowedOriginsFromEnv)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("GetAllowedOrigins completed")
	fmt.Println("AllowedOriginsFromEnv:",AllowedOriginsFromEnv)
	return AllowedOriginsFromEnv
}