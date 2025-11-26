package main

import (
	_ "challenge-app/docs"
	"challenge-app/internal/bootstrap"
	"challenge-app/internal/injector"
	"challenge-app/pkg/validation"
	"log"
	"os"
	"github.com/gin-contrib/cors"
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

	// Set up Gin router
	r := gin.Default()

	// Apply CORS middleware globally
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow requests from localhost:3000
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "ACCEPT"},
		AllowCredentials: true,
	}))

	// Example route for Swagger or health check endpoint
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "API is running",
		})
	})

	// Set up routes, inject your application, etc.
	// Example route:
	// app.Router.GET("/api/v1/users", getUserHandler)

	log.Printf("Starting server on port %s...", bootstrap.AppPort)

	// Start the server - use 0.0.0.0 for Docker
	if err := app.Router.Run("0.0.0.0:" + bootstrap.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
