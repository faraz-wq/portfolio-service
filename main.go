package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/faraz-wq/portfolio-service/handlers"
	"github.com/faraz-wq/portfolio-service/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Unable to ping the database: %v\n", err)
	}
}

func main() {
	defer db.Close()

	// Initialize handlers with database connection
	handlers.Init(db)

	// Create a new router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/projects", middleware.APIKeyAuth(handlers.GetProjects)).Methods("GET")
	r.HandleFunc("/projects/{id:[0-9]+}", middleware.APIKeyAuth(handlers.GetProject)).Methods("GET")
	r.HandleFunc("/projects", middleware.APIKeyAuth(handlers.CreateProject)).Methods("POST")
	r.HandleFunc("/projects/{id:[0-9]+}", middleware.APIKeyAuth(handlers.DeleteProject)).Methods("DELETE")

	// apply API key middleware to routes
	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
