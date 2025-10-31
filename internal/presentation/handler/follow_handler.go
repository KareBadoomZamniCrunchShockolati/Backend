package handler

import (
	"challenge-app/internal/application/dto"
	serviceinterface "challenge-app/internal/application/service/interface"
	"net/http"
	"strconv"

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
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	followers, err := h.FollowService.GetFollowers(uint(userID))
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
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	following, err := h.FollowService.GetFollowing(uint(userID))
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
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	followingID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	isFollowing, err := h.FollowService.IsFollowing(followerID.(uint), uint(followingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check follow status"})
		return
	}

	c.JSON(http.StatusOK, dto.FollowStatusResponse{
		IsFollowing: isFollowing,
	})
}

func (h *FollowHandlerImpl) GetFollowStats(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.FollowService.GetFollowStats(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get follow stats"})
		return
	}

	c.JSON(http.StatusOK, dto.FollowStatsResponse{
		FollowersCount: stats.FollowersCount,
		FollowingCount: stats.FollowingCount,
	})
}
