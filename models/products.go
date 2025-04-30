package models

import (
	"time"
)

type Product struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Price     float64   `json:"price"`  // Add price field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetProducts(limit int) ([]Product, error) {
	rows, err := DB.Query("SELECT id, name, status, price, created_at, updated_at FROM products LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Status, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func GetProductByID(id int) (Product, error) {
	var product Product
	err := DB.QueryRow("SELECT id, name, status, price, created_at, updated_at FROM products WHERE id = ?", id).Scan(
		&product.ID, &product.Name, &product.Status, &product.Price, &product.CreatedAt, &product.UpdatedAt)
	return product, err
}

func CreateProduct(name string, status string, price float64) (int, error) {
	if status == "" {
		status = "active"
	}
	
	result, err := DB.Exec(
		"INSERT INTO products (name, status, price, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", 
		name, status, price)
	if err != nil {
		return 0, err
	}
	
	id, err := result.LastInsertId()
	return int(id), err
}

func UpdateProduct(id int, name string, status string, price float64) error {
	_, err := DB.Exec(
		"UPDATE products SET name = ?, status = ?, price = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", 
		name, status, price, id)
	return err
}

func DeleteProduct(id int) error {
	_, err := DB.Exec("DELETE FROM products WHERE id = ?", id)
	return err
}