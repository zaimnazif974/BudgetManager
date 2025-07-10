package routes

import (
	"github.com/gorilla/mux"
	"github.com/zaimnazif974/budgeting-BE/pkg/controllers"
)

var AuthRoutes = func(router *mux.Router) {

	authRoutes := router.PathPrefix("/auth").Subrouter()

	authRoutes.HandleFunc("/signup", controllers.SignUp).Methods("POST")

	authRoutes.HandleFunc("/login", controllers.Login).Methods("POST")

	authRoutes.HandleFunc("/{provider}/callback", controllers.GoogleLogin).Methods("GET")

	authRoutes.HandleFunc("/{provider}", controllers.BeginGoogleAuth).Methods("GET")

	authRoutes.HandleFunc("/{provider}/logout", controllers.Logout).Methods("GET")
}
