package main

import (
	"fmt"
	"log"
	"os"

	"challenge-app/internal/presentation/handler"
	"challenge-app/internal/presentation/middleware"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/presentation/router"
	"challenge-app/internal/application/service"
	"challenge-app/internal/infrastructure/repository/postgres"

	"github.com/joho/godotenv"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	err = db.AutoMigrate(&model.UserModel{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("User table migration complete.")

	// --- 2. DEPENDENCY INJECTION (WIRING THE LAYERS) ---
	
	// Repository Layer
	userRepo := postgres.NewUserRepository(db)
	
	// Service Layer
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	userHandler = handler.NewUserHandler(userService)

	// --- 3. RUN APPLICATION ---
	
	// Instantiate the JWT Middleware
	jwtMiddleware := middleware.JWTAuthMiddleware()

	// Set up the Gin router and inject all necessary handlers and the middleware
	r := router.SetupRouter(userHandler, authHandler, jwtMiddleware)

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
