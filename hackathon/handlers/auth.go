package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"file-uploader/models"
	"file-uploader/utils"
)

// AuthHandler handles authentication operations
type AuthHandler struct {
	userModel *models.UserModel
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userModel *models.UserModel) *AuthHandler {
	return &AuthHandler{
		userModel: userModel,
	}
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token   string       `json:"token"`
	User    *models.User `json:"user"`
	Message string       `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON payload"})
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Username and password are required"})
		return
	}

	// Check minimum password length
	if len(req.Password) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Password must be at least 6 characters long"})
		return
	}

	// Create user
	user, err := h.userModel.Create(req.Username, req.Password)
	if err != nil {
		// Check if it's a duplicate username error
		if err.Error() == "UNIQUE constraint failed: users.username" {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Username already exists"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to generate token"})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		Token:   token,
		User:    user,
		Message: "User registered successfully",
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON payload"})
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Username and password are required"})
		return
	}

	// Get user by username
	user, err := h.userModel.GetByUsername(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid credentials"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Database error"})
		return
	}

	// Validate password
	if !user.ValidatePassword(req.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to generate token"})
		return
	}

	// Return success response
	json.NewEncoder(w).Encode(AuthResponse{
		Token:   token,
		User:    user,
		Message: "Login successful",
	})
}

// Revoke handles token revocation (logout)
func (h *AuthHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString, ok := r.Context().Value("token").(string)
	if !ok || tokenString == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No token provided"})
		return
	}

	blacklist := utils.GetTokenBlacklist()
	blacklist.RevokeToken(tokenString)

	// Return success response
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Token revoked successfully",
	})
}
