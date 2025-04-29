package middlewares

import (
	"context"
	"log"
	"net/http"
	"go-crud/models"
)


type ProductWithStatus struct {
	ID           int    `json:"id"`
	Name string `json:"name"`
	Status       string `json:"status"`
}


func ProductStatusMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Product Status Middleware is working...")

		
		productsWithStatus, err := GetActiveProducts()
		if err != nil {
			http.Error(w, "Error fetching active products: "+err.Error(), http.StatusInternalServerError)
			return
		}

		
		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, "products", productsWithStatus))

		
		next(w, r)
	}
}


func GetActiveProducts() ([]ProductWithStatus, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.status
		FROM spree_products p
		WHERE p.status = 'active'
		ORDER BY p.id
		LIMIT 100
	`


	rows, err := models.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	var activeProducts []ProductWithStatus
	for rows.Next() {
		var p ProductWithStatus
	
		if err := rows.Scan(&p.ID, &p.Name, &p.Status); err != nil {
			return nil, err
		}
		activeProducts = append(activeProducts, p)
	}


	return activeProducts, nil
}
