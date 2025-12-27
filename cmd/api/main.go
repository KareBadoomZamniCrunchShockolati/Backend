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

	// Initialize application with worker container
	container, err := injector.InitializeApplicationWithWorker()
	if err != nil {
		log.Fatalf("Error initializing dependencies: %v", err)
	}

	// Start the background worker
	container.Worker.Start()
	defer container.Worker.Stop()

	log.Println("✅ Background worker for challenge completion started")

	// Set up Gin router - use the router from the app
	r := container.App.Router

	// Apply CORS middleware globally (already applied in router setup, but adding here for safety)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "ACCEPT"},
		AllowCredentials: true,
	}))

	// Worker health endpoint
	r.GET("/api/v1/worker/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "worker_running",
			"message": "Challenge completion worker is active",
		})
	})

	log.Printf("Starting server on port %s...", bootstrap.AppPort)

	// Start the server
	if err := r.Run("0.0.0.0:" + bootstrap.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
