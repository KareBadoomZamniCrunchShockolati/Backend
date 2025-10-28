package service

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/pkg/email"
	"challenge-app/pkg/errs"
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
		Verified:     false,
	}

	// 4. Persist user
	err = s.UserRepo.CreateUser(user)
	if err != nil {
		panic(fmt.Errorf("failed to persist new user %s: %w", email, err))
	}

	// 5. Generate verification code
	code, err := generateVerificationCode1()
	if err != nil {
		return nil, "", &errs.InternalServerError{
			Err: errors.New("failed to generate verification code: %v"),
		}
	}

	ctx := context.Background()
	err = s.VerificationRepo.StoreVerificationCode(ctx, email, code, 5)
	if err != nil {
		return nil, "", &errs.InternalServerError{
			Err: errors.New("failed to store verification code: %v"),
		}
	}

	err = s.EmailService.SendVerificationEmail(email, code)
	if err != nil {
		return nil, "", &errs.InternalServerError{
			Err: errors.New("failed to send verification email: %v"),
		}
	}

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		panic(fmt.Errorf("JWT token generation failed, system state compromised: %w", err))
	}

	return user, token, nil
}

func (s *AuthService) LoginUser(email string, password string) (*model.UserModel, string, error) {
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

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		panic(fmt.Errorf("JWT token generation failed during login, system state compromised: %w", err))

	}

	return user, token, nil
}

func (s *AuthService) VerifyEmail(email, code string) (string, error) {
	ctx := context.Background()

	storedCode, err := s.VerificationRepo.GetVerificationCode(ctx, email)
	if err != nil {
		return "", &errs.BadRequestError{
			MessageValue: "Verification code expired or not found.",
		}
	}

	if storedCode != code {
		return "", &errs.BadRequestError{
			MessageValue: "Invalid verification code.",
		}
	}

	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return "", &errs.NotFoundError{
			Resource: fmt.Sprintf("User with email %s", email),
		}
	}

	user.Verified = true
	_, err = s.UserRepo.UpdateUser(user)
	if err != nil {
		return "", &errs.InternalServerError{
			Err: errors.New("Failed to update user verification: %v"),
		}
	}

	s.VerificationRepo.DeleteVerificationCode(ctx, email)

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		panic(fmt.Errorf("JWT generation failed after verification: %w", err))
	}

	return token, nil
}

func (s *AuthService) ResendVerificationEmail(email string) error {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return &errs.NotFoundError{
			Resource: fmt.Sprintf("User with email %s", email),
		}
	}

	if user.Verified {
		return &errs.ConflictError{
			MessageValue: "User is already verified.",
		}
	}

	code, err := generateVerificationCode1()
	if err != nil {
		return &errs.InternalServerError{
			Err: errors.New("failed to generate verification code: %v"),
		}
	}

	ctx := context.Background()
	err = s.VerificationRepo.StoreVerificationCode(ctx, email, code, 5)
	if err != nil {
		return &errs.InternalServerError{
			Err: errors.New("failed to store verification code: %v"),
		}
	}

	err = s.EmailService.SendVerificationEmail(email, code)
	if err != nil {
		return &errs.InternalServerError{
			Err: errors.New("failed to send verification email: %v"),
		}
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
