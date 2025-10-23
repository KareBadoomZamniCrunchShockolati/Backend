package handler

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/application/service"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// GetProfile (CRUD - Read Handler)
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed: User ID not found in context"})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Server processing failed",
			"details": "User ID format mismatch",
		})
		return
	}

	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID: user.ID, Username: user.Username, Email: user.Email, Bio: user.Bio,
	})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users", "details": err.Error()})
		return
	}
	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID: user.ID, Username: user.Username, Email: user.Email, Bio: user.Bio,
		})
	}

	c.JSON(http.StatusOK, userResponses)
}

// UpdateProfile (CRUD - Update Handler) - WITHOUT email change
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed: User ID not found in context"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error: User ID format mismatch"})
		return
	}

	// 2. Bind the request body to the DTO
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// 3. Call the business logic (without email change)
	user, err := h.UserService.UpdateUser(
		userIDUint,
		req.Username,
		req.Bio,
	)

	// 4. Error Handling and Status Mapping
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}

	// 5. Success Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully.",
		"user": dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Bio:      user.Bio,
		},
	})
}

// InitiateEmailChange starts the email change process
func (h *UserHandler) InitiateEmailChange(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed: User ID not found in context"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error: User ID format mismatch"})
		return
	}

	var req dto.InitiateEmailChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	err := h.UserService.InitiateEmailChange(userIDUint, req.NewEmail)
	if err != nil {
		if strings.Contains(err.Error(), "already in use") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "same as current") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate email change", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent to your new email address",
		"email":   req.NewEmail,
	})
}

// VerifyEmailChange completes the email change process
func (h *UserHandler) VerifyEmailChange(c *gin.Context) {
	var req dto.VerifyEmailChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	user, err := h.UserService.CompleteEmailChange(req.OldEmail, req.NewEmail, req.Code)
	if err != nil {
		if strings.Contains(err.Error(), "not found or expired") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email change request not found or expired"})
			return
		}
		if strings.Contains(err.Error(), "invalid verification code") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete email change", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email changed successfully",
		"user": dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Bio:      user.Bio,
		},
	})
}

// DeleteUser (CRUD - Delete Handler)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed: User ID not found in context"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error: User ID format mismatch"})
		return
	}

	if err := h.UserService.DeleteUser(userIDUint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
