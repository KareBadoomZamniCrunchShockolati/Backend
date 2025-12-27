package entity

import (
	"time"

	"gorm.io/gorm"
)

type ChallengeCompletionEntity struct {
	gorm.Model
	UserID      uint `gorm:"not null;index:idx_user_challenge"`
	ChallengeID uint `gorm:"not null;index:idx_user_challenge"`
	IsCompleted bool `gorm:"default:false"`
	CompletedAt time.Time

	User      UserEntity      `gorm:"foreignKey:UserID"`
	Challenge ChallengeEntity `gorm:"foreignKey:ChallengeID"`
}

func (ChallengeCompletionEntity) TableName() string {
	return "challenge_completions"
}
