package model

import (
	"challenge-app/internal/domain/enum"
	"time"
)

type ChallengeRequest struct {
	ID           uint
	ChallengeID  uint
	RequesterID  uint
	OwnerID      uint
	Status       enum.RequestStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}