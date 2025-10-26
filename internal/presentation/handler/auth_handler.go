package handler

import (
	"challenge-app/internal/application/dto"
	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/pkg/validation"
	"net/http"
	"strings"
	"errors"
	validator "github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
)
type AuthHandler struct {
	AuthService serviceinterface.AuthServicer
}

func NewAuthHandler(authService serviceinterface.AuthServicer) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}



// Signup godoc
// @Summary Register a new user
// @Description Creates a new user account with username, email, password, and optional bio
// @Tags auth
// @Accept json
// @Produce json
// @Param signup body dto.SignupRequest true "Signup details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /auth/signup [post]
// Signup (CRUD - Create Handler)
func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		// Check if it's a validation error
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			formatted := validation.FormatValidationError(verrs)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": formatted,
			})
			return
		}

		// For any other JSON binding issues (syntax, type mismatch, etc.)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": gin.H{
				"message": err.Error(),
			},
		})
		return
	}

	
	bio := strings.TrimSpace(req.Bio)
	if bio == "" {
		bio = "" 
	}

	user, token, err := h.AuthService.RegisterUser(req.Username, req.Email, req.Password, req.Bio)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Check your email for verification code.",
		"user": dto.LoginResponse{
			ID: user.ID, Username: user.Username, Email: user.Email, Bio: user.Bio, Token: token,
		}, 
	})
}



// Login godoc
// @Summary Login existing user
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
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
		"user_response": dto.LoginResponse{
			ID: user.ID, Username: user.Username, Email: user.Email, Bio: user.Bio, Token: token,
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
		"token":   token,
	})
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
