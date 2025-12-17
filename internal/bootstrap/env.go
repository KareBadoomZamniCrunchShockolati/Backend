package bootstrap

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Env struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Security SecurityConfig
	Email    EmailConfig
	Storage  StorageConfig
}

type AppConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Address  string
	Port     string
	Password string
	DB       int
}

type SecurityConfig struct {
	JWTSecret string
	JWTIssuer string
	TokenTTL  time.Duration
}

type EmailConfig struct {
	From     string
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
}

type StorageConfig struct {
	Provider  string
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
	PublicURL string
}

func LoadEnv() *Env {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	return &Env{
		App: AppConfig{
			Port: AppPort,
			Mode: getEnv("APP_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     mustGetEnv("DB_HOST"),
			Port:     mustGetEnv("DB_PORT"),
			User:     mustGetEnv("DB_USER"),
			Password: mustGetEnv("DB_PASSWORD"),
			Name:     mustGetEnv("DB_NAME"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Address:  getEnv("REDIS_ADDR", "localhost"),
			Port:     RedisPort,
			Password: getEnv("REDIS_PASS", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Security: SecurityConfig{
			JWTSecret: mustGetEnv("JWT_SECRET_KEY"),
			JWTIssuer: JWTIssuer,
			TokenTTL:  JWTTokenExpiry,
		},
		Email: EmailConfig{
			From:     mustGetEnv("EMAIL_FROM"),
			SMTPHost: mustGetEnv("SMTP_HOST"),
			SMTPPort: getEnvInt("SMTP_PORT", 587),
			SMTPUser: mustGetEnv("SMTP_USER"),
			SMTPPass: mustGetEnv("SMTP_PASS"),
		},
		Storage: StorageConfig{
			Provider:  getEnv("STORAGE_PROVIDER", "arvan"),
			Endpoint:  mustGetEnv("STORAGE_ENDPOINT"),
			Region:    mustGetEnv("STORAGE_REGION"),
			Bucket:    mustGetEnv("STORAGE_BUCKET"),
			AccessKey: mustGetEnv("STORAGE_ACCESS_KEY"),
			SecretKey: mustGetEnv("STORAGE_SECRET_KEY"),
			PublicURL: mustGetEnv("STORAGE_PUBLIC_URL"),
		},
	}
}

// --- Helper functions ---
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

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func GetDSN(cfg *Env) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)
}
