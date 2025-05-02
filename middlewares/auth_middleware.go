// Updated auth_middleware.go
package middlewares

import (
    "context"
    "encoding/base64"
    "go-crud/auth"
    "go-crud/models"
    "net/http"
    "strings"
)

type UserContext string

const UserIDKey UserContext = "user_id"

// AuthMiddleware checks for a valid JWT token or Basic Auth and adds user info to the request context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract token from the Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header is required", http.StatusUnauthorized)
            return
        }

        var userID int
        var authenticated bool

        // Check if using Bearer token (JWT)
        if strings.HasPrefix(authHeader, "Bearer ") {
            // Extract and validate the token
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            claims, err := auth.ValidateToken(tokenString)
            if err != nil {
                http.Error(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
                return
            }
            userID = claims.UserID
            authenticated = true
        } else if strings.HasPrefix(authHeader, "Basic ") {
            // Extract credentials from Basic Auth header
            credentials, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
            if err != nil {
                http.Error(w, "Invalid Basic Auth format", http.StatusUnauthorized)
                return
            }
            
            credParts := strings.SplitN(string(credentials), ":", 2)
            if len(credParts) != 2 {
                http.Error(w, "Invalid Basic Auth format", http.StatusUnauthorized)
                return
            }
            
            email := credParts[0]
            password := credParts[1]
            
            // Authenticate user with credentials
            user, err := models.LoginUser(email, password)
            if err != nil {
                http.Error(w, "Invalid credentials", http.StatusUnauthorized)
                return
            }
            
            userID = user.ID
            authenticated = true
        } else {
            http.Error(w, "Authorization header format must be Bearer {token} or Basic {credentials}", http.StatusUnauthorized)
            return
        }

        if authenticated {
            // Add user ID to request context
            ctx := context.WithValue(r.Context(), UserIDKey, userID)
            next(w, r.WithContext(ctx))
            return
        }

        http.Error(w, "Authentication failed", http.StatusUnauthorized)
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

        var userID int
        var authenticated bool

        // Check if using Bearer token (JWT)
        if strings.HasPrefix(authHeader, "Bearer ") {
            // Extract and validate the token
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            claims, err := auth.ValidateToken(tokenString)
            if err == nil {
                userID = claims.UserID
                authenticated = true
            }
        } else if strings.HasPrefix(authHeader, "Basic ") {
            // Extract credentials from Basic Auth header
            credentials, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
            if err == nil {
                credParts := strings.SplitN(string(credentials), ":", 2)
                if len(credParts) == 2 {
                    email := credParts[0]
                    password := credParts[1]
                    
                    // Authenticate user with credentials
                    user, err := models.LoginUser(email, password)
                    if err == nil {
                        userID = user.ID
                        authenticated = true
                    }
                }
            }
        }

        if authenticated {
            // Add user ID to request context
            ctx := context.WithValue(r.Context(), UserIDKey, userID)
            next(w, r.WithContext(ctx))
            return
        }

        // Failed authentication but we still continue
        next(w, r)
    }
}

// GetUserID extracts the user ID from the request context
func GetUserID(r *http.Request) (int, bool) {
    userID, ok := r.Context().Value(UserIDKey).(int)
    return userID, ok
}