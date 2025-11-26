// internal/domain/model/challenge_like.go
package model

import (
	"time"
)

type ChallengeLike struct {
	ID          uint
	ChallengeID uint
	UserID      uint
	CreatedAt   time.Time
}