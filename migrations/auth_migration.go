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
	fmt.Println("Starting auth migration...")

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}
	defer tx.Rollback() // Will be ignored if transaction is committed

	// Check if password column exists in users table
	var passwordExists int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='password'").Scan(&passwordExists)
	if err != nil {
		log.Fatal("Failed to check if password column exists:", err)
	}
	
	if passwordExists == 0 {
		fmt.Println("Adding password column to users table...")
		_, err = tx.Exec("ALTER TABLE users ADD COLUMN password TEXT")
		if err != nil {
			log.Fatal("Failed to add password column:", err)
		}
		fmt.Println("Password column added successfully.")
	} else {
		fmt.Println("Password column already exists.")
	}

	// Create refresh_tokens table if it doesn't exist
	fmt.Println("Creating refresh_tokens table if it doesn't exist...")
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL UNIQUE,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatal("Failed to create refresh_tokens table:", err)
	}
	fmt.Println("Refresh tokens table created successfully.")

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}

	fmt.Println("Auth migration completed successfully.")
}