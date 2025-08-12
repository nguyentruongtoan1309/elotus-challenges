package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"file-uploader/utils"
)

// AuthMiddleware validates JWT tokens and checks for revocation
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header or form data
		var tokenString string

		// First try Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = r.FormValue("token")
		}

		if tokenString == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing authorization token"})
			return
		}

		blacklist := utils.GetTokenBlacklist()
		if blacklist.IsRevoked(tokenString) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Token has been revoked"})
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token: " + err.Error()})
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "token", tokenString)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
