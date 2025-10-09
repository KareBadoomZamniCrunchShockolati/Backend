package service

import (
	"github.com/google/uuid"
	"challenge-app/internal/domain"
)

// AuthServicer defines the interface for all authentication and user management 
// business logic functions.
type AuthServicer interface {
    // Auth Operations
    RegisterUser(username, email, password, bio string) (*domain.User, error)
    LoginUser(email, password string) (*domain.User, error)
}

type UserServicer interface {
    // CRUD Operations
    UpdateUser(id uuid.UUID, username, bio, newEmail string) (*domain.User, error)
    DeleteUser(id uuid.UUID) error
    
    // Read Operations
    GetUserByID(id uuid.UUID) (*domain.User, error)
    GetAllUsers() ([]domain.User, error)
}