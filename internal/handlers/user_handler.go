package handlers

import (
	"net/http"
	"challenge-app/internal/dto"	
	"challenge-app/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
	"errors"
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

    userID, ok := userIDValue.(uuid.UUID)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error: User ID format mismatch"})
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

// UpdateProfile (CRUD - Update Handler)
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed: User ID not found in context"})
		return
	}

	userID, ok := userID.(uuid.UUID)
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
		userID.(uuid.UUID), 
		req.Username, 
		req.Bio, 
		req.NewEmail, // Field from the DTO
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
			ID: user.ID, 
			Username: user.Username, 
			Email: user.Email, 
			Bio: user.Bio,
		},
	})
}

// DeleteUser (CRUD - Delete Handler)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// TEMPORARY: Read ID from URL
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed: User ID not found in context"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
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
