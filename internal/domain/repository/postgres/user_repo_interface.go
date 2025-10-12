package postgres

import (
	"challenge-app/internal/domain/model"
)

// UserRepository defines the methods for interacting with user data storage.
// This interface lives in the Domain layer to ensure the Application (Service) layer
// only depends on Domain concepts, not Infrastructure details.
type UserRepository interface {
	// CreateUser inserts a new user model into the storage and returns the created model (with ID).
	CreateUser(user *model.UserModel)  error

	// GetUserByEmail retrieves a user by their unique email address.
	GetUserByEmail(email string) (*model.UserModel, error)

	// GetUserByID retrieves a user by their unique internal ID.
	GetUserByID(id uint) (*model.UserModel, error)

	// UpdateUser updates an existing user's details. It expects a Domain model
	// containing the ID and the fields to be updated.
	UpdateUser(user *model.UserModel) (*model.UserModel, error)

	// DeleteUser removes a user from the storage by their ID.
	DeleteUser(id uint) error

	// GetAllUsers retrieves all user records.
	GetAllUsers() ([]model.UserModel, error)
}
