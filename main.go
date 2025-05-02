package main

import (
	"database/sql"
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
	routes.RegisterRoutes()
	
	serverAddr := ":" + config.AppConfig.ServerPort
	log.Println("Server running at http://localhost" + serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}