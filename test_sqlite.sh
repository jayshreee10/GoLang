#!/bin/bash

echo "Testing Order CRUD API..."

# Start the server in the background (comment this out if your server is already running)
go run main.go &
SERVER_PID=$!

# Give the server a moment to start
sleep 2

# Base URL
BASE_URL="http://localhost:8080"

# Function to make HTTP requests and print response
function make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    echo -e "\n=== $method $endpoint ==="
    if [ -z "$data" ]; then
        response=$(curl -s -X $method $BASE_URL$endpoint)
    else
        response=$(curl -s -X $method -H "Content-Type: application/json" -d "$data" $BASE_URL$endpoint)
    fi
    
    echo $response | jq '.' 2>/dev/null || echo $response
}

# Test order endpoints
echo -e "\n\n=== Testing Order Endpoints ==="

# Get all orders
make_request "GET" "/orders"

# Place a new order
make_request "POST" "/orders/place" '{
  "user_id": 1,
  "items": [
    {"product_id": 1, "quantity": 2},
    {"product_id": 3, "quantity": 1}
  ]
}'

# Get all orders again to see the new order
make_request "GET" "/orders"

# Get a specific order (adjust ID as needed based on your database)
make_request "GET" "/orders/get?id=1"

# Update order status
make_request "PUT" "/orders/update-status" '{
  "id": 1,
  "status": "shipped"
}'

# Get the updated order
make_request "GET" "/orders/get?id=1"

# Visit the admin dashboard to see all data
echo -e "\n\n=== Admin Dashboard ==="
echo "Open http://localhost:8080/admin in your browser to see all data"

# Clean up - if you started the server in this script
echo -e "\n\nStopping server..."
kill $SERVER_PID 2>/dev/null

echo -e "\n\nTest completed."