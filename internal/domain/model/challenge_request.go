package model

import (
	"challenge-app/internal/domain/enum"
	"time"
)

type ChallengeRequest struct {
	ID          uint
	ChallengeID uint
	RequesterID uint
	Status      enum.RequestStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Challenge *ChallengeModel
	Requester *UserModel
}
