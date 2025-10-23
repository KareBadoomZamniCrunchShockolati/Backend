package main

import (
	_ "challenge-app/docs" // Swagger docs
	"challenge-app/internal/bootstrap"
	"challenge-app/internal/injector"
	"log"
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
	env := bootstrap.LoadEnv()

	app, err := injector.InitializeApplication()
	if err != nil {
		log.Fatalf("Error initializing dependencies: %v", err)
	}

	log.Printf("Starting server on port %s...", env.AppPort)
	if err := app.Router.Run(":" + env.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
