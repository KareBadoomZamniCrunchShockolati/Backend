package handler

import (
	"challenge-app/internal/application/dto"
	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/exception"
	"challenge-app/pkg/validation"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
		
		validationMap := validation.FormatValidationError(err)

		details := make(map[string]any, len(validationMap))
		for k, v := range validationMap {
			details[k] = v
		}

		c.Error(exception.NewBadRequestException(
			"Input validation failed. Please review the details for specific field issues.",
			"INPUT_VALIDATION_FAILED",
			details,
		))
		return

	}

	bio := strings.TrimSpace(req.Bio)
	if bio == "" {
		bio = ""
	}

	user, token, err := h.AuthService.RegisterUser(req.Username, req.Email, req.Password, req.Bio)
	if err != nil {
		var clientErr exception.ClientError
		if errors.As(err, &clientErr) {
			c.Error(err)
			return
		}
		
		panic(err)
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
		c.Error(exception.NewBadRequestException(
			fmt.Sprintf("Invalid request format: %s", err.Error()),
			"INVALID_JSON_FORMAT",
			nil,
		))
		return
	}

	// 2. Service: Authenticate user
	user, token, err := h.AuthService.LoginUser(req.Email, req.Password)

	if err != nil {
		var clientErr exception.ClientError
		if errors.As(err, &clientErr) {
			c.Error(err)
			return
		}

		panic(err)
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
		c.Error(exception.NewBadRequestException(
			fmt.Sprintf("Invalid request body: %s", err.Error()),
			"INVALID_JSON_FORMAT",
			nil,
		))
		return
	}

	token, err := h.AuthService.VerifyEmail(req.Email, req.Code)
	if err != nil {
		var clientErr exception.ClientError		
		if errors.As(err, &clientErr) {
			c.Error(err)
			return
		}
		
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
		"token":   token,
	})
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(exception.NewBadRequestException(
			fmt.Sprintf("Invalid request body: %s", err.Error()),
			"INVALID_JSON_FORMAT",
			nil,
		))
		return
	}

	err := h.AuthService.ResendVerificationEmail(req.Email)
	if err != nil {
		var clientErr exception.ClientError
		if errors.As(err, &clientErr) {
			c.Error(err)
			return
		}
		
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent successfully",
		"email":   req.Email,
	})
}
