package middlewares

import (
	"context"
	"go-crud/auth"
	"net/http"
	"strings"
)


type UserContext string

const UserIDKey UserContext = "user_id"

// AuthMiddleware checks for a valid JWT token and adds user info to the request context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		// Extract and validate the token
		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

// Optional middleware that tries to authenticate but doesn't fail if no token is provided
func OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No auth header, just continue
			next(w, r)
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Invalid auth header format, but we'll still continue
			next(w, r)
			return
		}

		// Extract and validate the token
		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err == nil {
			// Add user ID to request context only if token is valid
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next(w, r.WithContext(ctx))
			return
		}

		// Token is invalid but we still continue
		next(w, r)
	}
}

// GetUserID extracts the user ID from the request context
func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	return userID, ok
}