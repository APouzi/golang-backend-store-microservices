package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type Config struct{
	DB *sql.DB
}

const webport int = 8080


func initDB() (*sql.DB,*database.Models){
	fmt.Println("Start of Initializing DB")

	cfg := mysql.Config{
		User:                 "user",
		Passwd:               "example",
		Net:                  "tcp",
		Addr:                 "mysql:3306",
		DBName:               "database",
		MultiStatements:      true,
	}
	var db *sql.DB
	var err error
	count := 0
	fmt.Println("DB connection has successfully initialized")
	database := &database.Models{}
	return db , database
}




func main() {

	connection, _ := initDB()
	time.Sleep(time.Second*3)
	// flags to initailize this
	var initializeDB, initailizeView string

	flag.StringVar(&initializeDB, "initdb","","Initalize Database")
	flag.StringVar(&initailizeView,"initView","","Intialize Views")
	flag.Parse()

	app := Config{
		DB: connection,
		// Models: models,
		// Redis: rdb,
	}

	

    exeDir, err := filepath.Abs("./")
    if err != nil {
        log.Fatal(err)
    }

    // Load the .env file from the directory
    err = godotenv.Load(filepath.Join(exeDir, ".env"))
    if err != nil {
        log.Fatal("Error loading .env file\n","exeDir: ", exeDir)
    }

	if initializeDB == "t" || initializeDB == "T"{
		initializingpopulation.PopulateProductTables(app.DB)
		// InitiateAndPopulateUsers(app.DB)
		// InitAdminTables(app.DB)
		
	}
	fmt.Printf("Starting Store Backend on port %d \n", webport)

	serve := &http.Server{
		Addr:    fmt.Sprintf(":%d", webport),
		Handler: app.StartRouter(),
	}

	err = serve.ListenAndServe()
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