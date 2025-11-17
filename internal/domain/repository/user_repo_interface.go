package repository

import (
	"challenge-app/internal/domain/model"
)

type UserRepository interface {
	CreateUser(user *model.UserModel)  error
	GetUserByEmail(email string) (*model.UserModel, error)
	GetUserByName(email string) (*model.UserModel, error)
	GetUserByID(id uint) (*model.UserModel, error)
	UpdateUser(user *model.UserModel) (*model.UserModel, error)
	DeleteUser(id uint) error
	GetAllUsers() ([]model.UserModel, error)
}
