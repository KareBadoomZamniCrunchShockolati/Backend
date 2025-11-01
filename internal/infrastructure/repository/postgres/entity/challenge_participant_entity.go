package entity

import "gorm.io/gorm"

type ChallengeParticipantEntity struct {
	gorm.Model
	ChallengeID uint   `gorm:"not null;index"`
	UserID      uint   `gorm:"not null;index"`
	Status      string `gorm:"size:20;not null"` 
}

func (ChallengeParticipantEntity) TableName() string {
	return "challenge_participants"
}
