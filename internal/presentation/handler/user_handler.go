package handler

import (
	"challenge-app/internal/application/dto"
	serviceinterface "challenge-app/internal/application/service/interface"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserService serviceinterface.UserServicer
}

func NewUserHandler(userService serviceinterface.UserServicer) *UserHandler {
	return &UserHandler{UserService: userService}
}

// GetProfile godoc
// @Summary Get current user's profile
// @Description Returns the profile of the authenticated user
// @Tags Users
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/profile [get]
// @Security BearerAuth
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

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieves a list of all registered users
// @Tags Users
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Failure 500 {object} map[string]string
// @Router /users [get]
// @Security BearerAuth
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

// UpdateProfile (CRUD - Update Handler)

// UpdateProfile godoc
// @Summary Update user profile
// @Description Updates username, bio, and/or email of the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Updated user info"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/profile [put]
// @Security BearerAuth
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed: User ID not found in context"})
		return
	}

	userID, ok := userIDVal.(uint)
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

	// 3. Call the business logic
	user, err := h.UserService.UpdateUser(
		userID,
		req.Username,
		req.Bio,
		req.NewEmail, 
	)

	// 4. Error Handling and Status Mapping
	if err != nil {
		// Map errors returned from the Service layer to HTTP status codes
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
			return
		}
		if strings.Contains(err.Error(), "already in use") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // 409 Conflict for resource collision
			return
		}
		// Generic server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}

	// 5. Success Response (Returning resource + message)
	var successMessage string
	if req.NewEmail != "" && req.NewEmail != user.Email {
		successMessage = "Profile and email updated successfully. You may need to log in again."
	} else {
		successMessage = "Profile updated successfully."
	}

	c.JSON(http.StatusOK, gin.H{
		"message": successMessage,
		"user": dto.UserResponse{ // Return the updated resource DTO
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

// DeleteUser godoc
// @Summary Delete current user
// @Description Deletes the account of the authenticated user
// @Tags Users
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/profile [delete]
// @Security BearerAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed: User ID not found in context"})
		return
	}

	userIDUUID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error: User ID format mismatch"})
		return
	}

	if err := h.UserService.DeleteUser(userIDUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
