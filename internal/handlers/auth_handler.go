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

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// Signup (CRUD - Create Handler)
func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	user, err := h.AuthService.RegisterUser(req.Username, req.Email, req.Password, req.Bio)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.UserResponse{
		ID: user.ID, Username: user.Username, Email: user.Email, Bio: user.Bio,
	})
}

// GetProfile (CRUD - Read Handler)
func (h *AuthHandler) GetProfile(c *gin.Context) {
    // TEMPORARY: Read ID from URL until JWT is implemented
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := h.AuthService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID: user.ID, Username: user.Username, Email: user.Email, Bio: user.Bio,
	})
}


func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	users, err := h.AuthService.GetAllUsers()
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
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// 1. TEMPORARY: Read ID from URL (Will be replaced by JWT later)
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// 2. Bind the request body to the DTO
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// 3. Call the business logic
	user, err := h.AuthService.UpdateUser(
        userID, 
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
func (h *AuthHandler) DeleteUser(c *gin.Context) {
    // TEMPORARY: Read ID from URL
    userIDStr := c.Param("id")
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    if err := h.AuthService.DeleteUser(userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// 1. Bind the JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// 2. Service: Authenticate user
	user, token, err := h.AuthService.LoginUser(req.Email, req.Password)

	if err != nil {
		// Treat all authentication failures as 401 Unauthorized for security
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		// Catch any other server-side errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed due to server error"})
		return
	}

	// 3. Success: Respond with user details (JWT will be added here later)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token": token,
		"user": dto.UserResponse{
			ID: user.ID, Username: user.Username, Email: user.Email, Bio: user.Bio,
		},
	})
}