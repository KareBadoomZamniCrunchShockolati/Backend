package entity

import "gorm.io/gorm"

type ChallengeRequestEntity struct {
	gorm.Model
	ChallengeID uint `gorm:"not null;index"`
	UserID      uint `gorm:"not null;index"`
	Status      uint `gorm:"not null;default:1"`

	Challenge ChallengeEntity `gorm:"foreignKey:ChallengeID"`
	User      UserEntity      `gorm:"foreignKey:UserID"`
}

func (ChallengeRequestEntity) TableName() string {
	return "challenge_join_requests"
}
