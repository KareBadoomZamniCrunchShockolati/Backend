package handler

import "github.com/gin-gonic/gin"

type PostHandler interface {
	CreatePost(c *gin.Context)
	GetPost(c *gin.Context)
	UpdatePost(c *gin.Context)
	DeletePost(c *gin.Context)
	GetUserPosts(c *gin.Context)
	GetFeedPosts(c *gin.Context)
	AddComment(c *gin.Context)
	GetComments(c *gin.Context)
	LikeEntity(c *gin.Context)
	UnlikeEntity(c *gin.Context)
	PresignPostImages(c *gin.Context)
}
