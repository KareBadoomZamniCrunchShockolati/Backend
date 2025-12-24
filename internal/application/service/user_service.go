package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"challenge-app/internal/application/validator"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/storage"
	"challenge-app/pkg/email"
)

type UserService struct {
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
	emailService     email.EmailService
	objectStorage    storage.ObjectStorage
}

func NewUserService(repo repository.UserRepository, vRepo repository.VerificationRepository, emailSvc email.EmailService, storage storage.ObjectStorage) *UserService {
	return &UserService{
		userRepo:         repo,
		verificationRepo: vRepo,
		emailService:     emailSvc,
		objectStorage:    storage,
	}
}

// GetUserByID (CRUD - Read Logic)
func (s *UserService) GetUserByID(id uint) (*model.UserModel, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		var nf *exception.NotFoundException
		if errors.As(err, &nf) {
			return nil, nf
		}
		return nil, exception.NewRepositoryError(err)
	}
	if user == nil {
		return nil, exception.NewNotFoundException("User", fmt.Sprintf("%d", id), "USER_NOT_FOUND")
	}
	return user, nil
}

func (s *UserService) GetAllUsers() ([]model.UserModel, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		var nf *exception.NotFoundException
		if errors.As(err, &nf) {
			return nil, nf
		}
		return nil, exception.NewRepositoryError(err)
	}
	if users == nil {
		return nil, exception.NewInternalServerException("CODE_LOGIC_ERROR", map[string]any{
			"reason": "userRepo.GetAllUsers returned nil slice",
		}, errors.New("userRepo.GetAllUsers returned nil slice"))
	}
	return users, nil
}

// UpdateUser (CRUD - Update Logic)
func (s *UserService) UpdateUser(id uint, username, bio, newEmail string) (*model.UserModel, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	if user == nil {
		return nil, exception.NewNotFoundException("User", fmt.Sprintf("%d", id), "USER_NOT_FOUND")
	}
	if strings.TrimSpace(username) == "" {
		return nil, exception.NewBadRequestException("USERNAME_REQUIRED", map[string]any{
			"field": "username",
		})
	}

	// 4. Persist changes to the repository
	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint") {
			return nil, exception.NewConflictException(
				"USER_UPDATE_CONFLICT",
				"User",
				"username/email",
				"",
			)
		}
		return nil, exception.NewRepositoryUpdateError(err)
	}
	if updatedUser == nil {
		return nil, exception.NewInternalServerException("CODE_LOGIC_ERROR", map[string]any{
			"reason": fmt.Sprintf("userRepo.UpdateUser returned nil for ID %d", id),
		}, errors.New("userRepo.UpdateUser returned nil"))
	}

	_ = bio
	_ = newEmail

	return updatedUser, nil
}

// decided to implement changing emails in 2 parts. because of better data consistency and safer approach.
// starting the process of changing email
func (s *UserService) InitiateEmailChange(userID uint, newEmail string) error {
	// 1. getting the user with his ID
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if user == nil {
		return exception.NewNotFoundException("User", fmt.Sprintf("%d", userID), "USER_NOT_FOUND")
	}

	if strings.EqualFold(strings.TrimSpace(user.Email), strings.TrimSpace(newEmail)) {
		return exception.NewBadRequestException("EMAIL_SAME", map[string]any{
			"current": user.Email,
			"new":     newEmail,
		})
	}

	// 3. Check if new email is already taken by another user
	existingUser, err := s.userRepo.GetUserByEmail(newEmail)
	if err != nil {
		return exception.NewRepositoryError(err)
	}

	if existingUser != nil && existingUser.ID != user.ID {
		return exception.NewConflictException("EMAIL_TAKEN", "Email", "email", newEmail)
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

	err = s.verificationRepo.StoreVerificationCode(ctx, emailChangeKey, verificationData, 10) // 10 minutes expiration time
	if err != nil {
		return exception.NewRepositoryVerificationError(err)
	}
	//fmt.Printf("DEBUG: Successfully stored in Redis\n")

	// 6. Send verification email to the NEW email address
	//fmt.Printf("DEBUG: Sending verification email to: %s\n", newEmail)
	err = s.emailService.SendEmailChangeVerification(newEmail, code, user.Email)
	if err != nil {
		// Clean up the stored verification if email fails
		//fmt.Printf("DEBUG: Email sending failed, cleaning up Redis: %v\n", err)
		_ = s.verificationRepo.DeleteVerificationCode(ctx, emailChangeKey)
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

	storedData, err := s.verificationRepo.GetVerificationCode(ctx, emailChangeKey)
	if err != nil && storedData == "" {
		// If storedData is empty, it means the key was not found or expired.
		// If err is not nil, it's an infra failure. We treat the expiry as a BadRequest
		return nil, exception.NewBadRequestException("EMAIL_CHANGE_EXPIRED", map[string]any{
			"email": newEmail,
		})
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
		return nil, exception.NewInternalServerException("CODE_DATA_CORRUPT", map[string]any{
			"key": emailChangeKey,
		}, err)
	}

	if strings.TrimSpace(storedCode) != code {
		return nil, exception.NewBadRequestException("VERIFY_CODE_INVALID", map[string]any{
			"email": newEmail,
		})
	}

	// 3. Get the user by ID ( because  it is more reliable than email)
	fmt.Printf("DEBUG CompleteEmailChange: Looking for user with ID: %d\n", userID)
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if user == nil {
		return nil, exception.NewNotFoundException("User", fmt.Sprintf("%d", userID), "USER_NOT_FOUND")
	}
	fmt.Printf("DEBUG CompleteEmailChange: Found user - ID: %d, Current Email: %s\n", user.ID, user.Email)

	// Verify the old email matches ( for security purposes only)
	if user.Email != oldEmail {
		//fmt.Printf("DEBUG CompleteEmailChange: Security check failed - user email %s doesn't match provided old email %s\n", user.Email, oldEmail)
		return nil, exception.NewBadRequestException("EMAIL_CHANGE_MISMATCH", map[string]any{
			"expected_old": user.Email,
			"got_old":      oldEmail,
		})
	}

	// 4. Update user's email
	user.Email = newEmail
	user.Verified = true // since we just verified the new email

	//fmt.Printf("DEBUG CompleteEmailChange: Updating user email to: %s\n", newEmail)

	// 5. Save the updated user
	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		//fmt.Printf("DEBUG CompleteEmailChange: Failed to update user in database: %v\n", err)
		return nil, exception.NewRepositoryUpdateError(err)
	}
	if updatedUser == nil {
		return nil, exception.NewInternalServerException("CODE_LOGIC_ERROR", map[string]any{
			"reason": "userRepo.UpdateUser returned nil",
		}, errors.New("userRepo.UpdateUser returned nil"))
	}
	fmt.Printf("DEBUG CompleteEmailChange: User updated successfully - New Email: %s\n", updatedUser.Email)

	// 6. Clean up the verification data for speed up and memory saving
	err = s.verificationRepo.DeleteVerificationCode(ctx, emailChangeKey)
	if err != nil {
		fmt.Printf("DEBUG CompleteEmailChange: Warning - failed to delete verification code: %v\n", err)
	} else {
		fmt.Printf("DEBUG CompleteEmailChange: Verification data cleaned up\n")
	}

	// 7. Send confirmation email to the NEW email
	err = s.emailService.SendEmailChangeConfirmation(newEmail, oldEmail)
	if err != nil {
		fmt.Printf("DEBUG CompleteEmailChange: Warning - failed to send confirmation email: %v\n", err)
	} else {
		fmt.Printf("DEBUG CompleteEmailChange: Confirmation email sent\n")
	}

	return updatedUser, nil
}

// DeleteUser (CRUD - Delete Logic)
func (s *UserService) DeleteUser(id uint) error {
	err := s.userRepo.DeleteUser(id)
	if err != nil {
		var nf *exception.NotFoundException
		if errors.As(err, &nf) {
			return nf
		}
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *UserService) UploadProfilePicture(ctx context.Context, userID uint, file *multipart.FileHeader) (*model.UserModel, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if user == nil {
		return nil, exception.NewNotFoundException("User", fmt.Sprintf("%d", userID), "USER_NOT_FOUND")
	}
	data, mime, err := validator.ValidateProfileImage(file)
	if err != nil {
		return nil, exception.NewBadRequestException("INVALID_PROFILE_PICTURE", map[string]any{
			"reason": err.Error(),
		})
	}
	ext := extFromMIMEOrName(mime, file.Filename)
	if ext == "" {
		return nil, exception.NewBadRequestException("INVALID_PROFILE_PICTURE", map[string]any{
			"reason": "unsupported_extension",
			"mime":   mime,
			"name":   file.Filename,
		})
	}
	key := fmt.Sprintf("profiles/%d/%d%s", userID, time.Now().UTC().UnixNano(), ext)
	publicURL, err := s.objectStorage.Upload(ctx, key, mime, bytes.NewReader(data))
	if err != nil {
		return nil, exception.NewInternalServerException("S3_UPLOAD_FAILED", map[string]any{
			"reason": "failed to upload profile picture",
		}, err)
	}
	user.ProfilePicture = publicURL
	updated, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, exception.NewRepositoryUpdateError(err)
	}
	if updated == nil {
		return nil, exception.NewInternalServerException("CODE_LOGIC_ERROR", map[string]any{
			"reason": "userRepo.UpdateUser returned nil",
		}, errors.New("userRepo.UpdateUser returned nil"))
	}
	return updated, nil
}

func extFromMIMEOrName(mime, filename string) string {
	switch strings.ToLower(mime) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".jpeg" {
		return ".jpg"
	}
	if ext == ".jpg" || ext == ".png" || ext == ".webp" {
		return ext
	}
	return ""
}
