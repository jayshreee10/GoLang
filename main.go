package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	
	_ "github.com/mattn/go-sqlite3"
	"go-crud/models"
	"go-crud/routes"
	"go-crud/config"
)

func main() {
	// Load configuration from environment variables
	config.Initialize()
	
	// Log admin credentials for development purposes
	log.Printf("Admin credentials set to - Email: %s, Password: %s", 
		config.AppConfig.DefaultAdminEmail, 
		config.AppConfig.DefaultAdminPassword)
	log.Println("IMPORTANT: Change these credentials in production using environment variables!")
	
	dbPath := config.AppConfig.DBPath

	needInit := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		needInit = true
	}

	var err error
	models.DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to SQLite DB:", err)
	}
	defer models.DB.Close()

	if err := models.DB.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}
	
	if needInit {
		if err := models.InitDB(); err != nil {
			log.Fatal("Failed to initialize database:", err)
		}
		log.Println("Database initialized successfully.")
	}
	
	log.Println("Connected to SQLite DB at", dbPath)
	
	// Register API routes
	routes.RegisterRoutes()
	
	// Start HTTP server
	serverAddr := ":" + config.AppConfig.ServerPort
	fmt.Printf("\n✓ Server is running!\n")
	fmt.Printf("✓ Local:   http://localhost%s\n\n", serverAddr)
	fmt.Println("Authentication required for all endpoints:")
	fmt.Printf("✓ Username: %s\n", config.AppConfig.DefaultAdminEmail)
	fmt.Printf("✓ Password: %s\n\n", config.AppConfig.DefaultAdminPassword)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}