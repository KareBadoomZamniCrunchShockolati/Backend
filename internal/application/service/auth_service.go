package service

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/pkg/email"
	"challenge-app/pkg/security"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"gorm.io/gorm"
)

type AuthService struct {
	UserRepo         repository.UserRepository
	VerificationRepo repository.VerificationRepository
	PasswordSvc      security.PasswordService
	JwtService       security.JWTService
	EmailService     email.EmailService
}

func NewAuthService(repo repository.UserRepository, verificationRepo repository.VerificationRepository,
	jwtService security.JWTService, emailService email.EmailService,
	passwordSvc security.PasswordService) *AuthService {
	return &AuthService{
		UserRepo:         repo,
		VerificationRepo: verificationRepo,
		EmailService:     emailService,
		PasswordSvc:      passwordSvc,
		JwtService:       jwtService,
	}
}

// RegisterUser (CRUD - Create Logic)
func (s *AuthService) RegisterUser(username, email, password, bio string) (*model.UserModel, string, error) {
	// 1. Check if user already exists
	_, err := s.UserRepo.GetUserByEmail(email)
	if err == nil {
		return nil, "", fmt.Errorf("user with email %s already exists", email)
	}

	// 2. Hash the password (Security Rule)
	hash, err := s.PasswordSvc.HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("could not hash password: %w", err)
	}

	// 3. Create the Domain Entity
	user := &model.UserModel{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Bio:          bio,
		Verified:     false,
	}

	// 4. Persist user
	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("user creation failed: %w", err)
	}

	// 5. Generate verification code
	code, err := generateVerificationCode1()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate verification code: %w", err)
	}

	ctx := context.Background()
	err = s.VerificationRepo.StoreVerificationCode(ctx, email, code, 5)
	if err != nil {
		return nil, "", fmt.Errorf("failed to store verification code: %w", err)
	}

	err = s.EmailService.SendVerificationEmail(email, code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send verification email: %w", err)
	}

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("could not generate token: %w", err)
	}

	return user, token, nil
}

func (s *AuthService) LoginUser(email string, password string) (*model.UserModel, string, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid credentials")
		}
		panic(fmt.Sprintf("database error while fetching user: %v", err))
	}

	if !user.Verified {
		return nil, "", fmt.Errorf("email not verified")
	}

	match := s.PasswordSvc.CheckPasswordHash(password, user.PasswordHash)
	if !match {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("could not generate token: %w", err)
	}

	return user, token, nil
}

func (s *AuthService) VerifyEmail(email, code string) (string, error) {
	ctx := context.Background()

	storedCode, err := s.VerificationRepo.GetVerificationCode(ctx, email)
	if err != nil {
		return "", fmt.Errorf("code expired or not found")
	}

	if storedCode != code {
		return "", fmt.Errorf("invalid code")
	}

	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	user.Verified = true
	_, err = s.UserRepo.UpdateUser(user)
	if err != nil {
		return "", fmt.Errorf("failed to update user: %w", err)
	}

	s.VerificationRepo.DeleteVerificationCode(ctx, email)

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) ResendVerificationEmail(email string) error {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.Verified {
		return fmt.Errorf("user is already verified")
	}

	code, err := generateVerificationCode1()
	if err != nil {
		return fmt.Errorf("failed to generate verification code: %w", err)
	}

	ctx := context.Background()
	err = s.VerificationRepo.StoreVerificationCode(ctx, email, code, 5)
	if err != nil {
		return fmt.Errorf("failed to store verification code: %w", err)
	}

	err = s.EmailService.SendVerificationEmail(email, code)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func generateVerificationCode1() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}
