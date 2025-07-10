package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/zaimnazif974/budgeting-BE/pkg/config"
	"github.com/zaimnazif974/budgeting-BE/pkg/models"
	"github.com/zaimnazif974/budgeting-BE/pkg/routes"
)

func main() {

	//Reading .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 1. Connect to database
	config.ConnectDatabase()
	defer config.CloseDatabase()

	// 2. Run migrations
	db := config.GetDB()
	if err := db.AutoMigrate(
		&models.Budget{},
		&models.User{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed!")

	// 3. Setup routes
	router := setupRoutes()

	// 4. Start server
	port := getEnv("APP_PORT", "")
	log.Printf("Server starting on port %s", port)

	// Setup graceful shutdown
	go func() {
		if err := http.ListenAndServe(":"+port, router); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")
}

func setupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", healthCheck).Methods("GET")

	// API routes with version prefix
	api := router.PathPrefix("/api/v1").Subrouter()

	// Add routes from /routes
	routes.BudgetRoutes(api)
	routes.AuthRoutes(api)

	return router
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()
	sqlDB, err := db.DB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		http.Error(w, "Database ping failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "healthy",
		"database": "connected",
	})
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
