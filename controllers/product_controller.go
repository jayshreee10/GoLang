package controllers
import (
	"encoding/json"
	"net/http"
	"go-crud/models"
)

type productRequest struct {
	ID      int    `json:"id,omitempty"`
	Product string `json:"product,omitempty"`
}

type productResponse struct {
	Message string `json:"message"`
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := models.GetProducts(100)
	if err != nil {
		http.Error(w, "Error fetching products: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}