package handler

import "github.com/gin-gonic/gin"

type NotificationHandler interface {
	List(c *gin.Context)
	MarkRead(c *gin.Context)
	MarkAllRead(c *gin.Context)
	UnreadCount(c *gin.Context)
}
