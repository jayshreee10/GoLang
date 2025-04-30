package models

import (
	"database/sql"
	"time"
)

type Order struct {
	ID          int         `json:"id"`
	UserID      int         `json:"user_id"`
	AddressID   *int        `json:"address_id,omitempty"` // Using pointer to handle NULL values
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Items       []OrderItem `json:"items,omitempty"`
	Address     *Address    `json:"address,omitempty"` // Address details
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Product   Product `json:"product,omitempty"`
}

// For creating a new order
type OrderRequest struct {
	UserID    int           `json:"user_id"`
	AddressID int           `json:"address_id,omitempty"`
	Items     []ItemRequest `json:"items"`
}

type ItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// Get all orders with optional limit
func GetOrders(limit int) ([]Order, error) {
	rows, err := DB.Query(`
		SELECT id, user_id, address_id, total_amount, status, created_at, updated_at 
		FROM orders 
		ORDER BY created_at DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var orders []Order
	for rows.Next() {
		var o Order
		var addressID sql.NullInt64
		if err := rows.Scan(&o.ID, &o.UserID, &addressID, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		
		// Handle NULL address_id
		if addressID.Valid {
			addrID := int(addressID.Int64)
			o.AddressID = &addrID
		}
		
		// Get order items for this order
		items, err := GetOrderItems(o.ID)
		if err != nil {
			return nil, err
		}
		o.Items = items
		
		orders = append(orders, o)
	}
	return orders, nil
}

// Get order by ID
func GetOrderByID(id int) (Order, error) {
	var order Order
	var addressID sql.NullInt64
	
	err := DB.QueryRow(`
		SELECT id, user_id, address_id, total_amount, status, created_at, updated_at 
		FROM orders 
		WHERE id = ?`, id).Scan(
		&order.ID, &order.UserID, &addressID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	
	if err != nil {
		return order, err
	}
	
	// Handle NULL address_id
	if addressID.Valid {
		addrID := int(addressID.Int64)
		order.AddressID = &addrID
		
		// Get address details if address_id is not null
		address, err := GetAddressByID(addrID)
		if err == nil {
			order.Address = &address
		}
	}
	
	// Get order items
	items, err := GetOrderItems(order.ID)
	if err != nil {
		return order, err
	}
	order.Items = items
	
	return order, nil
}

// Get order items for a specific order
func GetOrderItems(orderID int) ([]OrderItem, error) {
	rows, err := DB.Query(`
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price,
		       p.name, p.status
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = ?`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var items []OrderItem
	for rows.Next() {
		var oi OrderItem
		var productName, productStatus string
		
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.Quantity, &oi.Price,
			&productName, &productStatus); err != nil {
			return nil, err
		}
		
		// Set basic product info
		oi.Product = Product{
			ID:     oi.ProductID,
			Name:   productName,
			Status: productStatus,
		}
		
		items = append(items, oi)
	}
	return items, nil
}

// Create a new order with items
func CreateOrder(userID int, items []ItemRequest, addressID ...int) (int, error) {
	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // Will be ignored if transaction is committed
	
	// Initialize total amount
	var totalAmount float64 = 0
	
	// Calculate total amount based on product prices and quantities
	for _, item := range items {
		// Get product price
		var price float64
		err := DB.QueryRow("SELECT COALESCE((SELECT price FROM products WHERE id = ?), 0)", item.ProductID).Scan(&price)
		if err != nil {
			return 0, err
		}
		
		// Add to total
		totalAmount += price * float64(item.Quantity)
	}
	
	var result sql.Result
	
	// Insert order (with or without address_id)
	if len(addressID) > 0 && addressID[0] > 0 {
		// With address
		result, err = tx.Exec(
			"INSERT INTO orders (user_id, address_id, total_amount, status, created_at, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", 
			userID, addressID[0], totalAmount, "pending")
	} else {
		// Without address
		result, err = tx.Exec(
			"INSERT INTO orders (user_id, total_amount, status, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", 
			userID, totalAmount, "pending")
	}
	
	if err != nil {
		return 0, err
	}
	
	// Get the order ID
	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	// Insert order items
	for _, item := range items {
		// Get product price again (for safety, though we calculated it above)
		var price float64
		err := DB.QueryRow("SELECT COALESCE((SELECT price FROM products WHERE id = ?), 0)", item.ProductID).Scan(&price)
		if err != nil {
			return 0, err
		}
		
		_, err = tx.Exec(
			"INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)", 
			orderID, item.ProductID, item.Quantity, price)
		if err != nil {
			return 0, err
		}
	}
	
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	
	return int(orderID), nil
}

// Update order status
func UpdateOrderStatus(id int, status string) error {
	_, err := DB.Exec(
		"UPDATE orders SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", 
		status, id)
	return err
}

// Update order address
func UpdateOrderAddress(id int, addressID int) error {
	_, err := DB.Exec(
		"UPDATE orders SET address_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", 
		addressID, id)
	return err
}

// Delete an order and its items
func DeleteOrder(id int) error {
	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Delete order items first (foreign key constraint)
	_, err = tx.Exec("DELETE FROM order_items WHERE order_id = ?", id)
	if err != nil {
		return err
	}
	
	// Delete the order
	_, err = tx.Exec("DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		return err
	}
	
	// Commit the transaction
	return tx.Commit()
}

// GetPendingOrdersByUserID retrieves all pending orders for a specific user
func GetPendingOrdersByUserID(userID int) ([]Order, error) {
	rows, err := DB.Query(`
		SELECT id, user_id, address_id, total_amount, status, created_at, updated_at 
		FROM orders 
		WHERE user_id = ? AND status = 'pending'
		ORDER BY created_at ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var orders []Order
	for rows.Next() {
		var o Order
		var addressID sql.NullInt64
		if err := rows.Scan(&o.ID, &o.UserID, &addressID, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		
		// Handle NULL address_id
		if addressID.Valid {
			addrID := int(addressID.Int64)
			o.AddressID = &addrID
		}
		
		// Get order items for this order
		items, err := GetOrderItems(o.ID)
		if err != nil {
			return nil, err
		}
		o.Items = items
		
		orders = append(orders, o)
	}
	return orders, nil
}

// BatchUpdateOrderStatus updates the status of multiple orders in a single transaction
func BatchUpdateOrderStatus(orderIDs []int, status string) (int, error) {
	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // Will be ignored if transaction is committed
	
	updatedCount := 0
	for _, id := range orderIDs {
		_, err := tx.Exec(
			"UPDATE orders SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", 
			status, id)
		if err != nil {
			return updatedCount, err
		}
		updatedCount++
	}
	
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return updatedCount, err
	}
	
	return updatedCount, nil
}