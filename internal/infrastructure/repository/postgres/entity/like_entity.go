package entity

import (
	"time"

	"gorm.io/gorm"
)

type LikeEntity struct {
	ID         uint   `gorm:"primaryKey"`
	EntityType string `gorm:"type:varchar(20);not null"` // challenge, post, comment
	EntityID   uint   `gorm:"not null"`
	UserID     uint   `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (LikeEntity) TableName() string {
	return "likes"
}

func AutoMigrateLikeEntity(db *gorm.DB) error {
	return db.AutoMigrate(&LikeEntity{})
}

func CreateLikeUniqueIndex(db *gorm.DB) error {
	return db.Exec(`
        CREATE UNIQUE INDEX IF NOT EXISTS idx_likes_unique 
        ON likes (entity_type, entity_id, user_id)
    `).Error
}
