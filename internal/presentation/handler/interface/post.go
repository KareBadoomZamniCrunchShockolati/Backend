package handler

import "github.com/gin-gonic/gin"

type PostHandler interface {
	CreatePost(c *gin.Context)
	GetPost(c *gin.Context)
	UpdatePost(c *gin.Context)
	DeletePost(c *gin.Context)
	GetUserPosts(c *gin.Context)
	GetFeedPosts(c *gin.Context)
	LikePost(c *gin.Context)
	UnlikePost(c *gin.Context)
	PresignPostImages(c *gin.Context)
	GetPostsByChallenge(c *gin.Context)
	UploadPostImages(c *gin.Context)

	// Post-specific comment methods
	AddPostComment(c *gin.Context)
	GetPostComments(c *gin.Context)
}
