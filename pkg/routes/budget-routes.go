package routes

import (
	"github.com/gorilla/mux"
	"github.com/zaimnazif974/budgeting-BE/pkg/controllers"
)

// AddBudgetRoutes adds budget-related routes to the router
var BudgetRoutes = func(router *mux.Router) {

	// Add prefix
	budgetRoutes := router.PathPrefix("/budget").Subrouter()

	// Create a new budget
	budgetRoutes.HandleFunc("/create", controllers.CreateBudget).Methods("POST")

	// get all budgets
	budgetRoutes.HandleFunc("/", controllers.GetBudgets).Methods("GET")

	// get budget by id
	budgetRoutes.HandleFunc("/budget", controllers.GetBudgetByID).Methods("GET")

	// Update budget by ID
	budgetRoutes.HandleFunc("/budget/{id}/update", controllers.EditBudget).Methods("PUT")

	// // Delete budget by ID
	// budgetRoutes.HandleFunc("/budget/{id}", controllers.DeleteBudget).Methods("DELETE")
}
