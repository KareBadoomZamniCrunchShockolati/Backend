package model

import (
	"time"
)

// User is the core entity.
type UserModel struct {
	ID             uint
	Username       string
	Email          string
	PasswordHash   string
	Bio            string
	ProfilePicture string
	CreatedAt      time.Time
}
