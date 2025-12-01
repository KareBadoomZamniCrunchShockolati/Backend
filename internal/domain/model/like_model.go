package model

import (
	"time"
)

type Like struct {
	ID         uint
	EntityType LikeType
	EntityID   uint
	UserID     uint
	CreatedAt  time.Time
}
