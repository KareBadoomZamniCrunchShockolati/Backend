package main

import (
	_ "challenge-app/docs"
	"challenge-app/internal/bootstrap"
	"challenge-app/internal/injector"
	"challenge-app/pkg/validation"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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
	// Load environment based on APP_MODE
	if os.Getenv("APP_MODE") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Register custom validation if necessary
	if err := validation.RegisterGinValidator(); err != nil {
		log.Fatalf("Failed to register custom validator: %v", err)
	}

	// Initialize application
	app, err := injector.InitializeApplication()
	if err != nil {
		log.Fatalf("Error initializing dependencies: %v", err)
	}

	log.Printf("Starting server on port %s...", bootstrap.AppPort)

	// Start the server - use 0.0.0.0 for Docker
	if err := app.Router.Run("0.0.0.0:" + bootstrap.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
