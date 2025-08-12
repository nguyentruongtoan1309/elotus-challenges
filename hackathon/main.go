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
	fileModel := models.NewFileModel(db)

	// Create tables
	if err := userModel.CreateTable(); err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	if err := fileModel.CreateTable(); err != nil {
		log.Fatal("Failed to create files table:", err)
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "/tmp" // Default fallback
	}
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userModel)
	uploadHandler := handlers.NewUploadHandler(fileModel)
	staticHandler := handlers.NewStaticHandler(fileModel)

	// Setup routes
	r := mux.NewRouter()

	apiV1Router := r.PathPrefix("/api/v1").Subrouter()

	// Static file routes
	r.HandleFunc("/files/{fileId:[0-9]+}", middleware.AuthMiddleware(staticHandler.ServeFile)).Methods("GET")
	r.HandleFunc("/public/files/{fileId:[0-9]+}", staticHandler.ServePublicFile).Methods("GET")

	// Auth routes
	apiV1Router.HandleFunc("/register", authHandler.Register).Methods("POST")
	apiV1Router.HandleFunc("/login", authHandler.Login).Methods("POST")
	apiV1Router.HandleFunc("/revoke", middleware.AuthMiddleware(authHandler.Revoke)).Methods("POST")

	// Upload routes
	apiV1Router.HandleFunc("/upload", middleware.AuthMiddleware(uploadHandler.Upload)).Methods("POST")
	// Simple HTML form for testing (as requested - not pretty)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<body>
			<h2>JWT Authentication & File Upload Test</h2>
			
			<h3>1. Register</h3>
			<form id="registerForm">
				<input type="text" id="regUsername" placeholder="Username" required><br><br>
				<input type="password" id="regPassword" placeholder="Password" required><br><br>
				<button type="submit">Register</button>
			</form>
			<div id="registerResult"></div>

			<h3>2. Login</h3>
			<form id="loginForm">
				<input type="text" id="loginUsername" placeholder="Username" required><br><br>
				<input type="password" id="loginPassword" placeholder="Password" required><br><br>
				<button type="submit">Login</button>
			</form>
			<div id="loginResult"></div>

			<h3>3. File Upload</h3>
			<form action="/api/v1/upload" method="post" enctype="multipart/form-data">
				<input type="file" name="data" accept="image/*" required><br><br>
				<input type="text" name="token" id="tokenField" placeholder="JWT Token (get from login)" required><br><br>
				<input type="submit" value="Upload Image">
			</form>

			<script>
				// Register form handler
				document.getElementById('registerForm').addEventListener('submit', async (e) => {
					e.preventDefault();
					const username = document.getElementById('regUsername').value;
					const password = document.getElementById('regPassword').value;
					
					try {
						const response = await fetch('/api/v1/register', {
							method: 'POST',
							headers: { 'Content-Type': 'application/json' },
							body: JSON.stringify({ username, password })
						});
						const result = await response.json();
						document.getElementById('registerResult').innerHTML = '<pre>' + JSON.stringify(result, null, 2) + '</pre>';
						if (result.token) {
							document.getElementById('tokenField').value = result.token;
						}
					} catch (error) {
						document.getElementById('registerResult').innerHTML = 'Error: ' + error.message;
					}
				});

				// Login form handler
				document.getElementById('loginForm').addEventListener('submit', async (e) => {
					e.preventDefault();
					const username = document.getElementById('loginUsername').value;
					const password = document.getElementById('loginPassword').value;
					
					try {
						const response = await fetch('/api/v1/login', {
							method: 'POST',
							headers: { 'Content-Type': 'application/json' },
							body: JSON.stringify({ username, password })
						});
						const result = await response.json();
						document.getElementById('loginResult').innerHTML = '<pre>' + JSON.stringify(result, null, 2) + '</pre>';
						if (result.token) {
							document.getElementById('tokenField').value = result.token;
						}
					} catch (error) {
						document.getElementById('loginResult').innerHTML = 'Error: ' + error.message;
					}
				});
			</script>
		</body>
		</html>`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}).Methods("GET")

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
