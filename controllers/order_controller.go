package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"go-crud/models"
)

type orderResponse struct {
	Message string      `json:"message"`
	OrderID int         `json:"order_id,omitempty"`
	Order   interface{} `json:"order,omitempty"`
}

// GetOrders handles retrieving all orders
func GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := models.GetOrders(100)
	if err != nil {
		http.Error(w, "Error fetching orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetOrderByID handles retrieving a specific order
func GetOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	
	order, err := models.GetOrderByID(id)
	if err != nil {
		http.Error(w, "Error fetching order: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// PlaceOrder handles creating a new order
func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.UserID == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	
	if len(req.Items) == 0 {
		http.Error(w, "Order must contain at least one item", http.StatusBadRequest)
		return
	}
	
	// Create the order
	var orderID int
	var err error
	
	if req.AddressID > 0 {
		// Create order with address
		orderID, err = models.CreateOrder(req.UserID, req.Items, req.AddressID)
	} else {
		// Create order without address
		orderID, err = models.CreateOrder(req.UserID, req.Items)
	}
	
	if err != nil {
		http.Error(w, "Error creating order: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get the created order details
	order, err := models.GetOrderByID(orderID)
	if err != nil {
		// Still report success even if we can't fetch the order details
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orderResponse{
			Message: "Order placed successfully",
			OrderID: orderID,
		})
		return
	}
	
	// Return the full order details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderResponse{
		Message: "Order placed successfully",
		OrderID: orderID,
		Order:   order,
	})
}

// UpdateOrderStatus handles updating an order's status
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	type updateRequest struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	}
	
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.ID == 0 {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}
	
	if req.Status == "" {
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}
	
	// Update the order status
	if err := models.UpdateOrderStatus(req.ID, req.Status); err != nil {
		http.Error(w, "Error updating order status: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderResponse{
		Message: "Order status updated successfully",
	})
}

// DeleteOrder handles deleting an order
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	type deleteRequest struct {
		ID int `json:"id"`
	}
	
	var req deleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.ID == 0 {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}
	
	// Delete the order
	if err := models.DeleteOrder(req.ID); err != nil {
		http.Error(w, "Error deleting order: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderResponse{
		Message: "Order deleted successfully",
	})
}