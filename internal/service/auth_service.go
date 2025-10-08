package service

import (
	"challenge-app/internal/domain"
	"challenge-app/pkg/security"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	UserRepo domain.UserRepository
}

func NewAuthService(repo domain.UserRepository) *AuthService {
	return &AuthService{UserRepo: repo}
}

// RegisterUser (CRUD - Create Logic)
func (s *AuthService) RegisterUser(username, email, password, bio string) (*domain.User, error) {
	// 1. Check if user already exists
	_, err := s.UserRepo.GetUserByEmail(email)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// 2. Hash the password (Security Rule)
	hash, err := security.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	// 3. Create the Domain Entity
	user := &domain.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
		Bio:          bio,
	}

	// 4. Persist the Domain Entity
	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("user creation failed: %w", err)
	}
	return user, nil
}

// GetUserByID (CRUD - Read Logic)
func (s *AuthService) GetUserByID(id uuid.UUID) (*domain.User, error) {
    return s.UserRepo.GetUserByID(id)
}


func (s *AuthService) GetAllUsers() ([]domain.User, error) {
    return s.UserRepo.GetAllUsers()
}

// UpdateUser (CRUD - Update Logic)
func (s *AuthService) UpdateUser(id uuid.UUID, username, bio, newEmail string) (*domain.User, error) {
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
func (s *AuthService) DeleteUser(id uuid.UUID) error {
    // Optionally add business logic here (e.g., check if user has active challenges)
    return s.UserRepo.DeleteUser(id)
}

func (s *AuthService) LoginUser(email, password string) (*domain.User, string, error) {
	// 1. Retrieve the user by email
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		// If user is not found or DB error, treat it as invalid credentials
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// 2. Check the password hash
	// Use the utility function to compare the plaintext password with the stored hash
	if !security.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// 3. Generate a JWT token
	token, err := security.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("could not generate token: %w", err)
	}
	// 4. Success: Return the user entity and token
	return user, token, nil
}