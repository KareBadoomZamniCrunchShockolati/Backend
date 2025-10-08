package domain

import (
	"time"
	"github.com/google/uuid"
)

// User is the core entity.
type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	Bio          string
	CreatedAt    time.Time
}

// UserRepository defines the methods for data persistence.
type UserRepository interface {
    // CRUD Operations
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id uuid.UUID) (*User, error)
	GetAllUsers() ([]User, error)
	UpdateUser(user *User) error
	DeleteUser(id uuid.UUID) error
}