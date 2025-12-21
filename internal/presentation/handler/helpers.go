package handler

import (
	"challenge-app/internal/domain/exception"

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
	var req T

	if err := ctx.ShouldBindUri(&req); err != nil {
		panic(exception.NewBadRequestException(
			err.Error(),
			exception.ErrorTypeInvalidJSONFormat,
			nil,
		))
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		panic(exception.NewBadRequestException(
			err.Error(),
			exception.ErrorTypeInvalidJSONFormat,
			nil,
		))
	}
	contentType := ctx.GetHeader("Content-Type")
	if contentType == "application/json" && ctx.Request.ContentLength > 0 {
		if err := ctx.ShouldBindJSON(&req); err != nil {
			panic(exception.NewBadRequestException(
				err.Error(),
				exception.ErrorTypeInvalidJSONFormat,
				nil,
			))
		}
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := map[string]any{"details": err.Error()}
		panic(exception.NewValidationFailedException(validationErrors))
	}

	return req
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
