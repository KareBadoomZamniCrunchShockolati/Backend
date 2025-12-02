package model

import (
	"challenge-app/internal/domain/enum"
	"time"
)

type ChallengeModel struct {
	ID              uint
	Title           string
	Description     string
	CategoryID      uint
	CreatorID       uint
	MaxParticipants uint
	Visibility      enum.ChallengeVisibility
	ImageURL        string
	Rule            string
	Timezone        string
	Location        string
	Goal            int
	StartTime       time.Time
	EndTime         *time.Time
	IsStopped       bool
	CommentsEnabled bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}