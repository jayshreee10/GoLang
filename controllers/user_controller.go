package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"go-crud/models"
	"go-crud/middlewares"
)

type userRequest struct {
	ID       int    `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"` // Added password field
}

type userResponse struct {
	Message string `json:"message"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers(100)
	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserProfile returns the currently authenticated user
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by AuthMiddleware)
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Error fetching user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if current user is an admin
	// This is just a placeholder. In a real app, you'd check user roles
	_, ok := middlewares.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate email
	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Create user with email only (for admin creation)
	// Password can be set later by the user
	id, err := models.CreateUser(req.Email)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{Message: "User created with ID " + strconv.Itoa(id)})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from token
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// If ID is provided in request, verify user is updating their own account
	// or implement admin check here
	if req.ID != 0 && req.ID != userID {
		// In a real app, check if user is admin before allowing this
		http.Error(w, "You can only update your own account", http.StatusForbidden)
		return
	}

	// If no ID provided, use the authenticated user's ID
	if req.ID == 0 {
		req.ID = userID
	}

	// Update email if provided
	if req.Email != "" {
		if err := models.UpdateUser(req.ID, req.Email); err != nil {
			http.Error(w, "Error updating user: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Update password if provided
	if req.Password != "" {
		if len(req.Password) < 6 {
			http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
			return
		}

		if err := models.UpdateUserPassword(req.ID, req.Password); err != nil {
			http.Error(w, "Error updating password: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{Message: "User updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from token
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// If ID is provided in request, verify user is deleting their own account
	// or implement admin check here
	if req.ID != 0 && req.ID != userID {
		// In a real app, check if user is admin before allowing this
		http.Error(w, "You can only delete your own account", http.StatusForbidden)
		return
	}

	// If no ID provided, use the authenticated user's ID
	if req.ID == 0 {
		req.ID = userID
	}

	if err := models.DeleteUser(req.ID); err != nil {
		http.Error(w, "Error deleting user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{Message: "User deleted successfully"})
}