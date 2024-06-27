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
func main() {

}