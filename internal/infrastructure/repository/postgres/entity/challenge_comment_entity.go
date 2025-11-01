package entity

import "gorm.io/gorm"

type ChallengeCommentEntity struct {
	gorm.Model
	ChallengeID uint   `gorm:"not null;index"`
	UserID      uint   `gorm:"not null;index"`
	Content     string `gorm:"size:220;not null"` 
}

func (ChallengeCommentEntity) TableName() string {
	return "challenge_comments"
}
