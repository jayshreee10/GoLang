package routes

import (
	"net/http"
	"go-crud/controllers"
	"go-crud/middlewares"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("hello"))
}

func RegisterRoutes() {
	// Auth routes - using AdminAuthMiddleware
	http.HandleFunc("/auth/register", middlewares.AdminAuthMiddleware(controllers.RegisterUser))
	http.HandleFunc("/auth/login", middlewares.AdminAuthMiddleware(controllers.LoginUser))
	
	// Protected auth route - requires both admin auth and JWT/session auth
	http.HandleFunc("/auth/me", middlewares.AdminAuthMiddleware(middlewares.AuthMiddleware(controllers.GetCurrentUser)))

	// Admin interface - protected by both admin auth and JWT auth
	http.HandleFunc("/admin", middlewares.AdminAuthMiddleware(middlewares.AuthMiddleware(controllers.AdminDashboard)))

	// User routes - protected by admin auth
	http.HandleFunc("/users", middlewares.AdminAuthMiddleware(controllers.GetUsers))
	http.HandleFunc("/users/profile", middlewares.AdminAuthMiddleware(controllers.GetUserProfile))
	http.HandleFunc("/users/create", middlewares.AdminAuthMiddleware(controllers.CreateUser))
	http.HandleFunc("/users/update", middlewares.AdminAuthMiddleware(controllers.UpdateUser))
	http.HandleFunc("/users/delete", middlewares.AdminAuthMiddleware(controllers.DeleteUser))
	
	// Product routes - protected by admin auth
	http.HandleFunc("/products", middlewares.AdminAuthMiddleware(
		middlewares.OptionalAuthMiddleware(
			middlewares.ProductMiddleware(controllers.GetProducts))))
	
	// Role routes - protected by admin auth
	http.HandleFunc("/roles", middlewares.AdminAuthMiddleware(controllers.GetRoles))
	http.HandleFunc("/roles/get", middlewares.AdminAuthMiddleware(controllers.GetRoleByID))
	http.HandleFunc("/roles/create", middlewares.AdminAuthMiddleware(controllers.CreateRole))
	http.HandleFunc("/roles/update", middlewares.AdminAuthMiddleware(controllers.UpdateRole))
	http.HandleFunc("/roles/delete", middlewares.AdminAuthMiddleware(controllers.DeleteRole))

	// Order routes - protected by admin auth
	http.HandleFunc("/orders", middlewares.AdminAuthMiddleware(controllers.GetOrders))
	http.HandleFunc("/orders/get", middlewares.AdminAuthMiddleware(controllers.GetOrderByID))
	http.HandleFunc("/orders/place", middlewares.AdminAuthMiddleware(controllers.PlaceOrder))
	http.HandleFunc("/orders/update-status", middlewares.AdminAuthMiddleware(controllers.UpdateOrderStatus))
	http.HandleFunc("/orders/delete", middlewares.AdminAuthMiddleware(controllers.DeleteOrder))
	
	// Address routes - protected by admin auth
	http.HandleFunc("/addresses", middlewares.AdminAuthMiddleware(controllers.GetAddresses))
	http.HandleFunc("/addresses/get", middlewares.AdminAuthMiddleware(controllers.GetAddressByID))
	http.HandleFunc("/addresses/create", middlewares.AdminAuthMiddleware(controllers.CreateAddress))
	http.HandleFunc("/addresses/update", middlewares.AdminAuthMiddleware(controllers.UpdateAddress))
	http.HandleFunc("/addresses/delete", middlewares.AdminAuthMiddleware(controllers.DeleteAddress))
	http.HandleFunc("/addresses/assign-to-order", middlewares.AdminAuthMiddleware(controllers.AssignAddressToOrder))
	
	// For backward compatibility with the original API - deprecated but still protected
	http.HandleFunc("/", middlewares.AdminAuthMiddleware(helloHandler))
	http.HandleFunc("/post", middlewares.AdminAuthMiddleware(controllers.CreateUser))
	http.HandleFunc("/update", middlewares.AdminAuthMiddleware(controllers.UpdateUser))
	http.HandleFunc("/delete", middlewares.AdminAuthMiddleware(controllers.DeleteUser))
}