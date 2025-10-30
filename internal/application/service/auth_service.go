package service

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/pkg/email"
	"challenge-app/internal/domain/exception"
	"challenge-app/pkg/security"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
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
	existingUser, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		panic(exception.NewInternalServerException(
			"Unexpected repository error during user check", 
			"DB_CHECK_FAIL",
		).Wrap(err))
	}
	if existingUser != nil {
		// DOMAIN CONFLICT -> RETURN ConflictException
		return nil, "", exception.NewConflictException("User", "email", "USER_CONFLICT_001")
	}
	// 2. Hash the password (Security Rule)
	hash, err := s.PasswordSvc.HashPassword(password)
	if err != nil {
		panic(exception.NewInternalServerException(
			"Password hashing failed, system state compromised",
			"SYSTEM_HASH_FAIL",
		).Wrap(err))
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
		panic(exception.NewInternalServerException(
			fmt.Sprintf("Failed to persist new user %s", email),
			"DB_PERSIST_FAIL",
		).Wrap(err))
	}

	// 5. Generate verification code
	code, err := generateVerificationCode1()
	if err != nil {
		return nil, "", exception.NewInternalServerException(
			"Failed to generate verification code",
			"VERIFY_CODE_GEN_FAIL",
		).Wrap(err)
	}

	ctx := context.Background()
	err = s.VerificationRepo.StoreVerificationCode(ctx, email, code, 5)
	if err != nil {
		return nil, "", exception.NewInternalServerException(
			"Failed to store verification code",
			"VERIFY_STORE_FAIL",
		).Wrap(err)
	}

	err = s.EmailService.SendVerificationEmail(email, code)
	if err != nil {
		return nil, "", exception.NewInternalServerException(
			"Failed to send verification email",
			"EMAIL_SEND_FAIL",
		).Wrap(err)
	}

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		panic(exception.NewInternalServerException(
			"JWT token generation failed, system state compromised",
			"SYSTEM_JWT_FAIL",
		).Wrap(err))
	}

	return user, token, nil
}

func (s *AuthService) LoginUser(email string, password string) (*model.UserModel, string, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		panic(exception.NewInternalServerException(
			"Database failure during login", 
			"DB_LOGIN_FAIL",
			).Wrap(err))
	}
	// 2. Check the password hash
	// Use the PasswordService to compare the plaintext password with the stored hash
	match := s.PasswordSvc.CheckPasswordHash(password, user.PasswordHash)
	if !match {
		panic(exception.NewInternalServerException(
			"Database failure during login", "DB_LOGIN_FAIL",
			).Wrap(err))
	}

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		panic(exception.NewInternalServerException(
			"JWT token generation failed during login, system state compromised", 
			"SYSTEM_JWT_FAIL_LOGIN",
		).Wrap(err))

	}

	return user, token, nil
}

func (s *AuthService) VerifyEmail(email, code string) (string, error) {
	ctx := context.Background()

	storedCode, err := s.VerificationRepo.GetVerificationCode(ctx, email)
	if err != nil {
		return "", exception.NewBadRequestException("Verification code expired or not found.",
		 "VERIFY_CODE_EXPIRED", nil).Wrap(err)
	}

	if storedCode != code {
		return "", exception.NewBadRequestException("Invalid verification code.", "VERIFY_CODE_INVALID", nil)
	}

	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		panic(exception.NewInternalServerException("Database failure during verification", 
		"DB_VERIFY_FAIL",
		).Wrap(err))
	}

	user.Verified = true
	_, err = s.UserRepo.UpdateUser(user)
	if err != nil {
		panic(exception.NewInternalServerException("Failed to update user verification status", 
		"DB_UPDATE_VERIFY_FAIL",
		).Wrap(err))
	}

	s.VerificationRepo.DeleteVerificationCode(ctx, email)

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		panic(exception.NewInternalServerException("JWT generation failed after verification", 
		"SYSTEM_JWT_FAIL_VERIFY",
		).Wrap(err))
	}

	return token, nil
}

func (s *AuthService) ResendVerificationEmail(email string) error {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		panic(exception.NewInternalServerException("Database failure during resend check", 
		"DB_RESEND_FAIL",
		).Wrap(err))
	}

	if user == nil {
		return exception.NewNotFoundException("User", email, "USER_NOT_FOUND_RESEND")
	}

	if user.Verified {
		return exception.NewConflictException("User", "verification status", "USER_ALREADY_VERIFIED")
	}

	code, err := generateVerificationCode1()
	if err != nil {
		return exception.NewInternalServerException(
			"Failed to generate verification code",
			"VERIFY_CODE_GEN_FAIL",
		).Wrap(err)
	}

	ctx := context.Background()
	err = s.VerificationRepo.StoreVerificationCode(ctx, email, code, 5)
	if err != nil {
		return exception.NewInternalServerException(
			"Failed to store verification code",
			"VERIFY_STORE_FAIL",
		).Wrap(err)
	}

	err = s.EmailService.SendVerificationEmail(email, code)
	if err != nil {
		return exception.NewInternalServerException(
			"Failed to send verification email",
			"EMAIL_SEND_FAIL",
		).Wrap(err)
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
