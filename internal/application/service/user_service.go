package service

import (
	"context"
	"strings"

	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/pkg/email"
	"errors"
	"fmt"
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
		if _, ok := err.(*exception.NotFoundException); ok {
			return nil, err
		}
		return nil, exception.NewRepositoryError(err)
	}
	if user == nil {
		return nil, exception.NewNotFoundException("User", fmt.Sprintf("%d", id), "USER_NOT_FOUND_001")
	}
	return user, nil
}

func (s *UserService) GetAllUsers() ([]model.UserModel, error) {
	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		if _, ok := err.(*exception.NotFoundException); ok {
			return nil, err
		}
		return nil, exception.NewRepositoryError(err)
	}
	if users == nil {
		return nil, exception.NewInternalServerException("UserRepo.GetAllUsers returned nil slice", "CODE_LOGIC_ERROR", nil)
	}
	return users, nil
}

// UpdateUser (CRUD - Update Logic)
func (s *UserService) UpdateUser(id uint, username, bio, newEmail string) (*model.UserModel, error) {
	user, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	if user == nil {
		return nil, exception.NewNotFoundException("User", fmt.Sprintf("%d", id), "USER_NOT_FOUND_002")
	}
	if username == "" {
		return nil, exception.NewBadRequestException("Username cannot be empty.", "INPUT_VALIDATION_001", nil)
	}

	// 4. Persist changes to the repository
	updatedUser, err := s.UserRepo.UpdateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, exception.NewConflictException("User field", "username/email", "USER_UPDATE_CONFLICT")
		}
		return nil, exception.NewRepositoryUpdateError(err)

	}
	if updatedUser == nil {
		return nil, exception.NewInternalServerException(fmt.Sprintf("UserRepo.UpdateUser returned nil for ID %d", id), "CODE_LOGIC_ERROR", nil)
	}

	return updatedUser, nil
}

// decided to implement changing emails in 2 parts. because of better data consistency and safer approach.
// starting the process of changing email
func (s *UserService) InitiateEmailChange(userID uint, newEmail string) error {
	// 1. getting the user with his ID
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if user == nil {
		return exception.NewNotFoundException("User", fmt.Sprintf("%d", userID), "USER_NOT_FOUND_003")
	}

	if user.Email == newEmail {
		return exception.NewBadRequestException("New email cannot be the same as the current email.", "EMAIL_SAME", nil)
	}

	// 3. Check if new email is already taken by another user
	existingUser, _ := s.UserRepo.GetUserByEmail(newEmail)
	if err != nil {
		return exception.NewRepositoryError(err)
	}

	if existingUser != nil && existingUser.ID != user.ID {
		return exception.NewConflictException("Email address", newEmail, "EMAIL_TAKEN")
	}

	// 4. Generate verification code
	code, err := generateVerificationCode1()
	if err != nil {
		return exception.NewVerificationCodeGenerationError(err)
	}

	// 5. Store email change request in Redis with User ID
	ctx := context.Background()

	emailChangeKey := fmt.Sprintf("emailchange:%s", newEmail)

	// Store the mapping: userID:code
	verificationData := fmt.Sprintf("%d:%s", user.ID, code)

	//fmt.Printf("DEBUG: Storing in Redis - Key: %s, Data: %s\n", emailChangeKey, verificationData)

	err = s.VerificationRepo.StoreVerificationCode(ctx, emailChangeKey, verificationData, 10) // 10 minutes expiration time
	if err != nil {
		return exception.NewRepositoryVerificationError(err)
	}
	//fmt.Printf("DEBUG: Successfully stored in Redis\n")

	// 6. Send verification email to the NEW email address
	//fmt.Printf("DEBUG: Sending verification email to: %s\n", newEmail)
	err = s.EmailService.SendEmailChangeVerification(newEmail, code, user.Email)
	if err != nil {
		// Clean up the stored verification if email fails
		//fmt.Printf("DEBUG: Email sending failed, cleaning up Redis: %v\n", err)
		_ = s.VerificationRepo.DeleteVerificationCode(ctx, emailChangeKey)
		return exception.NewEmailError(err)
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
	if err != nil && storedData == "" {
		// If storedData is empty, it means the key was not found or expired.
		// If err is not nil, it's an infra failure. We treat the expiry as a BadRequest
		return nil, exception.NewBadRequestException("Email change request not found or expired.", "EMAIL_CHANGE_EXPIRED", nil)
	}
	if err != nil {
		// Infrastructure failure on Redis lookup -> PANIC
		return nil, exception.NewRepositoryVerificationError(err)
	}
	fmt.Printf("DEBUG CompleteEmailChange: Found stored data: %s\n", storedData)

	// 2. Parse the stored data (format: "userID:code")
	var userID uint
	var storedCode string
	n, err := fmt.Sscanf(storedData, "%d:%s", &userID, &storedCode)
	if err != nil || n != 2 {
		//fmt.Printf("DEBUG CompleteEmailChange: Failed to parse stored data. Parsed %d items, error: %v\n", n, err)
		return nil, exception.NewInternalServerException(fmt.Sprintf("Corrupted verification data for key %s", emailChangeKey), "CODE_DATA_CORRUPT", err)
	}

	if strings.TrimSpace(storedCode) != code {
		return nil, exception.NewBadRequestException("Invalid verification code.", "VERIFY_CODE_INVALID", nil)
	}

	// 3. Get the user by ID ( because  it is more reliable than email)
	fmt.Printf("DEBUG CompleteEmailChange: Looking for user with ID: %d\n", userID)
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	fmt.Printf("DEBUG CompleteEmailChange: Found user - ID: %d, Current Email: %s\n", user.ID, user.Email)
	if user == nil {
		return nil, exception.NewNotFoundException("User", fmt.Sprintf("%d", userID), "USER_NOT_FOUND_004")
	}

	// Verify the old email matches ( for security purposes only)
	if user.Email != oldEmail {
		//fmt.Printf("DEBUG CompleteEmailChange: Security check failed - user email %s doesn't match provided old email %s\n", user.Email, oldEmail)
		return nil, exception.NewBadRequestException("Email change request mismatch with existing user.", "EMAIL_CHANGE_MISMATCH", nil)
	}

	// 4. Update user's email
	user.Email = newEmail
	user.Verified = true // since we just verified the new email

	//fmt.Printf("DEBUG CompleteEmailChange: Updating user email to: %s\n", newEmail)

	// 5. Save the updated user
	updatedUser, err := s.UserRepo.UpdateUser(user)
	if err != nil {
		//fmt.Printf("DEBUG CompleteEmailChange: Failed to update user in database: %v\n", err)
		return nil, exception.NewRepositoryUpdateError(err)
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
		var nf exception.NotFoundException
		if errors.As(err, &nf) {
			return nf
		}
		return exception.NewRepositoryError(err)

	}
	return nil
}
