package entity

import (
	"gorm.io/gorm"
)

type PostEntity struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index"`
	Description string `gorm:"type:text;not null"`
	ChallengeID *uint  `gorm:"index"`
	Pictures    string `gorm:"type:text"`

	User      UserEntity       `gorm:"foreignKey:UserID"`
	Challenge *ChallengeEntity `gorm:"foreignKey:ChallengeID"`
}

func (PostEntity) TableName() string {
	return "posts"
}
