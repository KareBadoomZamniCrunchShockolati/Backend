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
func (s *AuthService) RegisterUser(username, email, password, bio string) (*domain.User, string, error) {
	// 1. Check if user already exists
	_, err := s.UserRepo.GetUserByEmail(email)
	if err == nil {
		return nil, "", fmt.Errorf("user with email %s already exists", email)
	}

	// 2. Hash the password (Security Rule)
	hash, err := security.HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("could not hash password: %w", err)
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
	token, err := security.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("could not generate token: %w", err)
	}
	// 4. Persist the Domain Entity
	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("user creation failed: %w", err)
	}
	return user, token, nil
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