# Product API Documentation

## Endpoint

```
GET /products
```

## Description

This endpoint retrieves a paginated and filterable list of products from the system.

## Query Parameters

### Pagination Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `page` | Integer | The page number to retrieve | 1 |
| `per_page` | Integer | Number of items per page | 10 |

### Filter Parameters

| Parameter | Type | Description | Filter Type |
|-----------|------|-------------|------------|
| `name` | String | Filter products by name | Case-insensitive partial match |
| `status` | String | Filter products by status | Exact match |
| `min_id` | Integer | Filter products with ID greater than or equal to value | Range |
| `max_id` | Integer | Filter products with ID less than or equal to value | Range |

## Response Structure

```json
{
  "products": [
    {
      "id": 1,
      "name": "Product Name",
      "status": "active"
    },
    ...
  ],
  "total_count": 100,
  "current_page": 1,
  "total_pages": 10,
  "per_page": 10,
  "filters": {
    "name": "product",
    "status": "active"
  }
}
```

## Examples

### Get all products (first page, 10 per page)
```
GET /products
```

### Get page 2 with 20 items per page
```
GET /products?page=2&per_page=20
```

### Get active products only
```
GET /products?status=active
```

### Get products with name containing "shirt"
```
GET /products?name=shirt
```

### Get products with ID between 100 and 200
```
GET /products?min_id=100&max_id=200
```

### Combine filters and pagination
```
GET /products?status=active&name=shirt&page=2&per_page=15
```

## Error Responses

| Status Code | Description |
|-------------|-------------|
| 400 | Bad Request - Indicates an invalid parameter format |
| 500 | Internal Server Error - Something went wrong on the server |

## Notes

- The `name` filter uses a LIKE query with wildcards on both sides, making it a partial match
- Filters can be combined for more specific queries
- If no products match the filter criteria, an empty array is returned with pagination metadata
- The response includes the applied filters in the `filters` field

# Get all products with pagination
http://localhost:8080/products

# Filter by name (case-insensitive partial match)
http://localhost:8080/products?name=shirt

# Filter by status (exact match)
http://localhost:8080/products?status=active

# Filter by ID range
http://localhost:8080/products?min_id=100&max_id=200

# Combining filters with pagination
http://localhost:8080/products?status=active&name=shirt&page=2&per_page=15

# Order API Documentation

This document describes the Order API endpoints added to the Go CRUD application.

## Database Structure

The order system consists of two tables:

1. **orders** table:
   - `id` - Primary key
   - `user_id` - Foreign key to users table
   - `total_amount` - Total order amount
   - `status` - Order status (pending, processing, shipped, completed, etc.)
   - `created_at` - Creation timestamp
   - `updated_at` - Last update timestamp

2. **order_items** table:
   - `id` - Primary key
   - `order_id` - Foreign key to orders table
   - `product_id` - Foreign key to products table
   - `quantity` - Quantity of the product
   - `price` - Price of the product at the time of order

## API Endpoints

### 1. Get All Orders

**Endpoint:** `GET /orders`

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "total_amount": 99.99,
    "status": "completed",
    "created_at": "2025-04-30T10:00:00Z",
    "updated_at": "2025-04-30T10:30:00Z",
    "items": [
      {
        "id": 1,
        "order_id": 1,
        "product_id": 1,
        "quantity": 2,
        "price": 49.99,
        "product": {
          "id": 1,
          "name": "Product 1",
          "status": "active"
        }
      }
    ]
  }
]
```

### 2. Get Order by ID

**Endpoint:** `GET /orders/get?id=1`

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "total_amount": 99.99,
  "status": "completed",
  "created_at": "2025-04-30T10:00:00Z",
  "updated_at": "2025-04-30T10:30:00Z",
  "items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": 1,
      "quantity": 2,
      "price": 49.99,
      "product": {
        "id": 1,
        "name": "Product 1",
        "status": "active"
      }
    }
  ]
}
```

### 3. Place New Order

**Endpoint:** `POST /orders/place`

**Request:**
```json
{
  "user_id": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    },
    {
      "product_id": 3,
      "quantity": 1
    }
  ]
}
```

**Response:**
```json
{
  "message": "Order placed successfully",
  "order_id": 4,
  "order": {
    "id": 4,
    "user_id": 1,
    "total_amount": 129.97,
    "status": "pending",
    "created_at": "2025-04-30T15:30:00Z",
    "updated_at": "2025-04-30T15:30:00Z",
    "items": [
      {
        "id": 4,
        "order_id": 4,
        "product_id": 1,
        "quantity": 2,
        "price": 49.99,
        "product": {
          "id": 1,
          "name": "Product 1",
          "status": "active"
        }
      },
      {
        "id": 5,
        "order_id": 4,
        "product_id": 3,
        "quantity": 1,
        "price": 29.99,
        "product": {
          "id": 3,
          "name": "Product 3",
          "status": "active"
        }
      }
    ]
  }
}
```

### 4. Update Order Status

**Endpoint:** `PUT /orders/update-status`

**Request:**
```json
{
  "id": 1,
  "status": "shipped"
}
```

**Response:**
```json
{
  "message": "Order status updated successfully"
}
```

### 5. Delete Order

**Endpoint:** `DELETE /orders/delete`

**Request:**
```json
{
  "id": 1
}
```

**Response:**
```json
{
  "message": "Order deleted successfully"
}
```

## Admin Dashboard

The admin dashboard at `/admin` now includes order and order item tables to view all data.

## Testing

You can use the `test_orders.sh` script to test the Order API:

```bash
chmod +x test_orders.sh
./test_orders.sh
```

This will send requests to all the order API endpoints and display the responses.