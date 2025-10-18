package main

import (
	"challenge-app/internal/application/service"
	"challenge-app/internal/infrastructure/repository/postgres"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"challenge-app/internal/presentation/handler"
	"challenge-app/internal/presentation/middleware"
	"challenge-app/internal/presentation/router"
	"challenge-app/pkg/security"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func validateEnvVars() error {
	required := []string{
		"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"JWT_SECRET_KEY", "SMTP_HOST", "SMTP_USER", "SMTP_PASS",
		"EMAIL_FROM",
	}

	for _, env := range required {
		if os.Getenv(env) == "" {
			return fmt.Errorf("required environment variable %s is not set", env)
		}
	}
	return nil
}

func main() {
	// --- 1. CONFIGURATION & DATABASE SETUP ---

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file. Ensure it exists in the project root.")
	}

	// Validate required environment variables
	if err := validateEnvVars(); err != nil {
		log.Fatalf("Environment configuration error: %v", err)
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
	err = db.AutoMigrate(&entity.UserEntity{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("User table migration complete.")

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connection successful.")

	// --- 2. DEPENDENCY INJECTION (WIRING THE LAYERS) ---
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("FATAL: JWT_SECRET_KEY is not set in environment. Please set it.")
	}
	const tokenExpiry = time.Hour * 24
	jwtService := security.NewJWTService(jwtSecretKey, tokenExpiry)

	// Repository Layer
	userRepo := postgres.NewUserRepository(db)

	// Service Layer
	authService := service.NewAuthService(userRepo, jwtService, rdb)
	userService := service.NewUserService(userRepo)

	jwtMiddleware := middleware.JWTAuthMiddleware(jwtService)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// --- 3. RUN APPLICATION ---
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
