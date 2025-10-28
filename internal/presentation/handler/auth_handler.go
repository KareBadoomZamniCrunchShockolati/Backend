package handler

import (
	"challenge-app/internal/application/dto"
	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/pkg/errs"
	"challenge-app/pkg/validation"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
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
			c.Error(&errs.BadRequestError{
				MessageValue: fmt.Sprintf("Validation failed: %s", formatted),
			})
			return
		}

		// For any other JSON binding issues (syntax, type mismatch, etc.)
		c.Error(&errs.BadRequestError{
			MessageValue: fmt.Sprintf("Invalid request format: %s", err.Error()),
		})

	}

	bio := strings.TrimSpace(req.Bio)
	if bio == "" {
		bio = ""
	}

	user, token, err := h.AuthService.RegisterUser(req.Username, req.Email, req.Password, req.Bio)
	if err != nil {
		if clientErr, ok := err.(errs.ClientError); ok {
			c.Error(clientErr)
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
		c.Error(&errs.BadRequestError{
			MessageValue: fmt.Sprintf("Invalid request format: %s", err.Error()),
		})
		return
	}

	// 2. Service: Authenticate user
	user, token, err := h.AuthService.LoginUser(req.Email, req.Password)

	if err != nil {
		if clientErr, ok := err.(errs.ClientError); ok {
			c.Error(clientErr)
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
		c.Error(&errs.BadRequestError{
			MessageValue: fmt.Sprintf("Invalid request body: %s", err.Error()),
		})
		return
	}

	token, err := h.AuthService.VerifyEmail(req.Email, req.Code)
	if err != nil {
		if clientErr, ok := err.(errs.ClientError); ok {
			c.Error(clientErr)
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
		c.Error(&errs.BadRequestError{
			MessageValue: fmt.Sprintf("Invalid request body: %s", err.Error()),
		})
		return
	}

	err := h.AuthService.ResendVerificationEmail(req.Email)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.Error(&errs.NotFoundError{Resource: "User not found"})
			return
		case "user is already verified":
			c.Error(&errs.ConflictError{MessageValue: "User is already verified"})
			return
		default:
			panic(err) 
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent successfully",
		"email":   req.Email,
	})
}
