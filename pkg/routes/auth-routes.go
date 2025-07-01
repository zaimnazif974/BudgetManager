package routes

import (
	"github.com/gorilla/mux"
	"github.com/zaimnazif974/budgeting-BE/pkg/controllers"
)

var AuthRoutes = func(router *mux.Router) {

	authRoutes := router.PathPrefix("/auth").Subrouter()

	authRoutes.HandleFunc("/google/callback", controllers.GoogleLogin).Methods("POST")
}
