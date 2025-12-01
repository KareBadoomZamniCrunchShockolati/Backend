package model

import (
	"challenge-app/internal/domain/enum"
	"time"
)

type ChallengeInvite struct {
	ID          uint
	ChallengeID uint
	InviterID   uint
	InviteeID   uint
	Status      enum.RequestStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}