package handler

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/application/service"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	specialChars := "!@#$%^&*()-_=+[]{}|;:,.<>?/"

	for _, ch := range password {
		switch {
		case 'a' <= ch && ch <= 'z':
			hasLower = true
		case 'A' <= ch && ch <= 'Z':
			hasUpper = true
		case '0' <= ch && ch <= '9':
			hasDigit = true
		case strings.ContainsRune(specialChars, ch):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// Signup (CRUD - Create Handler)
func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	// Validate password strength
	if err := validatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _, err := h.AuthService.RegisterUser(req.Username, req.Email, req.Password, req.Bio)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registrated successfully. check your email for verification code.",
		"email":   user.Email,
		"user": dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Bio:      user.Bio,
		},
	})
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
		if err.Error() == "email not verified" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified. Please verify your email first."})
			return
		}
		// Catch any other server-side errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed due to server error"})
		return
	}

	// 3. Success: Respond with user details (JWT will be added here later)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Bio:      user.Bio,
		},
	})
}

func (h *AuthHandler) Verify(c *gin.Context) {
	var req dto.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	token, err := h.AuthService.VerifyEmail(req.Email, req.Code)
	if err != nil {
		if err.Error() == "invalid code" || err.Error() == "code expired or not found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification failed", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully", "token": token})
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	err := h.AuthService.ResendVerificationEmail(req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err.Error() == "user is already verified" {
			c.JSON(http.StatusConflict, gin.H{"error": "User is already verified"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resend verification email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent successfully",
		"email":   req.Email,
	})
}
