#!/bin/bash

echo "Altering existing SQLite database schema..."

DB_FILE="./sqlite_db.db"

# Check if the database file exists
if [ ! -f "$DB_FILE" ]; then
    echo "Database file not found: $DB_FILE"
    echo "Please run the application first to create the database."
    exit 1
fi

# Add price column to products table if it doesn't exist
echo "Adding price column to products table..."
sqlite3 $DB_FILE <<EOF
PRAGMA foreign_keys=off;
BEGIN TRANSACTION;

-- Check if price column exists and add it if it doesn't
SELECT CASE 
    WHEN COUNT(*) = 0 THEN
        'ALTER TABLE products ADD COLUMN price REAL DEFAULT 0.0;'
    ELSE
        'SELECT 1;'
END AS sql_to_run
FROM pragma_table_info('products') WHERE name='price'
|sqlite3 $DB_FILE;

-- Update existing products with sample prices
UPDATE products SET price = 49.99 WHERE id = 1 AND price IS NULL;
UPDATE products SET price = 149.95 WHERE id = 2 AND price IS NULL;
UPDATE products SET price = 29.99 WHERE id = 3 AND price IS NULL;

COMMIT;
PRAGMA foreign_keys=on;
EOF

# Create orders table if it doesn't exist
echo "Creating orders table if it doesn't exist..."
sqlite3 $DB_FILE <<EOF
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    total_amount REAL NOT NULL,
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
EOF

# Create order_items table if it doesn't exist
echo "Creating order_items table if it doesn't exist..."
sqlite3 $DB_FILE <<EOF
CREATE TABLE IF NOT EXISTS order_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price REAL NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);
EOF

# Add sample orders if orders table is empty
echo "Adding sample orders if the orders table is empty..."
sqlite3 $DB_FILE <<EOF
INSERT INTO orders (user_id, total_amount, status)
SELECT 1, 99.99, 'completed'
WHERE NOT EXISTS (SELECT 1 FROM orders LIMIT 1);

INSERT INTO orders (user_id, total_amount, status)
SELECT 2, 149.95, 'processing'
WHERE NOT EXISTS (SELECT 1 FROM orders LIMIT 1 OFFSET 1);

INSERT INTO orders (user_id, total_amount, status)
SELECT 3, 29.99, 'pending'
WHERE NOT EXISTS (SELECT 1 FROM orders LIMIT 1 OFFSET 2);
EOF

# Add sample order items if order_items table is empty
echo "Adding sample order items if the order_items table is empty..."
sqlite3 $DB_FILE <<EOF
INSERT INTO order_items (order_id, product_id, quantity, price)
SELECT 1, 1, 2, 49.99
WHERE EXISTS (SELECT 1 FROM orders WHERE id = 1)
AND NOT EXISTS (SELECT 1 FROM order_items WHERE order_id = 1);

INSERT INTO order_items (order_id, product_id, quantity, price)
SELECT 2, 2, 1, 149.95
WHERE EXISTS (SELECT 1 FROM orders WHERE id = 2)
AND NOT EXISTS (SELECT 1 FROM order_items WHERE order_id = 2);

INSERT INTO order_items (order_id, product_id, quantity, price)
SELECT 3, 3, 1, 29.99
WHERE EXISTS (SELECT 1 FROM orders WHERE id = 3)
AND NOT EXISTS (SELECT 1 FROM order_items WHERE order_id = 3);
EOF

echo "Database schema has been updated."
echo "You can now run your application with: go run main.go"