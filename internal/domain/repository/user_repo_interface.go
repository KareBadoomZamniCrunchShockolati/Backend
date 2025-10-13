package repository

import (
	"challenge-app/internal/domain/model"
)

// UserRepository defines the methods for interacting with user data storage.
// This interface lives in the Domain layer to ensure the Application (Service) layer
// only depends on Domain concepts, not Infrastructure details.
type UserRepository interface {
	CreateUser(user *model.UserModel)  error
	GetUserByEmail(email string) (*model.UserModel, error)
	GetUserByID(id uint) (*model.UserModel, error)
	UpdateUser(user *model.UserModel) (*model.UserModel, error)
	DeleteUser(id uint) error
	GetAllUsers() ([]model.UserModel, error)
}
