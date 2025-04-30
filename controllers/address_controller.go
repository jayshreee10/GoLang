package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"go-crud/models"
)

type addressRequest struct {
	ID          int    `json:"id,omitempty"`
	UserID      int    `json:"user_id"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	IsDefault   bool   `json:"is_default"`
}

type addressResponse struct {
	Message string `json:"message"`
	ID      int    `json:"id,omitempty"`
}

type assignAddressRequest struct {
	OrderID   int `json:"order_id"`
	AddressID int `json:"address_id"`
}

// GetAddresses handles retrieving all addresses
func GetAddresses(w http.ResponseWriter, r *http.Request) {
	addresses, err := models.GetAddresses(100)
	if err != nil {
		http.Error(w, "Error fetching addresses: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addresses)
}

// GetAddressByID handles retrieving a specific address
func GetAddressByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing address ID", http.StatusBadRequest)
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}
	
	address, err := models.GetAddressByID(id)
	if err != nil {
		http.Error(w, "Error fetching address: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

// CreateAddress handles creating a new address
func CreateAddress(w http.ResponseWriter, r *http.Request) {
	var req addressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.UserID == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	
	if req.StreetLine1 == "" {
		http.Error(w, "Street address is required", http.StatusBadRequest)
		return
	}
	
	if req.City == "" {
		http.Error(w, "City is required", http.StatusBadRequest)
		return
	}
	
	if req.State == "" {
		http.Error(w, "State is required", http.StatusBadRequest)
		return
	}
	
	if req.PostalCode == "" {
		http.Error(w, "Postal code is required", http.StatusBadRequest)
		return
	}
	
	if req.Country == "" {
		http.Error(w, "Country is required", http.StatusBadRequest)
		return
	}
	
	// Create the address
	id, err := models.CreateAddress(
		req.UserID,
		req.StreetLine1,
		req.StreetLine2,
		req.City,
		req.State,
		req.PostalCode,
		req.Country,
		req.IsDefault,
	)
	
	if err != nil {
		http.Error(w, "Error creating address: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Update pending orders status to processing for this user
	// Get all pending orders for this user
	rows, err := models.DB.Query(`
		SELECT id FROM orders 
		WHERE user_id = ? AND status = 'pending'`, req.UserID)
	if err != nil {
		http.Error(w, "Error fetching pending orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	// Update each order's status to 'processing'
	updatedCount := 0
	for rows.Next() {
		var orderID int
		if err := rows.Scan(&orderID); err != nil {
			http.Error(w, "Error processing orders: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		if err := models.UpdateOrderStatus(orderID, "processing"); err != nil {
			http.Error(w, "Error updating order status: "+err.Error(), http.StatusInternalServerError)
			return
		}
		updatedCount++
	}
	
	// Add info about updated orders to the response
	responseMsg := "Address created successfully"
	if updatedCount > 0 {
		responseMsg += ". " + strconv.Itoa(updatedCount) + " pending orders updated to processing"
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addressResponse{
		Message: responseMsg,
		ID:      id,
	})
}

// UpdateAddress handles updating an existing address
func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	var req addressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.ID == 0 {
		http.Error(w, "Address ID is required", http.StatusBadRequest)
		return
	}
	
	if req.UserID == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	
	if req.StreetLine1 == "" {
		http.Error(w, "Street address is required", http.StatusBadRequest)
		return
	}
	
	if req.City == "" {
		http.Error(w, "City is required", http.StatusBadRequest)
		return
	}
	
	if req.State == "" {
		http.Error(w, "State is required", http.StatusBadRequest)
		return
	}
	
	if req.PostalCode == "" {
		http.Error(w, "Postal code is required", http.StatusBadRequest)
		return
	}
	
	if req.Country == "" {
		http.Error(w, "Country is required", http.StatusBadRequest)
		return
	}
	
	// Update the address
	err := models.UpdateAddress(
		req.ID,
		req.UserID,
		req.StreetLine1,
		req.StreetLine2,
		req.City,
		req.State,
		req.PostalCode,
		req.Country,
		req.IsDefault,
	)
	
	if err != nil {
		http.Error(w, "Error updating address: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Update pending orders status to processing for this user
	// Get all pending orders for this user
	rows, err := models.DB.Query(`
		SELECT id FROM orders 
		WHERE user_id = ? AND status = 'pending'`, req.UserID)
	if err != nil {
		http.Error(w, "Error fetching pending orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	// Update each order's status to 'processing'
	updatedCount := 0
	for rows.Next() {
		var orderID int
		if err := rows.Scan(&orderID); err != nil {
			http.Error(w, "Error processing orders: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		if err := models.UpdateOrderStatus(orderID, "processing"); err != nil {
			http.Error(w, "Error updating order status: "+err.Error(), http.StatusInternalServerError)
			return
		}
		updatedCount++
	}
	
	// Add info about updated orders to the response
	responseMsg := "Address updated successfully"
	if updatedCount > 0 {
		responseMsg += ". " + strconv.Itoa(updatedCount) + " pending orders updated to processing"
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addressResponse{
		Message: responseMsg,
	})
}

// DeleteAddress handles deleting an address
func DeleteAddress(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     int `json:"id"`
		UserID int `json:"user_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.ID == 0 {
		http.Error(w, "Address ID is required", http.StatusBadRequest)
		return
	}
	
	if req.UserID == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	
	// Delete the address
	err := models.DeleteAddress(req.ID, req.UserID)
	if err != nil {
		if err == models.ErrAddressInUse {
			http.Error(w, "Address cannot be deleted because it is being used by one or more orders", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error deleting address: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addressResponse{
		Message: "Address deleted successfully",
	})
}

// AssignAddressToOrder handles assigning an address to an order
func AssignAddressToOrder(w http.ResponseWriter, r *http.Request) {
	var req assignAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.OrderID == 0 {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}
	
	if req.AddressID == 0 {
		http.Error(w, "Address ID is required", http.StatusBadRequest)
		return
	}
	
	// Get order status first
	var orderStatus string
	err := models.DB.QueryRow("SELECT status FROM orders WHERE id = ?", req.OrderID).Scan(&orderStatus)
	if err != nil {
		http.Error(w, "Error fetching order status: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Update the order with address
	err = models.UpdateOrderAddress(req.OrderID, req.AddressID)
	if err != nil {
		http.Error(w, "Error assigning address to order: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// If the order is 'pending', change it to 'processing'
	responseMsg := "Address assigned to order successfully"
	if orderStatus == "pending" {
		err = models.UpdateOrderStatus(req.OrderID, "processing")
		if err != nil {
			http.Error(w, "Warning: Address assigned but failed to update order status: "+err.Error(), http.StatusInternalServerError)
			return
		}
		responseMsg += ". Order status updated from 'pending' to 'processing'"
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addressResponse{
		Message: responseMsg,
	})
}