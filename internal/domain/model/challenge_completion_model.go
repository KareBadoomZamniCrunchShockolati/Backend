package model

import "time"

type ChallengeCompletion struct {
	ID          uint
	UserID      uint
	ChallengeID uint
	IsCompleted bool
	CompletedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
