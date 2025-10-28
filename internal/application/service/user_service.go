package service

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/pkg/errs"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}

// GetUserByID (CRUD - Read Logic)
func (s *UserService) GetUserByID(id uint) (*model.UserModel, error) {
	user, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errs.NotFoundError{
				Resource: fmt.Sprintf("User with ID %d", id),
			}
		}
		return nil, &errs.InternalServerError{
			Err: fmt.Errorf("database read failure for user ID %d: %w", id, err),
		}
	}
	if user == nil {
		panic(fmt.Sprintf("UserRepo.GetUserByID returned nil user for ID %d", id))
	}
	return user, nil
}

func (s *UserService) GetAllUsers() ([]model.UserModel, error) {
	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		return nil, &errs.InternalServerError{
			Err: fmt.Errorf("database read failure for all users: %w", err),
		}
	}
	if users == nil {
		panic("UserRepo.GetAllUsers returned nil slice unexpectedly")
	}
	return users, nil
}

// UpdateUser (CRUD - Update Logic)
func (s *UserService) UpdateUser(id uint, username, bio, newEmail string) (*model.UserModel, error) {
	user, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errs.NotFoundError{
				Resource: fmt.Sprintf("User with ID %d", id),
			}
		}
		return nil, &errs.InternalServerError{
			Err: fmt.Errorf("database lookup failure for update: %w", err),
		}
	}

	if user == nil {
		panic(fmt.Sprintf("Unexpected nil user returned during update for ID %d", id))
	}

	// --- 2. Handle Email Change ---
	if username == "" {
		return user, &errs.BadRequestError{
			MessageValue: "Username cannot be empty.",
		}
	}
	user.Username = username
	if bio != "" {
		user.Bio = bio
	}

	if newEmail == "" {
		return nil, &errs.BadRequestError{
			MessageValue: "Email address cannot be set to empty.",
		}
	}

	if newEmail != "" && newEmail != user.Email {
		existingUser, err := s.UserRepo.GetUserByEmail(newEmail)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errs.InternalServerError{
				Err: fmt.Errorf("database error checking email: %w", err),
			}
		}

		if existingUser != nil {
			return nil, &errs.ConflictError{
				MessageValue: fmt.Sprintf("Email address %s is already in use.", newEmail),
			}
		}
		user.Email = newEmail
	}

	// 4. Persist changes to the repository
	updatedUser, err := s.UserRepo.UpdateUser(user)
	if err != nil {
		// Check for possible duplicate key violation from repository (e.g., username change)
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, &errs.ConflictError{
				MessageValue: "The requested username or email is already in use.",
			}
		}
		// Recoverable DB write failure -> 500
		return nil, &errs.InternalServerError{
			Err: fmt.Errorf("failed to persist user update: %w", err),
		}
	}
	if updatedUser == nil {
		panic(fmt.Sprintf("UserRepo.UpdateUser returned nil for ID %d", id))
	}

	return updatedUser, nil
}

// DeleteUser (CRUD - Delete Logic)
func (s *UserService) DeleteUser(id uint) error {
	err := s.UserRepo.DeleteUser(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &errs.NotFoundError{
				Resource: fmt.Sprintf("User with ID %d", id),
			}
		}
		return &errs.InternalServerError{
			Err: fmt.Errorf("database delete failure for user ID %d: %w", id, err),
		}
	}
	return nil
}
