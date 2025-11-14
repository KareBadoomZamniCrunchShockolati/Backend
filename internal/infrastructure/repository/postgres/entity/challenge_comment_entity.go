package entity

import "gorm.io/gorm"

type ChallengeCommentEntity struct {
	gorm.Model
	ChallengeID uint   `gorm:"not null;index"`
	UserID      uint   `gorm:"not null;index"`
	Content     string `gorm:"type:text;not null"`

	Challenge ChallengeEntity `gorm:"foreignKey:ChallengeID"`
	User      UserEntity      `gorm:"foreignKey:UserID"`
}

func (ChallengeCommentEntity) TableName() string {
	return "challenge_comments"
}
