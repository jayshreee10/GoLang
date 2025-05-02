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
	// Auth routes - no middleware required
	http.HandleFunc("/auth/register", controllers.RegisterUser)
	http.HandleFunc("/auth/login", controllers.LoginUser)
	
	// Protected auth route - requires auth
	http.HandleFunc("/auth/me", middlewares.AuthMiddleware(controllers.GetCurrentUser))
	// http.HandleFunc("/auth/change-password", middlewares.AuthMiddleware(controllers.ChangePassword))

	// Admin interface - protected by auth middleware
	http.HandleFunc("/admin", middlewares.AuthMiddleware(controllers.AdminDashboard))

	// User routes - protected
	http.HandleFunc("/users", middlewares.AuthMiddleware(controllers.GetUsers))
	http.HandleFunc("/users/profile", middlewares.AuthMiddleware(controllers.GetUserProfile)) // Changed to use renamed function
	http.HandleFunc("/users/create", middlewares.AuthMiddleware(controllers.CreateUser))
	http.HandleFunc("/users/update", middlewares.AuthMiddleware(controllers.UpdateUser))
	http.HandleFunc("/users/delete", middlewares.AuthMiddleware(controllers.DeleteUser))
	
	// Product routes - product browsing can be public
	http.HandleFunc("/products", middlewares.OptionalAuthMiddleware(middlewares.ProductMiddleware(controllers.GetProducts)))
	
	// Role routes - protected
	http.HandleFunc("/roles", middlewares.AuthMiddleware(controllers.GetRoles))
	http.HandleFunc("/roles/get", middlewares.AuthMiddleware(controllers.GetRoleByID))
	http.HandleFunc("/roles/create", middlewares.AuthMiddleware(controllers.CreateRole))
	http.HandleFunc("/roles/update", middlewares.AuthMiddleware(controllers.UpdateRole))
	http.HandleFunc("/roles/delete", middlewares.AuthMiddleware(controllers.DeleteRole))

	// Order routes - protected
	http.HandleFunc("/orders", middlewares.AuthMiddleware(controllers.GetOrders))
	http.HandleFunc("/orders/get", middlewares.AuthMiddleware(controllers.GetOrderByID))
	http.HandleFunc("/orders/place", middlewares.AuthMiddleware(controllers.PlaceOrder))
	http.HandleFunc("/orders/update-status", middlewares.AuthMiddleware(controllers.UpdateOrderStatus))
	http.HandleFunc("/orders/delete", middlewares.AuthMiddleware(controllers.DeleteOrder))
	
	// Address routes - protected
	http.HandleFunc("/addresses", middlewares.AuthMiddleware(controllers.GetAddresses))
	http.HandleFunc("/addresses/get", middlewares.AuthMiddleware(controllers.GetAddressByID))
	http.HandleFunc("/addresses/create", middlewares.AuthMiddleware(controllers.CreateAddress))
	http.HandleFunc("/addresses/update", middlewares.AuthMiddleware(controllers.UpdateAddress))
	http.HandleFunc("/addresses/delete", middlewares.AuthMiddleware(controllers.DeleteAddress))
	http.HandleFunc("/addresses/assign-to-order", middlewares.AuthMiddleware(controllers.AssignAddressToOrder))
	
	// For backward compatibility with the original API - deprecated
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/post", middlewares.AuthMiddleware(controllers.CreateUser))
	http.HandleFunc("/update", middlewares.AuthMiddleware(controllers.UpdateUser))
	http.HandleFunc("/delete", middlewares.AuthMiddleware(controllers.DeleteUser))
}