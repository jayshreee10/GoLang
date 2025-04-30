package models

import (
	"database/sql"
	"errors"
	"log"
)

var DB *sql.DB
// Define a custom error for "no rows"
var ErrNoRows = errors.New("no rows found")

func InitDB() error {
	// Create users table
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}
	log.Println("Users table created successfully")


	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		status TEXT DEFAULT 'active',
		price REAL DEFAULT 0.0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}
	log.Println("Products table created successfully")
	
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS roles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}
	log.Println("Roles table created successfully")

	// Create orders table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		address_id INTEGER DEFAULT NULL,
		total_amount REAL NOT NULL,
		status TEXT DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`)
	if err != nil {
		return err
	}
	log.Println("Orders table created successfully")

	// Create order_items table for order-product relationships
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS order_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		price REAL NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id),
		FOREIGN KEY (product_id) REFERENCES products(id)
	)`)
	if err != nil {
		return err
	}
	log.Println("Order items table created successfully")

	// Create addresses table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS addresses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		street_line1 TEXT NOT NULL,
		street_line2 TEXT,
		city TEXT NOT NULL,
		state TEXT NOT NULL,
		postal_code TEXT NOT NULL,
		country TEXT NOT NULL,
		is_default BOOLEAN DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`)
	if err != nil {
		return err
	}
	log.Println("Addresses table created successfully")

	err = seedInitialData()
	if err != nil {
		return err
	}

	return nil
}

func seedInitialData() error {
	// Add sample users
	users := []struct {
		Email string
	}{
		{"admin@example.com"},
		{"user1@example.com"},
		{"user2@example.com"},
	}

	for _, user := range users {
		_, err := DB.Exec("INSERT OR IGNORE INTO users (email) VALUES (?)", user.Email)
		if err != nil {
			return err
		}
	}

	// Add sample products
	products := []struct {
		Name   string
		Status string
		Price  float64
	}{
		{"Product 1", "active", 49.99},
		{"Product 2", "inactive", 149.95},
		{"Product 3", "active", 29.99},
	}
	
	for _, product := range products {
		_, err := DB.Exec("INSERT OR IGNORE INTO products (name, status, price) VALUES (?, ?, ?)", 
			product.Name, product.Status, product.Price)
		if err != nil {
			return err
		}
	}

	roles := []struct {
		Name        string
		Description string
	}{
		{"Admin", "Full system access"},
		{"Editor", "Can edit content"},
		{"Viewer", "Read-only access"},
	}

	for _, role := range roles {
		_, err := DB.Exec("INSERT OR IGNORE INTO roles (name, description) VALUES (?, ?)", 
			role.Name, role.Description)
		if err != nil {
			return err
		}
	}

	// Add sample orders
	orders := []struct {
		UserID      int
		TotalAmount float64
		Status      string
	}{
		{1, 99.99, "completed"},
		{2, 149.95, "processing"},
		{3, 29.99, "pending"},
	}

	for _, order := range orders {
		result, err := DB.Exec("INSERT OR IGNORE INTO orders (user_id, total_amount, status) VALUES (?, ?, ?)", 
			order.UserID, order.TotalAmount, order.Status)
		if err != nil {
			return err
		}
		// Get the order ID for inserting order items
		orderID, _ := result.LastInsertId()
		
		// Add sample order items
		if orderID == 1 {
			_, err = DB.Exec("INSERT OR IGNORE INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)", 
				orderID, 1, 2, 49.99)
			if err != nil {
				return err
			}
		} else if orderID == 2 {
			_, err = DB.Exec("INSERT OR IGNORE INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)", 
				orderID, 2, 1, 149.95)
			if err != nil {
				return err
			}
		} else if orderID == 3 {
			_, err = DB.Exec("INSERT OR IGNORE INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)", 
				orderID, 3, 1, 29.99)
			if err != nil {
				return err
			}
		}
	}

	// Add sample addresses
	addresses := []struct {
		UserID      int
		StreetLine1 string
		StreetLine2 string
		City        string
		State       string
		PostalCode  string
		Country     string
		IsDefault   bool
	}{
		{1, "123 Main St", "Apt 4B", "New York", "NY", "10001", "USA", true},
		{2, "456 Oak Ave", "", "Los Angeles", "CA", "90001", "USA", true},
	}
	
	for _, address := range addresses {
		_, err := DB.Exec(`
			INSERT OR IGNORE INTO addresses 
			(user_id, street_line1, street_line2, city, state, postal_code, country, is_default) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			address.UserID, address.StreetLine1, address.StreetLine2, address.City, 
			address.State, address.PostalCode, address.Country, address.IsDefault)
		if err != nil {
			return err
		}
	}

	log.Println("Initial data seeded successfully")
	return nil
}