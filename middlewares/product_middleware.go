package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"go-crud/models"
)

type ProductWithStatus struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PaginatedProducts struct {
	Products    []ProductWithStatus `json:"products"`
	TotalCount  int                 `json:"total_count"`
	CurrentPage int                 `json:"current_page"`
	TotalPages  int                 `json:"total_pages"`
	PerPage     int                 `json:"per_page"`
	Filters     map[string]string   `json:"filters,omitempty"`
}

// FilterParams stores the filter parameters
type FilterParams struct {
	Name   string
	Status string
	MinID  int
	MaxID  int
}

func ProductMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Product Middleware is working...")

		// Get pagination parameters from query string
		page := 1
		perPage := 10 // Default items per page
		
		// Parse page parameter
		if pageParam := r.URL.Query().Get("page"); pageParam != "" {
			if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
				page = parsedPage
			}
		}
		
		// Parse perPage parameter
		if perPageParam := r.URL.Query().Get("per_page"); perPageParam != "" {
			if parsedPerPage, err := strconv.Atoi(perPageParam); err == nil && parsedPerPage > 0 {
				perPage = parsedPerPage
			}
		}
		
		// Parse filter parameters
		filters := FilterParams{}
		
		// Name filter (case-insensitive partial match)
		filters.Name = r.URL.Query().Get("name")
		
		// Status filter (exact match)
		filters.Status = r.URL.Query().Get("status")
		
		// ID range filters
		if minIDParam := r.URL.Query().Get("min_id"); minIDParam != "" {
			if parsedMinID, err := strconv.Atoi(minIDParam); err == nil {
				filters.MinID = parsedMinID
			}
		}
		
		if maxIDParam := r.URL.Query().Get("max_id"); maxIDParam != "" {
			if parsedMaxID, err := strconv.Atoi(maxIDParam); err == nil {
				filters.MaxID = parsedMaxID
			}
		}
		
		// Get filtered products with pagination
		paginatedProducts, err := GetFilteredProducts(page, perPage, filters)
		if err != nil {
			http.Error(w, "Error fetching products: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Add pagination data to context
		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, "paginatedProducts", paginatedProducts))

		// Call the next handler
		next(w, r)
	}
}

func GetFilteredProducts(page, perPage int, filters FilterParams) (PaginatedProducts, error) {
	offset := (page - 1) * perPage
	
	// Build WHERE clause based on filters
	whereClause, args := buildWhereClause(filters)
	
	// Count total products after filtering
	var totalCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	if err := models.DB.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return PaginatedProducts{}, err
	}

	// Calculate total pages
	totalPages := (totalCount + perPage - 1) / perPage

	// Query for filtered and paginated products
	query := fmt.Sprintf(`
		SELECT p.id, p.name, p.status
		FROM products p
		%s
		ORDER BY p.id
		LIMIT ? OFFSET ?
	`, whereClause)
	
	// Add pagination parameters to args
	queryArgs := append(args, perPage, offset)
	
	rows, err := models.DB.Query(query, queryArgs...)
	if err != nil {
		return PaginatedProducts{}, err
	}
	defer rows.Close()

	var products []ProductWithStatus
	for rows.Next() {
		var p ProductWithStatus
		if err := rows.Scan(&p.ID, &p.Name, &p.Status); err != nil {
			return PaginatedProducts{}, err
		}
		products = append(products, p)
	}
	
	// Create filters map for response
	filtersMap := make(map[string]string)
	if filters.Name != "" {
		filtersMap["name"] = filters.Name
	}
	if filters.Status != "" {
		filtersMap["status"] = filters.Status
	}
	if filters.MinID > 0 {
		filtersMap["min_id"] = strconv.Itoa(filters.MinID)
	}
	if filters.MaxID > 0 {
		filtersMap["max_id"] = strconv.Itoa(filters.MaxID)
	}

	return PaginatedProducts{
		Products:    products,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
		PerPage:     perPage,
		Filters:     filtersMap,
	}, nil
}

// buildWhereClause constructs SQL WHERE clause and arguments based on filters
func buildWhereClause(filters FilterParams) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	
	if filters.Name != "" {
		conditions = append(conditions, "LOWER(name) LIKE LOWER(?)")
		args = append(args, "%"+filters.Name+"%")
	}
	
	if filters.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filters.Status)
	}
	
	if filters.MinID > 0 {
		conditions = append(conditions, "id >= ?")
		args = append(args, filters.MinID)
	}
	
	if filters.MaxID > 0 {
		conditions = append(conditions, "id <= ?")
		args = append(args, filters.MaxID)
	}
	
	if len(conditions) > 0 {
		return "WHERE " + strings.Join(conditions, " AND "), args
	}
	
	return "", args
}