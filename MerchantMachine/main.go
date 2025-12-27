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
	"time"

	firebase "firebase.google.com/go"
	"github.com/APouzi/MerchantMachinee/routes"
	"github.com/APouzi/MerchantMachinee/routes/checkout"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v82"
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
	const webport = 8000
	var initializeDB, initailizeView string

	flag.StringVar(&initializeDB, "initdb", "", "Initalize Database")
	flag.StringVar(&initailizeView, "initView", "", "Intialize Views")
	flag.Parse()

	// if TestInitCreateThenDelete(app.DB) == false{
	// 	log.Fatal("Connection Test had failed")
	// }
	app := Config{}

	fmt.Printf("Starting Store Backend on port %d \n", webport)
	fbDB, err :=fireBaseInit()

	if err != nil {
		log.Panic(err)
	}

	sc := NewStripeClient()
	redis_connection, err := ConnectToRedis()
	if err != nil {
		log.Panic(err)
	}
	serve := &http.Server{
		Addr:    fmt.Sprintf(":%d", webport),
		Handler: app.StartRouter(fbDB, sc, redis_connection),
	}


	err = serve.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func (app *Config) StartRouter(firebase *firebase.App, stripeclient *stripe.Client, redis_client *redis.Client) http.Handler { // Change the receiver to (*Config)
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   app.GetAllowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Link","Content-Type","Accept","Accept", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           301,
	}))
	// Test if this is working, maybe for microservice
	mux.Use(middleware.Heartbeat("/ping"))

	
	checkout_config := checkout.Config{
		STRIPE_KEY: os.Getenv("STRIPE_KEY"),
		STRIPE_WEBHOOK_KEY: os.Getenv("STRIPE_WEBHOOK_KEY"),
	}

	//Pass the mux to routes to use.
	routes.RouteDigest(mux,firebase, stripeclient, checkout_config, redis_client)
	return mux
}

func ConnectToRedis() (*redis.Client, error) {
	redis_port := os.Getenv("REDIS_PORT")
	if redis_port == "" {
		redis_port = "6379"
	}
	redis_host := os.Getenv("REDIS_HOST")
	if redis_host == "" || redis_host == "user" {
		redis_host = "redis"
	}
	redis_password := os.Getenv("REDIS_PASSWORD")
	redis_username := os.Getenv("REDIS_USERNAME")
	if redis_username == "" || redis_username == "user" {
		redis_username = ""
	}

	fmt.Println("Connecting to Redis at ", redis_host, ":", redis_port)
	rdb := redis.NewClient(&redis.Options{
		Username: redis_username,
		Addr:     fmt.Sprintf("%s:%s", redis_host, redis_port),
		Password: redis_password,
		DB:       0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}

	return rdb, nil
}

func NewStripeClient() *stripe.Client {
    return stripe.NewClient(stripe.Key)
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