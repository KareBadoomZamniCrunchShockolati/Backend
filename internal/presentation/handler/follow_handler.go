package handler

import (
	"challenge-app/internal/application/dto"
	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/exception"
	
	"net/http"
	"github.com/gin-gonic/gin"
)

type FollowHandlerImpl struct {
	FollowService serviceinterface.FollowServicer
}

func NewFollowHandler(followService serviceinterface.FollowServicer) *FollowHandlerImpl {
	return &FollowHandlerImpl{FollowService: followService}
}

func (h *FollowHandlerImpl) Follow(c *gin.Context) {
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	var req dto.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.FollowService.Follow(followerID.(uint), req.FollowingID)
	if err != nil {
		switch err.Error() {
		case "cannot follow yourself":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "user to follow not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "already following this user":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.FollowResponse{
		Message: "Successfully followed user",
	})
}

func (h *FollowHandlerImpl) Unfollow(c *gin.Context) {
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	var req dto.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.FollowService.Unfollow(followerID.(uint), req.FollowingID)
	if err != nil {
		if err.Error() == "follow relationship not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.FollowResponse{
		Message: "Successfully unfollowed user",
	})
}

func (h *FollowHandlerImpl) RemoveFollower(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	var req dto.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.FollowService.RemoveFollower(userID.(uint), req.FollowingID)
	if err != nil {
		if err.Error() == "user is not following you" || err.Error() == "follower user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove follower"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.FollowResponse{
		Message: "Successfully removed follower",
	})
}

func (h *FollowHandlerImpl) GetFollowers(c *gin.Context) {
	var uri  dto.UserURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(exception.NewMissingUserIDException())
		return
	}

	followers, err := h.FollowService.GetFollowers(uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get followers"})
		return
	}

	var userResponses []dto.UserResponse
	for _, user := range followers {
		userResponses = append(userResponses, dto.UserResponse{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			Bio:            user.Bio,
			ProfilePicture: user.ProfilePicture,
		})
	}

	c.JSON(http.StatusOK, dto.UserListResponse{
		Users: userResponses,
		Count: len(userResponses),
	})
}

func (h *FollowHandlerImpl) GetFollowing(c *gin.Context) {
	var uri  dto.UserURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(exception.NewMissingUserIDException())
		return
	}

	following, err := h.FollowService.GetFollowing(uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get following"})
		return
	}

	var userResponses []dto.UserResponse
	for _, user := range following {
		userResponses = append(userResponses, dto.UserResponse{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			Bio:            user.Bio,
			ProfilePicture: user.ProfilePicture,
		})
	}

	c.JSON(http.StatusOK, dto.UserListResponse{
		Users: userResponses,
		Count: len(userResponses),
	})
}

func (h *FollowHandlerImpl) CheckFollowStatus(c *gin.Context) {
	followerIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	followerID, ok := followerIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid context user ID"})
		return
	}

	var uri dto.UserURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in URI"})
		return
	}

	if followerID == uri.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	isFollowing, err := h.FollowService.IsFollowing(followerID, uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check follow status"})
		return
	}

	c.JSON(http.StatusOK, dto.FollowStatusResponse{
		IsFollowing: isFollowing,
	})
}

func (h *FollowHandlerImpl) GetFollowStats(c *gin.Context) {
	var uri dto.UserURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in URI"})
		return
	}

	stats, err := h.FollowService.GetFollowStats(uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get follow stats"})
		return
	}

	c.JSON(http.StatusOK, dto.FollowStatsResponse{
		FollowersCount: stats.FollowersCount,
		FollowingCount: stats.FollowingCount,
	})
}
