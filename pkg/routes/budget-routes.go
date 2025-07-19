package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zaimnazif974/budgeting-BE/pkg/controllers"
	"github.com/zaimnazif974/budgeting-BE/pkg/middlewares"
)

// AddBudgetRoutes adds budget-related routes to the router
var BudgetRoutes = func(router *mux.Router) {

	// Add prefix
	budgetRoutes := router.PathPrefix("/budget").Subrouter()

	// Create a new budget
	budgetRoutes.Handle("/create", middlewares.JWTMiddleware(http.HandlerFunc(controllers.CreateBudget))).Methods("POST")

	// get all budgets
	budgetRoutes.Handle("/", middlewares.JWTMiddleware(http.HandlerFunc(controllers.GetBudgets))).Methods("GET")

	// get budget by id
	budgetRoutes.Handle("/budget", middlewares.JWTMiddleware(http.HandlerFunc(controllers.GetBudgetByID))).Methods("GET")

	// Update budget by ID
	budgetRoutes.Handle("/budget/{id}/update", middlewares.JWTMiddleware(http.HandlerFunc(controllers.EditBudget))).Methods("PUT")

	budgetRoutes.Handle("/calendar", middlewares.JWTMiddleware(http.HandlerFunc(controllers.GetCalendar))).Methods("GET")

	// // Delete budget by ID
	// budgetRoutes.HandleFunc("/budget/{id}", controllers.DeleteBudget).Methods("DELETE")
}
