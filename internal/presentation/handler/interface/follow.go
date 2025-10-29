package handler

import "github.com/gin-gonic/gin"

type FollowHandler interface {
	//crud
	Follow(c *gin.Context)
	Unfollow(c *gin.Context)
	RemoveFollower(c *gin.Context)
	//read
	GetFollowers(c *gin.Context)
	GetFollowing(c *gin.Context)
	CheckFollowStatus(c *gin.Context)
	GetFollowStats(c *gin.Context)
}
