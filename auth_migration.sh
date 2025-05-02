#!/bin/bash

echo "Running authentication system migration..."

# Compile and run the migration Go code
go run migrations/auth_migration.go

# Install necessary dependencies
echo "Installing dependencies..."
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt

echo "Migration completed!"
echo "You can now use the authentication system with your Go application."