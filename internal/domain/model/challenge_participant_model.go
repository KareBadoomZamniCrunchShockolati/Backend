package model

import (
	"challenge-app/internal/domain/enum"
	"time"
)

type ChallengeParticipant struct {
	ID          uint
	ChallengeID uint
	UserID      uint
	Username    string
	Status      enum.ParticipantStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
