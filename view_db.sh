#!/bin/bash

DB_FILE="sqlite_db.db"  # Change this to your database file name

if [ ! -f "$DB_FILE" ]; then
    echo "Database file not found: $DB_FILE"
    exit 1
fi

echo "=== Database Tables ==="
sqlite3 $DB_FILE ".tables"

echo -e "\n=== Table Schemas ==="
for table in $(sqlite3 $DB_FILE ".tables"); do
    echo -e "\nTable: $table"
    sqlite3 $DB_FILE ".schema $table"
done

echo -e "\n=== Data Preview ==="
for table in $(sqlite3 $DB_FILE ".tables"); do
    echo -e "\nTable: $table"
    sqlite3 $DB_FILE "SELECT * FROM $table LIMIT 5;"
done