package handler

import (
	"challenge-app/internal/application/dto"
	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/exception"
	"errors"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
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
		c.Error(exception.NewMissingUserIDException())
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		panic(exception.NewContextCastError(errors.New("userID context value was not uint")))
	}

	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		var clientErr exception.ClientError
		if errors.As(err, &clientErr) {
			c.Error(err)
			return
		}		
		panic(err)
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID: user.ID, Username: user.Username, Email: user.Email, Bio: user.Bio, ProfilePicture: user.ProfilePicture,
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
		var clientErr exception.ClientError
		if errors.As(err, &clientErr) {
			c.Error(err)
			return
		}		
		panic(err)
	}
	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID: user.ID, Username: user.Username, Email: user.Email, Bio: user.Bio, ProfilePicture: user.ProfilePicture,
		})
	}

	c.JSON(http.StatusOK, userResponses)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	requesterID := c.GetUint("userID")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, e := h.UserService.GetUserByID(uint(id))
	if e != nil {
		switch e := e.(type) {
		case *exception.NotFoundException:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
		}
		return
	}
	if requesterID != user.ID {
		user.Email = "" // hide private info
	}
	c.JSON(http.StatusOK, dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
		ProfilePicture: user.ProfilePicture,
	})

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
		c.Error(exception.NewMissingUserIDException())
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		panic(exception.NewContextCastError(errors.New("userID context value was not uint")))
	}

	// 2. Bind the request body to the DTO
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(exception.NewInvalidRequestBodyException(err))
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
		var clientErr exception.ClientError
		if errors.As(err, &clientErr) {
			c.Error(err)
			return
		}
		panic(err)
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
		c.Error(exception.NewMissingUserIDException())
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		panic(exception.NewContextCastError(errors.New("userID context value was not uint")))
	}

	var req dto.InitiateEmailChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(exception.NewInvalidRequestBodyException(err))
		return
	}

	err := h.UserService.InitiateEmailChange(userIDUint, req.NewEmail)
	if err != nil {
		var clientErr exception.ClientError
		if errors.As(err, &clientErr) {
			c.Error(err)
			return
		}
		panic(err)
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
		c.Error(exception.NewInvalidRequestBodyException(err))
		return
	}

	user, err := h.UserService.CompleteEmailChange(req.OldEmail, req.NewEmail, req.Code)
	if err != nil {
		var clientErr exception.ClientError
		if errors.As(err, &clientErr) {
			c.Error(err)
			return
		}
		panic(err)
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
		c.Error(exception.NewMissingUserIDException())
		return
	}

	userIDUUID, ok := userID.(uint)
	if !ok {
		panic(exception.NewContextCastError(errors.New("userID context value was not uint")))
	}

	if err := h.UserService.DeleteUser(userIDUUID); err != nil {
		var clientErr exception.ClientError
		if errors.As(err, &clientErr) {
			c.Error(err)
			return
		}
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) UploadProfilePicture(c *gin.Context) {
	userID := c.GetUint("userID")

	file, err := c.FormFile("file")
	if err != nil {
		c.Error(exception.NewBadRequestException("file is required", "FILE_REQUIRED", nil))
		return
	}

	user, e := h.UserService.UploadProfilePicture(c.Request.Context(), userID, file)
	if e != nil {
		c.Error(e)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Profile picture uploaded successfully",
		"profile_picture": user.ProfilePicture,
	})
}