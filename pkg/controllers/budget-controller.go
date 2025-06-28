package controllers

import (
	"log"
	"net/http"

	"github.com/zaimnazif974/budgeting-BE/pkg/config"
	"github.com/zaimnazif974/budgeting-BE/pkg/models"
	"github.com/zaimnazif974/budgeting-BE/pkg/utils"
)

// CreateBudget creates a new budget
func CreateBudget(w http.ResponseWriter, r *http.Request) {
	var budget models.Budget
	utils.ParseBody(r, &budget)

	// Validate input
	if budget.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Budget Name must be filled")
		return
	}
	if budget.Amount <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "Budget Name must be filled")
		return
	}

	db := config.GetDB()
	if err := db.Create(&budget).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create budget")
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, budget, "Budget successfully created")
}

// GetBudgets without query
func GetBudgets(w http.ResponseWriter, r *http.Request) {
	var budgets []models.Budget

	db := config.GetDB()

	err := db.Find(&budgets).Error

	if err != nil {
		utils.WriteError(w, http.StatusBadGateway, "Couldn't get budgets")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, budgets, "Sucessfully fetch budgets")
}

// GetBudget by id
func GetBudgetByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var budget models.Budget

	utils.ParseBody(r, &id)

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "Couldn't get the id")
		return
	}

	db := config.GetDB()

	err := db.First(&budget, id).Error

	if err != nil {
		utils.WriteError(w, http.StatusBadGateway, "Couldn't get Budget")
		log.Printf("GetBudgetByID Error: %v", err)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, budget, "Sucessfully fetch budgets")
}
