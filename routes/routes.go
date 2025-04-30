package routes

import (
	"net/http"
	"go-crud/controllers"
	"go-crud/middlewares"
)

func RegisterRoutes() {
	// Admin interface
	http.HandleFunc("/admin", controllers.AdminDashboard)

	// User routes
	http.HandleFunc("/users", controllers.GetUsers)
	http.HandleFunc("/users/create", controllers.CreateUser)
	http.HandleFunc("/users/update", controllers.UpdateUser)
	http.HandleFunc("/users/delete", controllers.DeleteUser)
	
	// Product routes
	http.HandleFunc("/products", middlewares.ProductMiddleware(controllers.GetProducts))
	
	// Role routes
	http.HandleFunc("/roles", controllers.GetRoles)
	http.HandleFunc("/roles/get", controllers.GetRoleByID)
	http.HandleFunc("/roles/create", controllers.CreateRole)
	http.HandleFunc("/roles/update", controllers.UpdateRole)
	http.HandleFunc("/roles/delete", controllers.DeleteRole)

	// Order routes
	http.HandleFunc("/orders", controllers.GetOrders)
	http.HandleFunc("/orders/get", controllers.GetOrderByID)
	http.HandleFunc("/orders/place", controllers.PlaceOrder)
	http.HandleFunc("/orders/update-status", controllers.UpdateOrderStatus)
	http.HandleFunc("/orders/delete", controllers.DeleteOrder)
	
	// Address routes - NEW
	http.HandleFunc("/addresses", controllers.GetAddresses)
	http.HandleFunc("/addresses/get", controllers.GetAddressByID)
	http.HandleFunc("/addresses/create", controllers.CreateAddress)
	http.HandleFunc("/addresses/update", controllers.UpdateAddress)
	http.HandleFunc("/addresses/delete", controllers.DeleteAddress)
	http.HandleFunc("/addresses/assign-to-order", controllers.AssignAddressToOrder)
	
	// For backward compatibility with the original API
	http.HandleFunc("/", controllers.GetUsers)
	http.HandleFunc("/post", controllers.CreateUser)
	http.HandleFunc("/update", controllers.UpdateUser)
	http.HandleFunc("/delete", controllers.DeleteUser)
}