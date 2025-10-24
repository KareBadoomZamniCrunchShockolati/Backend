package bootstrap

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AppPort      string
	JWTSecretKey string
	DBHost       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBPort       string
	SSLMode      string
}

// LoadEnv loads variables from .env file and environment
func LoadEnv() *Env {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables.")
	}

	env := &Env{
		AppPort:      getEnv("APP_PORT", "8080"),
		JWTSecretKey: mustGetEnv("JWT_SECRET_KEY"),
		DBHost:       mustGetEnv("DB_HOST"),
		DBUser:       mustGetEnv("DB_USER"),
		DBPassword:   mustGetEnv("DB_PASSWORD"),
		DBName:       mustGetEnv("DB_NAME"),
		DBPort:       mustGetEnv("DB_PORT"),
		SSLMode:      getEnv("SSL_MODE", "disable"),
	}

	return env
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return value
}
