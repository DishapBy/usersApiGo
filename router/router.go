package router

import (
	"userApiGo/services"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/users", services.GetAllUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{id}", services.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newuser", services.CreateNewUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/updateuser/{id}", services.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteuser/{id}", services.DeleteUser).Methods("DELETE", "OPTIONS")

	return router
}
