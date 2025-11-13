package handler

import (
	"github.com/gin-gonic/gin"
) 

type CategoryHandler interface {
	CreateCategory(c *gin.Context)
	GetCategory(c *gin.Context)
	GetAllCategories(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}