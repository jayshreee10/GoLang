package main

import (
	"database/sql"
	"log"
	"net/http"
	_ "github.com/lib/pq"
	"go-crud/models"
	"go-crud/routes"
)

func main() {
	var err error
	connStr := "host=localhost port=5432 user=jayshree password=password dbname=spree_development sslmode=disable"
	models.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	if err := models.DB.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}

	log.Println("Connected to DB.")
	routes.RegisterRoutes()

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
