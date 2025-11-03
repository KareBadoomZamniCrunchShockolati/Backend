package entity

import (
	"challenge-app/internal/domain/enum"
	"time"

	"gorm.io/gorm"
)

type ChallengeEntity struct {
	gorm.Model
	Title           string                   `gorm:"size:150;not null"`
	Description     string                   `gorm:"size:280;not null"`
	Category        string                   `gorm:"size:100"`
	CreatorID       uint                     `gorm:"not null;index"`
	MaxParticipants uint                     `gorm:"default:0"`
	Visibility      enum.ChallengeVisibility `gorm:"size:20;not null"`
	ImageURL        string                   `gorm:"type:text"`
	Rule            string                   `gorm:"size:255"`
	Timezone        string                   `gorm:"size:50"`
	StartTime       *time.Time
	EndTime         *time.Time
	Stopped         bool

	// Relations
	Participants []*ChallengeParticipantEntity `gorm:"foreignKey:ChallengeID"`
	Comments     []*ChallengeCommentEntity     `gorm:"foreignKey:ChallengeID"`
}

func (ChallengeEntity) TableName() string {
	return "challenges"
}
