package service

import (
	"context"
	"strings"

	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/internal/domain/exception"
	"errors"
	"challenge-app/pkg/email"
	"fmt"

	"gorm.io/gorm"
)

type UserService struct {
	UserRepo         repository.UserRepository
	VerificationRepo repository.VerificationRepository
	EmailService     email.EmailService
}

func NewUserService(repo repository.UserRepository, vRepo repository.VerificationRepository, emailSvc email.EmailService) *UserService {
	return &UserService{
		UserRepo:         repo,
		VerificationRepo: vRepo,
		EmailService:     emailSvc,
	}
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

// decided to implement changing emails in 2 parts. because of better data consistency and safer approach.
// starting the process of changing email
func (s *UserService) InitiateEmailChange(userID uint, newEmail string) error {
	// 1. getting the user with his ID
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &errs.NotFoundError{Resource: fmt.Sprintf("User with ID %d", userID)}
		}
		return &errs.InternalServerError{
			Err: fmt.Errorf("failed to retrieve user for email change: %w", err),
		}
	}
	if user == nil {
		panic(fmt.Sprintf("UserRepo.GetUserByID returned nil user for ID %d", userID))
	}

	if user.Email == newEmail {
		return &errs.BadRequestError{MessageValue: "New email cannot be the same as the current email."}
	}
	

	// 3. Check if new email is already taken by another user
	existingUser, _ := s.UserRepo.GetUserByEmail(newEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return &errs.InternalServerError{
			Err: fmt.Errorf("failed to check existing email %s: %w", newEmail, err),
		}
	}

	if existingUser != nil && existingUser.ID != user.ID {
		return &errs.ConflictError{
			MessageValue: fmt.Sprintf("Email address %s is already in use.", newEmail),
		}
	}


	// 4. Generate verification code
	code, err := generateVerificationCode1()
	if err != nil {
		//fmt.Printf("DEBUG: Failed to generate verification code: %v\n", err)
		panic(fmt.Errorf("failed to generate verification code: %w", err))
	}

	// 5. Store email change request in Redis with User ID
	ctx := context.Background()

	emailChangeKey := fmt.Sprintf("emailchange:%s", newEmail)

	// Store the mapping: userID:code
	verificationData := fmt.Sprintf("%d:%s", user.ID, code)

	//fmt.Printf("DEBUG: Storing in Redis - Key: %s, Data: %s\n", emailChangeKey, verificationData)

	err = s.VerificationRepo.StoreVerificationCode(ctx, emailChangeKey, verificationData, 10) // 10 minutes expiration time
	if err != nil {
		//fmt.Printf("DEBUG: Failed to store in Redis: %v\n", err)
		return &errs.InternalServerError{
			Err: fmt.Errorf("failed to store email change verification: %w", err),
		}
	}
	//fmt.Printf("DEBUG: Successfully stored in Redis\n")

	// 6. Send verification email to the NEW email address
	//fmt.Printf("DEBUG: Sending verification email to: %s\n", newEmail)
	err = s.EmailService.SendEmailChangeVerification(newEmail, code, user.Email)
	if err != nil {
		// Clean up the stored verification if email fails
		//fmt.Printf("DEBUG: Email sending failed, cleaning up Redis: %v\n", err)
		_ = s.VerificationRepo.DeleteVerificationCode(ctx, emailChangeKey)
		return &errs.InternalServerError{
			Err: fmt.Errorf("failed to send verification email: %w", err),
		}
	}
	fmt.Printf("DEBUG: Verification email sent successfully\n")

	return nil
}

// completes the process of email changing and really changes that if conditions are true.
func (s *UserService) CompleteEmailChange(oldEmail, newEmail, code string) (*model.UserModel, error) {
	ctx := context.Background()

	//fmt.Printf("DEBUG CompleteEmailChange: Starting - oldEmail: %s, newEmail: %s, code: %s\n", oldEmail, newEmail, code)

	// 1. Retrieve the email change verification data
	// Use the SAME key format that includes the "verify:" prefix
	emailChangeKey := fmt.Sprintf("emailchange:%s", newEmail)
	//fmt.Printf("DEBUG CompleteEmailChange: Looking for Redis key: %s\n", emailChangeKey)

	storedData, err := s.VerificationRepo.GetVerificationCode(ctx, emailChangeKey)
	if err != nil {
		//fmt.Printf("DEBUG CompleteEmailChange: Redis key not found or error: %v\n", err)
		return nil, &errs.BadRequestError{MessageValue: "Email change request not found or expired."}
	}
	fmt.Printf("DEBUG CompleteEmailChange: Found stored data: %s\n", storedData)

	// 2. Parse the stored data (format: "userID:code")
	var userID uint
	var storedCode string
	n, err := fmt.Sscanf(storedData, "%d:%s", &userID, &storedCode)
	if err != nil || n != 2 {
		//fmt.Printf("DEBUG CompleteEmailChange: Failed to parse stored data. Parsed %d items, error: %v\n", n, err)
		panic(fmt.Errorf("corrupted verification data format for key %s: %v", emailChangeKey, err))
	}
	fmt.Printf("DEBUG CompleteEmailChange: Parsed - userID: %d, storedCode: '%s'\n", userID, storedCode)

	if strings.TrimSpace(storedCode) != code {
		return nil, &errs.BadRequestError{MessageValue: "Invalid verification code."}
	}


	// 3. Get the user by ID ( because  it is more reliable than email)
	fmt.Printf("DEBUG CompleteEmailChange: Looking for user with ID: %d\n", userID)
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errs.NotFoundError{Resource: fmt.Sprintf("User with ID %d", userID)}
		}
		return nil, &errs.InternalServerError{
			Err: fmt.Errorf("failed to retrieve user during email change: %w", err),
		}
	}
	fmt.Printf("DEBUG CompleteEmailChange: Found user - ID: %d, Current Email: %s\n", user.ID, user.Email)
	if user == nil {
		panic(fmt.Sprintf("UserRepo.GetUserByID returned nil for ID %d", userID))
	}

	// Verify the old email matches ( for security purposes only)
	if user.Email != oldEmail {
		//fmt.Printf("DEBUG CompleteEmailChange: Security check failed - user email %s doesn't match provided old email %s\n", user.Email, oldEmail)
		return nil, &errs.BadRequestError{
			MessageValue: "Email change request mismatch with existing user.",
		}
	}

	// 4. Update user's email
	user.Email = newEmail
	user.Verified = true // since we just verified the new email

	//fmt.Printf("DEBUG CompleteEmailChange: Updating user email to: %s\n", newEmail)

	// 5. Save the updated user
	updatedUser, err := s.UserRepo.UpdateUser(user)
	if err != nil {
		//fmt.Printf("DEBUG CompleteEmailChange: Failed to update user in database: %v\n", err)
		return nil, &errs.InternalServerError{
			Err: fmt.Errorf("failed to update user email: %w", err),
		}
	}
	fmt.Printf("DEBUG CompleteEmailChange: User updated successfully - New Email: %s\n", updatedUser.Email)

	// 6. Clean up the verification data for speed up and memory saving
	err = s.VerificationRepo.DeleteVerificationCode(ctx, emailChangeKey)
	if err != nil {
		fmt.Printf("DEBUG CompleteEmailChange: Warning - failed to delete verification code: %v\n", err)
	} else {
		fmt.Printf("DEBUG CompleteEmailChange: Verification data cleaned up\n")
	}

	// 7. Send confirmation email to the NEW email
	err = s.EmailService.SendEmailChangeConfirmation(newEmail, oldEmail)
	if err != nil {
		fmt.Printf("DEBUG CompleteEmailChange: Warning - failed to send confirmation email: %v\n", err)
	} else {
		fmt.Printf("DEBUG CompleteEmailChange: Confirmation email sent\n")
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