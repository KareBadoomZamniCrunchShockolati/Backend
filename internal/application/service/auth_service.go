// Updated internal/application/service/auth_service.go
package service

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/pkg/security"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gopkg.in/gomail.v2"
)

type AuthService struct {
	UserRepo   repository.UserRepository
	JwtService security.JWTService
	Redis      *redis.Client
}

func NewAuthService(repo repository.UserRepository, jwtService security.JWTService, rdb *redis.Client) *AuthService {
	return &AuthService{
		UserRepo:   repo,
		JwtService: jwtService,
		Redis:      rdb,
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
	hash, err := security.HashPassword(password)
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

	// 4. Persist the Domain Entity
	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("user creation failed: %w", err)
	}

	// Generate and send verification code
	code, err := generateVerificationCode()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate verification code: %w", err)
	}

	ctx := context.Background()
	err = s.Redis.Set(ctx, "verify:"+email, code, 5*time.Minute).Err()
	if err != nil {
		return nil, "", fmt.Errorf("failed to store verification code: %w", err)
	}

	err = sendVerificationEmail(email, code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send verification email: %w", err)
	}

	return user, "", nil //***** No token until verified ******
}

func (s *AuthService) LoginUser(email, password string) (*model.UserModel, string, error) {
	// 1. Retrieve the user by email
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		// If user is not found or DB error, treat it as invalid credentials
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Check if verified
	if !user.Verified {
		return nil, "", fmt.Errorf("email not verified")
	}

	// 2. Check the password hash
	if !security.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// 3. Generate a JWT token
	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("could not generate token: %w", err)
	}
	// 4. Success: Return the user entity and token
	return user, token, nil
}

func (s *AuthService) VerifyEmail(email, code string) (string, error) {
	ctx := context.Background()
	storedCode, err := s.Redis.Get(ctx, "verify:"+email).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("code expired or not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve code: %w", err)
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

	s.Redis.Del(ctx, "verify:"+email)

	token, err := s.JwtService.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return token, nil
}

func generateVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func sendVerificationEmail(to, code string) error {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s\nIt expires in 5 minutes.", code))

	return d.DialAndSend(m)
}

func (s *AuthService) ResendVerificationEmail(email string) error {
	// 1. Check if user exists and is not verified
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// 2. Check if user is already verified
	if user.Verified {
		return fmt.Errorf("user is already verified")
	}

	// 3. Generate new verification code
	code, err := generateVerificationCode()
	if err != nil {
		return fmt.Errorf("failed to generate verification code: %w", err)
	}

	// 4. Store new code in Redis (overwrites old one)
	ctx := context.Background()
	err = s.Redis.Set(ctx, "verify:"+email, code, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to store verification code: %w", err)
	}

	// 5. Send new verification email
	err = sendVerificationEmail(email, code)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}
