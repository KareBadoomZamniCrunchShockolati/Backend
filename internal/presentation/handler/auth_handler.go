package handler

import (
	"net/http"
	"challenge-app/internal/application/dto"
	"challenge-app/internal/application/service"
	"github.com/gin-gonic/gin"
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

	user, token, err := h.AuthService.RegisterUser(req.Username, req.Email, req.Password, req.Bio)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user, "token": token})
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