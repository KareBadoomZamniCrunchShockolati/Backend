package driver

import (
	"fmt"
	"log"

	"challenge-app/internal/infrastructure/repository/postgres/entity"

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

	// Auto-migrate all entities
	if err := db.AutoMigrate(
		&entity.UserEntity{},
		&entity.FollowEntity{},
		&entity.ChallengeEntity{},
		&entity.ChallengeParticipantEntity{},
		&entity.CommentEntity{},
		&entity.ChallengeCategoryEntity{},
		&entity.ChallengeParticipationRequestEntity{},
		&entity.LikeEntity{},
		&entity.UserDailyFeelingEntity{},
		&entity.UserDailyGoalProgressEntity{},
		&entity.UserDailyNoteEntity{},
		&entity.PostEntity{},
		&entity.ChallengeCompletionEntity{},
	); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
		return nil, err
	}
	fmt.Println("Database migrated successfully.")

	// Create unique constraints after migration
	if err := CreateUniqueConstraints(db); err != nil {
		log.Fatalf("Failed to create unique constraints: %v", err)
		return nil, err
	}
	fmt.Println("Unique constraints created successfully.")

	return db, nil
}

func CreateUniqueConstraints(db *gorm.DB) error {
	// Unique constraint for likes - prevent duplicate likes
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_likes_unique 
		ON likes (entity_type, entity_id, user_id)
	`).Error; err != nil {
		return fmt.Errorf("failed to create likes unique index: %w", err)
	}

	// Optional: Add other unique constraints you might need
	// For example, prevent users from following the same person multiple times
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_follows_unique 
		ON follows (follower_id, following_id)
	`).Error; err != nil {
		return fmt.Errorf("failed to create follows unique index: %w", err)
	}

	// Optional: Prevent duplicate challenge participants
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_challenge_participants_unique 
		ON challenge_participants (challenge_id, user_id)
	`).Error; err != nil {
		return fmt.Errorf("failed to create challenge participants unique index: %w", err)
	}

	fmt.Println("✅ All unique constraints created successfully")
	return nil
}
