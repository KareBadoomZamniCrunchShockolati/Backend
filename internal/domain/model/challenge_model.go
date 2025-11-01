package model

import (
	"challenge-app/internal/domain/enum"
	"time"
)

type Challenge struct {
	ID              uint
	Title           string
	Description     string
	Category        string
	CreatorID       uint
	MaxParticipants uint
	Visibility      enum.ChallengeVisibility
	ImageURL        string
	Rule            string
	Timezone        string
	StartTime       time.Time
	EndTime         *time.Time
	Stopped         bool
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Participants []ChallengeParticipant
	Comments     []ChallengeComment
}
