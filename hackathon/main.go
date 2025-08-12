package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"file-uploader/handlers"
	"file-uploader/middleware"
	"file-uploader/models"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./app.db" // Default fallback
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatal("Failed to create database directory:", err)
	}

	// Initialize database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize models
	userModel := models.NewUserModel(db)

	// Create tables
	if err := userModel.CreateTable(); err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userModel)

	// Setup routes
	r := mux.NewRouter()

	apiV1Router := r.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	apiV1Router.HandleFunc("/register", authHandler.Register).Methods("POST")
	apiV1Router.HandleFunc("/login", authHandler.Login).Methods("POST")
	apiV1Router.HandleFunc("/revoke", middleware.AuthMiddleware(authHandler.Revoke)).Methods("POST")

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
