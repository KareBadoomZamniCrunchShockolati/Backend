package model

import (
	"time"
)

type Post struct {
	ID          uint
	UserID      uint
	Description string
	ChallengeID *uint
	Pictures    []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
