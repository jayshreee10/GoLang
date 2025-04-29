package routes

import (
	"net/http"
	"go-crud/controllers"
	"go-crud/middlewares"
)

func RegisterRoutes() {
	http.HandleFunc("/", controllers.GetUsers)
	http.HandleFunc("/products", middlewares.ProductStatusMiddleware(controllers.GetProducts))
	http.HandleFunc("/post", controllers.CreateUser)
	http.HandleFunc("/update", controllers.UpdateUser)
	http.HandleFunc("/delete", controllers.DeleteUser)
}
