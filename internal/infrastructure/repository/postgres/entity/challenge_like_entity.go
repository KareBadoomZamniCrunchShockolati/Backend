package entity

import (
	"gorm.io/gorm"
)

type ChallengeLikeEntity struct {
	gorm.Model
	ChallengeID uint `gorm:"not null;index:idx_challenge_user,unique"`
	UserID      uint `gorm:"not null;index:idx_challenge_user,unique"`
	
	Challenge ChallengeEntity `gorm:"foreignKey:ChallengeID"`
	User      UserEntity      `gorm:"foreignKey:UserID"`
}