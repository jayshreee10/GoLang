// Package config handles environment configuration for the application
package config

import (
	"log"
	"os"
	"strconv"
	"time"
	
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	ServerPort string
	
	// Database configuration
	DBPath string
	
	// JWT configuration
	JWTSecret       string
	JWTExpiration   time.Duration // in hours
	
	// Default admin account
	DefaultAdminEmail    string
	DefaultAdminPassword string
}

// Global application configuration
var AppConfig Config

// Initialize loads configuration from environment variables
func Initialize() {
	// Try to load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
	
	// Initialize with defaults
	AppConfig = Config{
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		DBPath:              getEnv("DB_PATH", "./sqlite_db.db"),
		JWTSecret:           getEnv("JWT_SECRET", "your-default-secret-key-for-development-only"),
		JWTExpiration:       time.Duration(getEnvAsInt("JWT_EXPIRATION_HOURS", 24)) * time.Hour,
		DefaultAdminEmail:   getEnv("DEFAULT_ADMIN_EMAIL", "admin@example.com"),
		DefaultAdminPassword: getEnv("DEFAULT_ADMIN_PASSWORD", "admin123"),
	}
	
	log.Println("Configuration loaded successfully")
}

// Helper function to get an environment variable or a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get an environment variable as an integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %d\n", key, defaultValue)
		return defaultValue
	}
	
	return value
}