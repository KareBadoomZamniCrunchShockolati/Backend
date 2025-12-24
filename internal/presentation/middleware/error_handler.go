// internal/middleware/error_handler.go
package middleware

import (
	"challenge-app/internal/domain/exception"
	domainLoc "challenge-app/internal/domain/localization"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorMiddleware struct{}
func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

func (e *ErrorMiddleware) PanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf(
					"!!! PANIC RECOVERED !!! Request: %s %s Error: %v\nStack: %s",
					c.Request.Method, c.Request.URL.Path, r, debug.Stack(),
				)

				var tr domainLoc.TranslatorInstance
				if v, ok := c.Get(TranslatorContextKey); ok {
					if t, ok2 := v.(domainLoc.TranslatorInstance); ok2 {
						tr = t
					}
				}

				translate := func(key string, params ...string) string {
					if tr == nil {
						return ""
					}
					if s, err := tr.Translate(key, params...); err == nil && s != "" {
						return s
					}
					return ""
				}

				var err error
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("%v", r)
				}

				var ex exception.Error
				if errors.As(err, &ex) {
					status := ex.HTTPStatus()
					msg := ex.Code()
					if status < 500 {
						if m := translate("errors." + ex.Code()); m != "" {
							msg = m
						}
					}

					resp := gin.H{"status": status, "message": msg, "error_code": ex.Code()}

					if ex.Code() == exception.ErrorTypeInputValidationFailed {
						if d := translateValidationDetails(tr, ex.Details()); d != nil {
							resp["details"] = d
						}
					} else if status < 500 && ex.Details() != nil && len(ex.Details()) > 0 {
						resp["details"] = ex.Details()
					}

					c.AbortWithStatusJSON(status, resp)
					return
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "INTERNAL_SERVER_ERROR",
					"error_code": "INTERNAL_SERVER_ERROR",
				})
			}
		}()
		c.Next()
	}
}

func (e *ErrorMiddleware) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() 
		if len(c.Errors) > 0 && !c.Writer.Written() {
			lastErr := c.Errors.Last().Err

			var tr domainLoc.TranslatorInstance
			if v, ok := c.Get(TranslatorContextKey); ok {
				if t, ok2 := v.(domainLoc.TranslatorInstance); ok2 {
					tr = t
				}
			}

			translate := func(key string, params ...string) string {
				if tr == nil {
					return ""
				}
				if s, err := tr.Translate(key, params...); err == nil && s != "" {
					return s
				}
				return ""
			}

			var ex exception.Error
			if errors.As(lastErr, &ex) {
				status := ex.HTTPStatus()
				msg := ex.Code()
				if status < 500 {
					if m := translate("errors." + ex.Code()); m != "" {
						msg = m
					}
				}

				resp := gin.H{"status": status, "message": msg, "error_code": ex.Code()}

				if ex.Code() == exception.ErrorTypeInputValidationFailed {
					if d := translateValidationDetails(tr, ex.Details()); d != nil {
						resp["details"] = d
					}
				} else if status < 500 && ex.Details() != nil && len(ex.Details()) > 0 {
					resp["details"] = ex.Details()
				}

				c.AbortWithStatusJSON(status, resp)
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":     http.StatusInternalServerError,
				"message":    "INTERNAL_SERVER_ERROR",
				"error_code": "INTERNAL_SERVER_ERROR",
			})
		}
	}
}


func extractIdentifier(details map[string]any) string {
	if details == nil {
		return ""
	}
	if v, ok := details["identifier"]; ok {
		if s, ok2 := v.(string); ok2 {
			return s
		}
	}
	return ""
}

func normalizeFieldName(field string) string {
	return strings.ToLower(field)
}

func translateValidationDetails(tr domainLoc.TranslatorInstance, details map[string]any) map[string]any {
	if tr == nil || details == nil {
		return nil
	}

	raw, ok := details["fields"]
	if !ok || raw == nil {
		return nil
	}

	items, ok := raw.([]any)
	if !ok {
		if typed, ok2 := raw.([]map[string]any); ok2 {
			items = make([]any, len(typed))
			for i, m := range typed {
				items[i] = m
			}
		} else {
			return nil
		}
	}

	out := map[string]string{}

	for _, it := range items {
		m, ok := it.(map[string]any)
		if !ok || m == nil {
			continue
		}

		field, _ := m["field"].(string)
		tag, _ := m["tag"].(string)
		paramAny := m["param"]

		fieldKey := normalizeFieldName(field)
		fieldLabel, _ := tr.Translate(fieldKey)
		if fieldLabel == "" {
			fieldLabel = fieldKey
		}

		param := ""
		switch v := paramAny.(type) {
		case string:
			param = v
		case int, int64, float64:
			param = fmt.Sprintf("%v", v)
		default:
			if paramAny != nil {
				param = fmt.Sprintf("%v", paramAny)
			}
		}

		var msg string
		switch tag {
		case "required":
			msg = mustTranslate(tr, "errors.required", fieldLabel)
		case "min":
			msg = mustTranslate(tr, "errors.min", fieldLabel, param)
		case "max":
			msg = mustTranslate(tr, "errors.max", fieldLabel, param)
		case "oneof":
			msg = mustTranslate(tr, "errors.oneof", fieldLabel)
		case "email":
			msg = mustTranslate(tr, "errors.email")
		case "datetime", "time", "date":
			msg = mustTranslate(tr, "errors.datetime", fieldLabel)
		case "latitude", "longitude":
			msg = mustTranslate(tr, "errors."+tag, fieldLabel)
		default:
			msg = mustTranslate(tr, "errors.generic")
		}

		if msg == "" {
			msg = mustTranslate(tr, "errors.generic")
		}

		out[fieldKey] = msg
	}

	return map[string]any{"fields": out}
}

func mustTranslate(tr domainLoc.TranslatorInstance, key string, params ...string) string {
	if tr == nil {
		return ""
	}
	if s, err := tr.Translate(key, params...); err == nil && s != "" {
		return s
	}
	return ""
}