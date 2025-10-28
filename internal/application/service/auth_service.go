package service

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/pkg/errs"
	"challenge-app/pkg/security"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AuthService struct {
	UserRepo    repository.UserRepository
	PasswordSvc security.PasswordService
	JwtService  security.JWTService
}

func NewAuthService(repo repository.UserRepository, passwordSvc security.PasswordService, jwtService security.JWTService) *AuthService {
	return &AuthService{
		UserRepo:    repo,
		PasswordSvc: passwordSvc,
		JwtService:  jwtService,
	}
}

// RegisterUser (CRUD - Create Logic)
func (s *AuthService) RegisterUser(username, email, password, bio string) (*model.UserModel, string, error) {
	// 1. Check if user already exists
	_, err := s.UserRepo.GetUserByEmail(email)
	if err == nil {
		return nil, "", &errs.ConflictError{
			MessageValue: fmt.Sprintf("User with email %s already exists.", email),
		}
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(fmt.Errorf("database failure while checking existing user: %w", err))
	}
	// 2. Hash the password (Security Rule)
	hash, err := s.PasswordSvc.HashPassword(password)
	if err != nil {
		panic(fmt.Errorf("UNRECOVERABLE ERROR: Password hashing failed, system state compromised: %w", err))
	}

	// 3. Create the Domain Entity
	user := &model.UserModel{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Bio:          bio,
	}

	// 4. Persist the Domain Entity
	err = s.UserRepo.CreateUser(user)
	if err != nil {
		panic(fmt.Errorf("failed to persist new user %s: %w", email, err))
	}

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		panic(fmt.Errorf("JWT token generation failed, system state compromised: %w", err))
	}
	return user, token, nil
}

func (s *AuthService) LoginUser(email, password string) (*model.UserModel, string, error) {
	// 1. Retrieve the user by email
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", &errs.UnAuthorizedError{
				MessageValue: "Invalid credentials.",
			}
		}

		panic(fmt.Errorf("database failure during login for email %s: %w", email, err))
	}
	// 2. Check the password hash
	// Use the PasswordService to compare the plaintext password with the stored hash
	match := s.PasswordSvc.CheckPasswordHash(password, user.PasswordHash)
	if !match {
		return nil, "", &errs.UnAuthorizedError{
			MessageValue: "Invalid credentials.",
		}
	}

	// 3. Generate a JWT token
	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		panic(fmt.Errorf("JWT token generation failed during login, system state compromised: %w", err))

	}
	// 4. Success: Return the user entity and token
	return user, token, nil
}
