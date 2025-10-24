package driver

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgresDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB instance: %v", err)
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	fmt.Println("PostgreSQL connected successfully.")
	return db, nil
}
