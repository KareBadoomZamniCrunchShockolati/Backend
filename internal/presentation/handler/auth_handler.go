package handler

import (
	"challenge-app/internal/application/dto"
	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/exception"
	"challenge-app/pkg/validation"
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
func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		formatted := validation.FormatValidationError(err)
		errorsMap := make(map[string]any, len(formatted))
		for k, v := range formatted {
			errorsMap[k] = v
		}
		c.Error(exception.NewValidationFailedException(map[string]any{
			"fields": errorsMap,
		}))
		return
	}

	bio := strings.TrimSpace(req.Bio)

	user, err := h.AuthService.RegisterUser(req.Username, req.Email, req.Password, bio)
	if err != nil {
		c.Error(err)
		return
	}
	message := GetTranslatedSuccessMessage(c, "USER_REGISTERED")

	Response(c, http.StatusCreated, message, dto.LoginResponse{
		ID: user.ID, Username: user.Username, Email: user.Email, Bio: user.Bio,
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(exception.NewInvalidRequestBodyException(err))
		return
	}

	user, token, err := h.AuthService.LoginUser(req.Email, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	Response(c, http.StatusOK, "Login successful", dto.LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
		Token:    token,
	})
}

func (h *AuthHandler) Verify(c *gin.Context) {
	var req dto.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(exception.NewInvalidRequestBodyException(err))
		return
	}

	token, err := h.AuthService.VerifyEmail(req.Email, req.Code)
	if err != nil {
		c.Error(err)
		return
	}

	Response(c, http.StatusOK, "Email verified successfully", gin.H{
		"token": token,
	})
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(exception.NewInvalidRequestBodyException(err))
		return
	}

	err := h.AuthService.ResendVerificationEmail(req.Email)
	if err != nil {
		c.Error(err)
		return
	}

	Response(c, http.StatusOK, "Verification email sent successfully", gin.H{
		"email": req.Email,
	})
}