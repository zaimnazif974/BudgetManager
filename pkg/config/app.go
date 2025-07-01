package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnv("DB_PORT", "15001"),
		User:     GetEnv("POSTGRES_USER", "postgres"),
		Password: GetEnv("POSTGRES_PASSWORD", ""),
		DBName:   GetEnv("POSTGRES_DB", "budgeting"),
		SSLMode:  GetEnv("DB_SSLMODE", "disable"),
	}
}

// ConnectDatabase establishes connection to PostgreSQL database
func ConnectDatabase() {
	config := LoadDatabaseConfig()

	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	// Connect to database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Database connected successfully!")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting database instance: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("Database connection closed successfully")
	}
}

// GetEnv gets environment variable with fallback to default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
