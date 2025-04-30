package controllers

import (
	"encoding/json"
	"net/http"
	"go-crud/middlewares"
	"go-crud/models"
)

type productRequest struct {
	ID      int    `json:"id,omitempty"`
	Product string `json:"product,omitempty"`
}

type productResponse struct {
	Message string `json:"message"`
}

// GetProducts handles retrieving products with filters and pagination
func GetProducts(w http.ResponseWriter, r *http.Request) {
	// Get paginated products from context
	paginatedProducts, ok := r.Context().Value("paginatedProducts").(middlewares.PaginatedProducts)
	if !ok {
		// Fallback to direct DB query if middleware didn't work
		products, err := models.GetProducts(100)
		if err != nil {
			http.Error(w, "Error fetching products: "+err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(products)
		return
	}

	// Check if we have any products
	if len(paginatedProducts.Products) == 0 {
		// We still return the pagination structure with empty products array
		// This is better than returning an error, as it's a valid state
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paginatedProducts)
		return
	}

	// Return the paginated and filtered products
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedProducts)
}