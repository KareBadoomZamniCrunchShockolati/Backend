package service

import (
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"fmt"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}


// GetUserByID (CRUD - Read Logic)
func (s *UserService) GetUserByID(id uint) (*model.UserModel, error) {
    return s.UserRepo.GetUserByID(id)
}


func (s *UserService) GetAllUsers() ([]model.UserModel, error) {
    return s.UserRepo.GetAllUsers()
}

// UpdateUser (CRUD - Update Logic)
func (s *UserService) UpdateUser(id uint, username, bio, newEmail string) (*model.UserModel, error) {
	// 1. Retrieve the existing user
	user, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		// Return the error from the repository (e.g., gorm.ErrRecordNotFound)
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	// --- 2. Handle Email Change ---
	if newEmail != "" {
		// Check if the new email is already taken
		existingUser, _ := s.UserRepo.GetUserByEmail(newEmail)
		if existingUser != nil {
			return nil, fmt.Errorf("email address %s is already in use", newEmail)
		}

		// Apply the new email address
		user.Email = newEmail
		// NOTE: In a production app, you would set a `pending_email` and
		// send a verification link here instead of updating directly.
	}

	// --- 3. Handle Username/Bio Updates ---
	if username == "" {
    	return user, fmt.Errorf("username cannot be empty")
	}
	user.Username = username
	if bio != "" {
		user.Bio = bio
	}

	// 3. Persist changes to the repository
	updatedUser, err := s.UserRepo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user in repository: %w", err)
	}

	return updatedUser, nil
}

// decided to implement changing emails in 2 parts. because of better data consistency and safer approach.
// starting the process of changing email
func (s *UserService) InitiateEmailChange(userID uint, newEmail string) error {
	// 1. getting the user with his ID
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		//fmt.Printf("DEBUG: User not found with ID: %d, error: %v\n", userID, err)
		return fmt.Errorf("user not found: %w", err)
	}
	fmt.Printf("DEBUG: Found user - ID: %d, Current Email: %s, New Email: %s\n", user.ID, user.Email, newEmail)

	// 2. Check if new email is the same as current email
	if user.Email == newEmail {
		//fmt.Printf("DEBUG: New email is same as current email: %s\n", newEmail)
		return fmt.Errorf("new email cannot be the same as current email")
	}

	// 3. Check if new email is already taken by another user
	existingUser, _ := s.UserRepo.GetUserByEmail(newEmail)
	if existingUser != nil && existingUser.ID != user.ID {
		//fmt.Printf("DEBUG: Email already taken by user ID: %d\n", existingUser.ID)
		return fmt.Errorf("email address %s is already in use", newEmail)
	}
	fmt.Printf("DEBUG: New email is available: %s\n", newEmail)

	// 4. Generate verification code
	code, err := generateVerificationCode2()
	if err != nil {
		//fmt.Printf("DEBUG: Failed to generate verification code: %v\n", err)
		return fmt.Errorf("failed to generate verification code: %w", err)
	}
	fmt.Printf("DEBUG: Generated verification code: %s\n", code)

	// 5. Store email change request in Redis with User ID
	ctx := context.Background()

	emailChangeKey := fmt.Sprintf("emailchange:%s", newEmail)

	// Store the mapping: userID:code
	verificationData := fmt.Sprintf("%d:%s", user.ID, code)

	//fmt.Printf("DEBUG: Storing in Redis - Key: %s, Data: %s\n", emailChangeKey, verificationData)

	err = s.VerificationRepo.StoreVerificationCode(ctx, emailChangeKey, verificationData, 10) // 10 minutes expiration time
	if err != nil {
		//fmt.Printf("DEBUG: Failed to store in Redis: %v\n", err)
		return fmt.Errorf("failed to store email change verification: %w", err)
	}
	//fmt.Printf("DEBUG: Successfully stored in Redis\n")

	// 6. Send verification email to the NEW email address
	//fmt.Printf("DEBUG: Sending verification email to: %s\n", newEmail)
	err = s.EmailService.SendEmailChangeVerification(newEmail, code, user.Email)
	if err != nil {
		// Clean up the stored verification if email fails
		//fmt.Printf("DEBUG: Email sending failed, cleaning up Redis: %v\n", err)
		s.VerificationRepo.DeleteVerificationCode(ctx, emailChangeKey)
		return fmt.Errorf("failed to send verification email: %w", err)
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
		return nil, fmt.Errorf("email change request not found or expired")
	}
	fmt.Printf("DEBUG CompleteEmailChange: Found stored data: %s\n", storedData)

	// 2. Parse the stored data (format: "userID:code")
	var userID uint
	var storedCode string
	n, err := fmt.Sscanf(storedData, "%d:%s", &userID, &storedCode)
	if err != nil || n != 2 {
		//fmt.Printf("DEBUG CompleteEmailChange: Failed to parse stored data. Parsed %d items, error: %v\n", n, err)
		return nil, fmt.Errorf("invalid verification data format")
	}
	fmt.Printf("DEBUG CompleteEmailChange: Parsed - userID: %d, storedCode: '%s'\n", userID, storedCode)

	storedCode = strings.TrimSpace(storedCode)

	// invalid code handling
	if storedCode != code {
		//fmt.Printf("DEBUG CompleteEmailChange: Code mismatch - expected: '%s', got: '%s'\n", code, storedCode)
		return nil, fmt.Errorf("invalid verification code")
	}

	// 3. Get the user by ID ( because  it is more reliable than email)
	fmt.Printf("DEBUG CompleteEmailChange: Looking for user with ID: %d\n", userID)
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		fmt.Printf("DEBUG CompleteEmailChange: User not found with ID: %d, error: %v\n", userID, err)
		return nil, fmt.Errorf("user not found")
	}
	fmt.Printf("DEBUG CompleteEmailChange: Found user - ID: %d, Current Email: %s\n", user.ID, user.Email)

	// Verify the old email matches ( for security purposes only)
	if user.Email != oldEmail {
		//fmt.Printf("DEBUG CompleteEmailChange: Security check failed - user email %s doesn't match provided old email %s\n", user.Email, oldEmail)
		return nil, fmt.Errorf("email change request mismatch")
	}

	// 4. Update user's email
	user.Email = newEmail
	user.Verified = true // since we just verified the new email

	//fmt.Printf("DEBUG CompleteEmailChange: Updating user email to: %s\n", newEmail)

	// 5. Save the updated user
	updatedUser, err := s.UserRepo.UpdateUser(user)
	if err != nil {
		//fmt.Printf("DEBUG CompleteEmailChange: Failed to update user in database: %v\n", err)
		return nil, fmt.Errorf("failed to update user email: %w", err)
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
	return s.UserRepo.DeleteUser(id)
}

func generateVerificationCode2() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}
