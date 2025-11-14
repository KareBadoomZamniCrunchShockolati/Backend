package entity

import "gorm.io/gorm"

type ChallengeParticipantEntity struct {
	gorm.Model
	ChallengeID uint `gorm:"not null;index"`
	UserID      uint `gorm:"not null;index"`
	Status      uint `gorm:"not null;default:1"`

	Challenge ChallengeEntity `gorm:"foreignKey:ChallengeID"`
	User      UserEntity      `gorm:"foreignKey:UserID"`
}

func (ChallengeParticipantEntity) TableName() string {
	return "challenge_participants"
}
