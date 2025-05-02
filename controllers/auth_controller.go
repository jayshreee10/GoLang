package controllers

import (
	"encoding/json"
	"go-crud/auth"
	"go-crud/models"
	"go-crud/middlewares"
	"net/http"
	"time"
)

// TokenExpiration is the default JWT token expiration time
const TokenExpiration = 24 * time.Hour

// RegisterUser handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userSignup models.UserSignup
	if err := json.NewDecoder(r.Body).Decode(&userSignup); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input
	if userSignup.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if userSignup.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Minimum password length validation
	if len(userSignup.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	// Register the user
	userID, err := models.RegisterUser(userSignup.Email, userSignup.Password)
	if err != nil {
		if err == models.ErrUserExists {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Error registering user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the user object (without password)
	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Error fetching user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, TokenExpiration)
	if err != nil {
		http.Error(w, "Error generating auth token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return user data and token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		User:  user,
		Token: token,
	})
}

// LoginUser handles user login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var userLogin models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input
	if userLogin.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if userLogin.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Authenticate the user
	user, err := models.LoginUser(userLogin.Email, userLogin.Password)
	if err != nil {
		if err == models.ErrInvalidLogin {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error during login: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, TokenExpiration)
	if err != nil {
		http.Error(w, "Error generating auth token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return user data (without password) and token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		User:  user,
		Token: token,
	})
}

// GetCurrentUser returns the current authenticated user's information
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from context (set by AuthMiddleware)
    userID, ok := middlewares.GetUserID(r)
    if !ok {
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }

    // Get user information
    user, err := models.GetUserByID(userID)
    if err != nil {
        http.Error(w, "Error fetching user data: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// ChangePassword handles password changes for authenticated users
// func ChangePassword(w http.ResponseWriter, r *http.Request) {
// 	// Extract user ID from context (set by AuthMiddleware)
// 	userID, ok := middlewares.GetUserID(r)
// 	if !ok  {
// 		http.Error(w, "User not authenticated", http.StatusUnauthorized)
// 		return
// 	}

// 	userID, ok := userID.(int)
// 	if !ok {
// 		http.Error(w, "Invalid user ID in token", http.StatusInternalServerError)
// 		return
// 	}

// 	// Parse request
// 	var req struct {
// 		CurrentPassword string `json:"current_password"`
// 		NewPassword     string `json:"new_password"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Validate input
// 	if req.CurrentPassword == "" {
// 		http.Error(w, "Current password is required", http.StatusBadRequest)
// 		return
// 	}

// 	if req.NewPassword == "" {
// 		http.Error(w, "New password is required", http.StatusBadRequest)
// 		return
// 	}

// 	if len(req.NewPassword) < 6 {
// 		http.Error(w, "New password must be at least 6 characters long", http.StatusBadRequest)
// 		return
// 	}

// 	// Get user with password
// 	var user models.User
// 	err := models.DB.QueryRow("SELECT id, email, password FROM users WHERE id = ?", userID).Scan(
// 		&user.ID, &user.Email, &user.Password)
// 	if err != nil {
// 		http.Error(w, "Error fetching user data: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Verify current password
// 	if err := auth.CheckPassword(user.Password, req.CurrentPassword); err != nil {
// 		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
// 		return
// 	}

// 	// Update password
// 	if err := models.UpdateUserPassword(userID, req.NewPassword); err != nil {
// 		http.Error(w, "Error updating password: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message": "Password updated successfully",
// 	})
// }