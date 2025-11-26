package entity

import (
	"challenge-app/internal/domain/enum"
	"time"

	"gorm.io/gorm"
)

type ChallengeEntity struct {
	gorm.Model
	Title           string                   `gorm:"size:150;not null;index"`
	Description     string                   `gorm:"size:280;not null"`
	CategoryID      uint                     `gorm:"not null;index"`
	CreatorID       uint                     `gorm:"not null;index"`
	MaxParticipants uint                     `gorm:"default:0"`
	Visibility      enum.ChallengeVisibility `gorm:"type:string;not null;index:idx_creator_visibility"`
	ImageURL        string                   `gorm:"type:text"`
	Rule            string                   `gorm:"size:255;not null"`
	Timezone        string                   `gorm:"size:50"`
	StartTime       *time.Time               `gorm:"index"`
	EndTime         *time.Time
	IsStopped       bool `gorm:"default:false"`
	CommentsEnabled bool `gorm:"default:false"`
	LikeCount       uint `gorm:"default:0"`
	// Relations
	Participants []*ChallengeParticipantEntity `gorm:"foreignKey:ChallengeID"`
	Comments     []*ChallengeCommentEntity     `gorm:"foreignKey:ChallengeID"`
	Likes        []*ChallengeLikeEntity        `gorm:"foreignKey:ChallengeID"`
	Category     ChallengeCategoryEntity       `gorm:"foreignKey:CategoryID"`
}

func (ChallengeEntity) TableName() string {
	return "challenges"
}
