package model

import "time"

type ChallengeComment struct {
	ID          uint
	ChallengeID uint
	UserID      uint
	Username    string
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
