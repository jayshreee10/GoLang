package routes

import (
	"net/http"

	"gorm-crud/controllers"
)

func RegisterRoutes() {
	http.HandleFunc("/", controllers.GetUsers)
	http.HandleFunc("/post", controllers.CreateUser)
	http.HandleFunc("/update", controllers.UpdateUser)
	http.HandleFunc("/delete", controllers.DeleteUser)
}
