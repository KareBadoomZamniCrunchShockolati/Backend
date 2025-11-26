package service

import (
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/repository/postgres"
	"challenge-app/pkg/email"
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
func (s *AuthService) RegisterUser(username, email, password, bio string) (*model.UserModel, error) {
	existingUser, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		if !postgres.IsNotFoundError(err) {
			return nil, exception.NewRepositoryError("Failed to check existing user by email", err)
		}
	}

	if existingUser != nil {
		return nil, exception.NewConflictException("User", username, "USER_ALREADY_EXIST")
	}

	existingUser, err = s.UserRepo.GetUserByName(username)
	if err != nil {
		if !postgres.IsNotFoundError(err) {
			return nil, exception.NewRepositoryError("Failed to check existing user by name", err)
		}
	}

	if existingUser != nil {
		return nil, exception.NewConflictException("User", username, "USER_ALREADY_EXIST")
	}

	// 2. Hash the password
	hash, err := s.PasswordSvc.HashPassword(password)
	if err != nil {
		return nil, exception.NewHashedPasswordError(err)
	}

	// 3. Create the user domain entity
	user := &model.UserModel{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Bio:          bio,
		Verified:     true,
	}

	// 4. Persist the user
	if err := s.UserRepo.CreateUser(user); err != nil {
		return nil, exception.NewRepositoryError("Failed to persist new user", err)
	}

	// 5. Generate verification code
	code, err := generateVerificationCode1()
	if err != nil {
		return nil, exception.NewVerificationCodeGenerationError(err)
	}

	// 6. Store verification code
	ctx := context.Background()
	if err := s.VerificationRepo.StoreVerificationCode(ctx, email, code, 5); err != nil {
		return nil, exception.NewVerificationError(err)
	}

	// 7. Send verification email
	if err := s.EmailService.SendVerificationEmail(email, code); err != nil {
		return nil, exception.NewEmailError(err)
	}

	return user, nil
}

func (s *AuthService) LoginUser(email string, password string) (*model.UserModel, string, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		if postgres.IsNotFoundError(err) {
			return nil,"", exception.NewNotFoundException("User", email, "USER_EMAIL_NOT_FOUND")
		}
		return nil, "", err
	}
	// 2. Check the password hash
	if user == nil {
		return nil, "", exception.NewAuthInvalidCredentials()
	}

	match := s.PasswordSvc.CheckPasswordHash(password, user.PasswordHash)
	if !match {
		return nil, "", exception.NewAuthInvalidCredentials()
	}

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", exception.NewJWTError(err)
	}

	return user, token, nil
}

func (s *AuthService) VerifyEmail(email, code string) (string, error) {
	ctx := context.Background()

	storedCode, err := s.VerificationRepo.GetVerificationCode(ctx, email)
	if err != nil {
		return "", exception.NewVerificationCodeExpired(err)
	}

	if storedCode != code {
		return "", exception.NewInvalidVerificationCode()
	}

	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	user.Verified = true
	_, err = s.UserRepo.UpdateUser(user)
	if err != nil {
		return "", err
	}

	s.VerificationRepo.DeleteVerificationCode(ctx, email)

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		return "", exception.NewJWTError(err)
	}

	return token, nil
}

func (s *AuthService) ResendVerificationEmail(email string) error {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if user == nil {
		return exception.NewNotFoundException("User", email, "USER_NOT_FOUND_RESEND")
	}

	if user.Verified {
		return exception.NewConflictException("User", "verification status", "USER_ALREADY_VERIFIED")
	}

	code, err := generateVerificationCode1()
	if err != nil {
		return exception.NewVerificationCodeGenerationError(err)
	}

	ctx := context.Background()
	err = s.VerificationRepo.StoreVerificationCode(ctx, email, code, 5)
	if err != nil {
		return exception.NewVerificationError(err)
	}

	err = s.EmailService.SendVerificationEmail(email, code)
	if err != nil {
		return exception.NewEmailError(err)
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
