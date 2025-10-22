package main

import (
	_ "challenge-app/docs" // Swagger docs
	"challenge-app/internal/bootstrap"
	"challenge-app/internal/infrastructure/repository/postgres/driver"
	"challenge-app/internal/injector"
	"fmt"
	"log"
	"time"
)

// @title Challenge App API
// @version 1.0
// @description Backend API for the Challenge App project built with Go and Gin.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/mobina-hoshiaripour/challenge-app
// @contact.email mobina@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// --- 1. LOAD CONFIG ---
	env := bootstrap.LoadEnv()

	// --- 2. SETUP DATABASE (later will be moved to a constructor) ---
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		env.DBHost, env.DBUser, env.DBPassword,
		env.DBName, env.DBPort, env.SSLMode,
	)

	db := driver.InitPostgresDB(dsn) // (next step we’ll create this function)
	log.Println("Database connected and migrated successfully.")

	// --- 3. DEPENDENCY INJECTION ---
	const tokenExpiry = time.Hour * 24
	r, err := injector.InitializeRouter(db, env.JWTSecretKey, tokenExpiry)
	if err != nil {
		log.Fatalf("Error initializing application dependencies: %v", err)
	}

	// Register custom validators (password etc.)
	if err := bootstrap.SetupValidator(); err != nil {
		log.Fatalf("failed to setup validators: %v", err)
	}
	// --- 4. START SERVER ---
	log.Printf("Starting server on port %s...", env.AppPort)
	if err := r.Run(":" + env.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
