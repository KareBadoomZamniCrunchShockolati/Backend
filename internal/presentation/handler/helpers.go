package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Response(ctx *gin.Context, statusCode int, message string, data interface{}) {
	if message != "" {
		ctx.JSON(statusCode, gin.H{
			"status":  statusCode,
			"message": message,
			"data":    data,
		})
	} else {
		ctx.JSON(statusCode, gin.H{
			"status": statusCode,
			"data":   data,
		})
	}
}

func Validated[T any](ctx *gin.Context) T {
	var params T
	if err := ctx.ShouldBindUri(&params); err != nil {
		panic("Validation failed: " + err.Error())
	}

	if ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" || ctx.Request.Method == "PATCH" {
		if err := ctx.ShouldBind(&params); err != nil {
			panic("Validation failed: " + err.Error())
		}
	} else {
		if err := ctx.ShouldBindQuery(&params); err != nil {
			panic("Validation failed: " + err.Error())
		}
	}
	if err := validate.Struct(params); err != nil {
		panic("Validation failed: " + err.Error())
	}

	return params
}

func GetOffsetLimit(page, pageSize, defaultPage, defaultPageSize int) (int, int) {
	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	offset := (page - 1) * pageSize
	limit := pageSize

	return offset, limit
}
