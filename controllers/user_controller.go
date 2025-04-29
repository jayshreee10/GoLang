package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"go-crud/models"
)

type userRequest struct {
	ID    int    `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
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
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	id, err := models.CreateUser(req.Email)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(userResponse{Message: "User created with ID " + strconv.Itoa(id)})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := models.UpdateUser(req.ID, req.Email); err != nil {
		http.Error(w, "Error updating user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(userResponse{Message: "User updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := models.DeleteUser(req.ID); err != nil {
		http.Error(w, "Error deleting user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(userResponse{Message: "User deleted successfully"})
}