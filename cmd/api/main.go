package main

import (
	"challenge-app/internal/application/service"
	"challenge-app/internal/infrastructure/repository/postgres"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	redisRepo "challenge-app/internal/infrastructure/repository/redis"
	"challenge-app/internal/presentation/handler"
	"challenge-app/internal/presentation/middleware"
	"challenge-app/internal/presentation/router"
	"challenge-app/pkg/email"
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
		log.Println("Warning: .env file not found. Using system environment variables.")
	}

	// Validate required environment variables
	if err := validateEnvVars(); err != nil {
		log.Fatalf("Environment configuration error: %v", err)
	}

	// Construct the Database Connection String (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("SSL_MODE"))

	// Connect to PostgreSQL
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
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connection successful.")

	// --- 2. DEPENDENCY INJECTION (WIRING THE LAYERS) ---

	// JWT Service
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("FATAL: JWT_SECRET_KEY is not set in environment. Please set it.")
	}

	const tokenExpiry = time.Hour * 24
	jwtService := security.NewJWTService(jwtSecretKey, tokenExpiry)

	// Repository Layer
	userRepo := postgres.NewUserRepository(db)
	verificationRepo := redisRepo.NewVerificationRepository(rdb)

	// External Services
	emailService := email.NewEmailService()

	// Service Layer with proper dependencies
	authService := service.NewAuthService(
		userRepo,         // UserRepository
		verificationRepo, // VerificationRepository
		emailService,     // EmailService
		jwtService,       // JWTService
	)
	userService := service.NewUserService(userRepo, emailService, verificationRepo)

	// Middleware
	jwtMiddleware := middleware.JWTAuthMiddleware(jwtService)

	// Handler Layer
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
	log.Printf("API endpoints available at: http://localhost:%s/api/v1", port)
	log.Println("Available routes:")
	log.Println("  POST   /api/v1/auth/signup")
	log.Println("  POST   /api/v1/auth/login")
	log.Println("  POST   /api/v1/verify")
	log.Println("  POST   /api/v1/resend-verification")
	log.Println("  GET    /api/v1/users/profile (protected)")
	log.Println("  PUT    /api/v1/users/profile (protected)")
	log.Println("  DELETE /api/v1/users/profile (protected)")
	log.Println("  GET    /api/v1/users (protected)")

	if err := r.Run(":" + port); err != nil {
	_ "challenge-app/docs" // Swagger docs
	"challenge-app/internal/bootstrap"
	"challenge-app/internal/injector"
	"challenge-app/pkg/validation"
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

	if err := validation.RegisterGinValidator(); err != nil {
	log.Fatalf("Failed to register custom validator: %v", err)
	}

	app, err := injector.InitializeApplication()
	if err != nil {
		log.Fatalf("Error initializing dependencies: %v", err)
	}

	log.Printf("Starting server on port %s...", env.AppPort)
	if err := app.Router.Run(":" + env.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
