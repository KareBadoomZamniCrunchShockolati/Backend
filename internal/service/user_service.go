package service

import (
	"challenge-app/internal/domain"
	"github.com/google/uuid"
	"fmt"
)

type UserService struct {
	UserRepo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}


// GetUserByID (CRUD - Read Logic)
func (s *UserService) GetUserByID(id uuid.UUID) (*domain.User, error) {
    return s.UserRepo.GetUserByID(id)
}


func (s *UserService) GetAllUsers() ([]domain.User, error) {
    return s.UserRepo.GetAllUsers()
}

// UpdateUser (CRUD - Update Logic)
func (s *UserService) UpdateUser(id uuid.UUID, username, bio, newEmail string) (*domain.User, error) {
	// 1. Retrieve the existing user
	user, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		// Return the error from the repository (e.g., gorm.ErrRecordNotFound)
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	// --- 2. Handle Email Change ---
	if newEmail != "" {
		// Check if the new email is already taken
		existingUser, _ := s.UserRepo.GetUserByEmail(newEmail)
		if existingUser != nil {
			return nil, fmt.Errorf("email address %s is already in use", newEmail)
		}

		// Apply the new email address
		user.Email = newEmail
		// NOTE: In a production app, you would set a `pending_email` and
		// send a verification link here instead of updating directly.
	}

	// --- 3. Handle Username/Bio Updates ---
	if username != "" {
		user.Username = username
	}
	if username == "" {
		return user, fmt.Errorf("username cannot be empty")
	}
	if bio != "" {
		user.Bio = bio
	}


	// 4. Persist changes to the repository
	if err := s.UserRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("failed to update user in repository: %w", err)
	}
    
	return user, nil
}

// DeleteUser (CRUD - Delete Logic)
func (s *UserService) DeleteUser(id uuid.UUID) error {
    return s.UserRepo.DeleteUser(id)
}