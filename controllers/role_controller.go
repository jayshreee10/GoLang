package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"go-crud/models"
)

type roleRequest struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type roleResponse struct {
	Message string `json:"message"`
}

func GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := models.GetRoles(100)
	if err != nil {
		http.Error(w, "Error fetching roles: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func GetRoleByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing role ID", http.StatusBadRequest)
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}
	
	role, err := models.GetRoleByID(id)
	if err != nil {
		http.Error(w, "Error fetching role: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func CreateRole(w http.ResponseWriter, r *http.Request) {
	var req roleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.Name == "" {
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}
	
	id, err := models.CreateRole(req.Name, req.Description)
	if err != nil {
		http.Error(w, "Error creating role: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roleResponse{Message: "Role created with ID " + strconv.Itoa(id)})
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	var req roleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.ID == 0 {
		http.Error(w, "Role ID is required", http.StatusBadRequest)
		return
	}
	
	if req.Name == "" {
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}
	
	if err := models.UpdateRole(req.ID, req.Name, req.Description); err != nil {
		http.Error(w, "Error updating role: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roleResponse{Message: "Role updated successfully"})
}

func DeleteRole(w http.ResponseWriter, r *http.Request) {
	var req roleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.ID == 0 {
		http.Error(w, "Role ID is required", http.StatusBadRequest)
		return
	}
	
	if err := models.DeleteRole(req.ID); err != nil {
		http.Error(w, "Error deleting role: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roleResponse{Message: "Role deleted successfully"})
}