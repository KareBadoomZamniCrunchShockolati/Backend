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
	Rule            string
	Timezone        string
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Address         string  `json:"address"`
	Goal            int
	StartTime       time.Time
	EndTime         *time.Time
	IsStopped       bool
	CommentsEnabled bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CoverImage      string
}
