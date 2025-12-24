package handler

import (
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/localization"
	"challenge-app/internal/presentation/middleware"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			return fld.Name
		}
		name = strings.Split(name, ",")[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})
}

func Response(ctx *gin.Context, statusCode int, message string, data interface{}) {
	if message != "" {
		ctx.JSON(statusCode, gin.H{
			"status":  statusCode,
			"message": message,
			"data":    data,
		})
		return
	}
	ctx.JSON(statusCode, gin.H{
		"status": statusCode,
		"data":   data,
	})
}

func GetTranslatedSuccessMessage(c *gin.Context, key string) string {
	tr, exists := c.Get(middleware.TranslatorContextKey)
	if !exists {
		return key
	}

	if translator, ok := tr.(localization.TranslatorInstance); ok {
		if msg, err := translator.Translate("success." + key); err == nil && msg != "" {
			return msg
		}
	}
	return key
}

func hasTag(t reflect.Type, tagName string) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < t.NumField(); i++ {
		v := t.Field(i).Tag.Get(tagName)
		if v != "" && v != "-" {
			return true
		}
	}
	return false
}

func Validated[T any](ctx *gin.Context) T {
	var req T
	if hasTag(reflect.TypeOf(req), "uri") {
		if err := ctx.ShouldBindUri(&req); err != nil {
			panic(exception.NewBadRequestException(
				"INVALID_URI_PARAMS",
				map[string]any{"error": err.Error()},
			))
		}
	}
	ct := ctx.ContentType()
	ctRaw := strings.ToLower(ctx.GetHeader("Content-Type"))
	if ct == "application/json" || strings.Contains(ctRaw, "json") {
		if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
			panic(exception.NewInvalidRequestBodyException(err))
		}
	}
	if err := validate.Struct(req); err != nil {
		ves, ok := err.(validator.ValidationErrors)
		if !ok {
			panic(exception.NewValidationFailedException(map[string]any{
				"error": err.Error(),
			}))
		}

		fields := make([]map[string]any, 0, len(ves))
		for _, fe := range ves {
			fields = append(fields, map[string]any{
				"field": fe.Field(),
				"tag":   fe.Tag(),
				"param": fe.Param(),
			})
		}

		panic(exception.NewValidationFailedException(map[string]any{
			"fields": fields,
		}))
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
