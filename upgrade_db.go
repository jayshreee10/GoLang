package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := "./sqlite_db.db"

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatal("Database file does not exist. Please run the main application first.")
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to SQLite DB:", err)
	}
	defer db.Close()

	fmt.Println("Connected to SQLite database successfully.")
	fmt.Println("Starting schema upgrade...")

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}
	defer tx.Rollback() // Will be ignored if transaction is committed

	// Create addresses table
	_, err = tx.Exec(`
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
		log.Fatal("Failed to create addresses table:", err)
	}
	fmt.Println("Addresses table created successfully.")
	
	// Add address_id column to orders table if it doesn't exist
	var addressIDExists int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('orders') WHERE name='address_id'").Scan(&addressIDExists)
	if err != nil {
		log.Fatal("Failed to check if address_id column exists:", err)
	}
	
	if addressIDExists == 0 {
		_, err = tx.Exec(`
		ALTER TABLE orders 
		ADD COLUMN address_id INTEGER DEFAULT NULL
		REFERENCES addresses(id)`)
		if err != nil {
			log.Fatal("Failed to add address_id column to orders table:", err)
		}
		fmt.Println("Added address_id column to orders table.")
	} else {
		fmt.Println("Address_id column already exists in orders table.")
	}
	
	// Add sample addresses for testing
	_, err = tx.Exec(`
	INSERT INTO addresses 
	(user_id, street_line1, street_line2, city, state, postal_code, country, is_default)
	SELECT 1, '123 Main St', 'Apt 4B', 'New York', 'NY', '10001', 'USA', 1
	WHERE NOT EXISTS (
		SELECT 1 FROM addresses WHERE user_id = 1 AND street_line1 = '123 Main St'
	)`)
	if err != nil {
		log.Fatal("Failed to insert first sample address:", err)
	}
	
	_, err = tx.Exec(`
	INSERT INTO addresses 
	(user_id, street_line1, street_line2, city, state, postal_code, country, is_default)
	SELECT 2, '456 Oak Ave', '', 'Los Angeles', 'CA', '90001', 'USA', 1
	WHERE NOT EXISTS (
		SELECT 1 FROM addresses WHERE user_id = 2 AND street_line1 = '456 Oak Ave'
	)`)
	if err != nil {
		log.Fatal("Failed to insert second sample address:", err)
	}
	fmt.Println("Sample addresses added.")

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}

	fmt.Println("Database schema upgrade completed successfully!")
	fmt.Println("You can now run your application with: go run main.go")
}