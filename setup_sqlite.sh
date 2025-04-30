#!/bin/bash

# Check if sqlite3 is installed
if command -v sqlite3 &> /dev/null; then
    echo "SQLite is already installed. Version:"
    sqlite3 --version
else
    echo "SQLite is not installed. Installing..."
    
    # Check the operating system
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        sudo apt-get update
        sudo apt-get install -y sqlite3
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        brew install sqlite3
    else
        echo "Please install SQLite manually for your operating system."
        exit 1
    fi
    
    echo "SQLite installed successfully. Version:"
    sqlite3 --version
fi

# Add Go SQLite driver to the project
echo "Adding SQLite driver to Go project..."
cd "$(dirname "$0")"  # Navigate to script directory
go get github.com/mattn/go-sqlite3

echo "Installation and setup complete!"