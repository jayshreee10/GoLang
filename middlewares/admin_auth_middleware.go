package middlewares

import (
	"encoding/base64"
	"go-crud/config"
	"net/http"
	"strings"
)

// AdminAuthMiddleware checks if the provided credentials match the admin credentials from environment variables
func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract credentials from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Only accept Basic auth for admin validation
		if !strings.HasPrefix(authHeader, "Basic ") {
			http.Error(w, "Authorization header must be Basic authentication", http.StatusUnauthorized)
			return
		}

		// Extract and decode credentials
		credentials, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		if err != nil {
			http.Error(w, "Invalid Basic Auth format", http.StatusUnauthorized)
			return
		}

		// Split into email and password
		credParts := strings.SplitN(string(credentials), ":", 2)
		if len(credParts) != 2 {
			http.Error(w, "Invalid Basic Auth format", http.StatusUnauthorized)
			return
		}

		email := credParts[0]
		password := credParts[1]

		// Validate against environment variables
		if email != config.AppConfig.DefaultAdminEmail || password != config.AppConfig.DefaultAdminPassword {
			http.Error(w, "Invalid admin credentials", http.StatusUnauthorized)
			return
		}

		// Credentials match, proceed to the next handler
		next(w, r)
	}
}