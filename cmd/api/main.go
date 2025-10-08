package main

import (
	"fmt"
	"log"
	"os"

	"challenge-app/internal/handlers"
	"challenge-app/internal/repository/postgres"
	"challenge-app/internal/router"
	"challenge-app/internal/service"
	"challenge-app/internal/repository/models"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	gormPostgres "gorm.io/driver/postgres"
)

func main() {
	// --- 1. CONFIGURATION & DATABASE SETUP ---
	
	// Load environment variables from .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file. Ensure it exists in the project root.")
	}

	// Construct the Database Connection String (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("SSL_MODE"))

	db, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection successful.")

	// Run Migrations
	// NOTE: In a production app, use dedicated migration tools. AutoMigrate is for rapid development.
	err = db.AutoMigrate(&models.UserModel{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("User table migration complete.")

	// --- 2. DEPENDENCY INJECTION (WIRING THE LAYERS) ---
	
	// Repository Layer (Infrastructure): Instantiate GORM-backed repositories
	userRepo := postgres.NewUserRepository(db)
	
	// Service Layer (Business Logic): Instantiate services and inject repositories
	authService := service.NewAuthService(userRepo)
	
	// Handler Layer (Presentation): Instantiate handlers and inject services
	authHandler := handlers.NewAuthHandler(authService)

	// --- 3. RUN APPLICATION ---
	
	// Set up the Gin router and inject all necessary handlers
	r := router.SetupRouter(authHandler)

	// Determine the port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	
	// Start Server
	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}