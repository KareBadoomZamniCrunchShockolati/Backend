package serviceinterface

import (
	"challenge-app/internal/domain/model"
)

type UserServicer interface {
    // CRUD Operations
    UpdateUser(id uint, username, bio, newEmail string) (*model.UserModel, error)
    DeleteUser(id uint) error

    // Read Operations
    GetUserByID(id uint) (*model.UserModel, error)
    GetAllUsers() ([]model.UserModel, error)
}