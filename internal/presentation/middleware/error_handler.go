package middleware

import (
	"challenge-app/internal/domain/exception"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

type ErrorMiddleware struct{}

func NewErrorProvider() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

func (p *ErrorMiddleware) PanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("!!! PANIC RECOVERED !!! Request: %s %s Error: %v\nStack: %s",
					c.Request.Method, c.Request.URL.Path, r, debug.Stack())

				err, ok := r.(error)
				if !ok {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"status":     http.StatusInternalServerError,
						"message":    "An unexpected server error occurred.",
						"error_code": "unknown_panic",
					})
					return
				}

				// Try to map to known exception types
				var ie *exception.InternalServerException
				var ce *exception.ConflictException
				var fe *exception.ForbiddenException
				var ne *exception.NotFoundException
				var ue *exception.UnauthorizedException
				var be *exception.BadRequestException

				switch {
				case errors.As(err, &be):
					status := be.HTTPStatus()
					resp := gin.H{
						"status":     status,
						"message":    be.Error(),
						"error_code": be.Code(),
					}
					if d := be.Details(); d != nil && len(d) > 0 {
						resp["details"] = d
					}
					c.AbortWithStatusJSON(status, resp)

				case errors.As(err, &ue):
					status := ue.HTTPStatus()
					c.AbortWithStatusJSON(status, gin.H{
						"status":     status,
						"message":    ue.Error(),
						"error_code": ue.Code(),
					})

				case errors.As(err, &ce):
					status := ce.HTTPStatus()
					resp := gin.H{
						"status":     status,
						"message":    ce.Error(),
						"error_code": ce.Code(),
					}
					if d := ce.Details(); d != nil && len(d) > 0 {
						resp["details"] = d
					}
					c.AbortWithStatusJSON(status, resp)

				case errors.As(err, &fe):
					status := fe.HTTPStatus()
					c.AbortWithStatusJSON(status, gin.H{
						"status":     status,
						"message":    fe.Error(),
						"error_code": fe.Code(),
					})

				case errors.As(err, &ne):
					status := ne.HTTPStatus()
					resp := gin.H{
						"status":     status,
						"message":    ne.Error(),
						"error_code": ne.Code(),
					}
					if d := ne.Details(); d != nil && len(d) > 0 {
						resp["details"] = d
					}
					c.AbortWithStatusJSON(status, resp)

				case errors.As(err, &ie):
					status := ie.HTTPStatus()
					c.AbortWithStatusJSON(status, gin.H{
						"status":     status,
						"message":    "Internal server error",
						"error_code": ie.Code(),
					})

				default:
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"status":     http.StatusInternalServerError,
						"message":    "An unexpected server error occurred.",
						"error_code": "unknown_internal_error",
					})
				}
			}
		}()
		c.Next()
	}
}

func (p *ErrorMiddleware) APIErrorTranslator() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var clientErr exception.ClientError
		if errors.As(err, &clientErr) {
			status := clientErr.HTTPStatus()

			resp := gin.H{
				"status":     status,
				"message":    clientErr.Error(),
				"error_code": clientErr.Code(),
			}
			if d := clientErr.Details(); d != nil && len(d) > 0 {
				resp["details"] = d
			}

			c.AbortWithStatusJSON(status, resp)
			return
		}

		log.Printf("UNHANDLED STANDARD ERROR: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":     http.StatusInternalServerError,
			"message":    "An unexpected internal error occurred.",
			"error_code": "unhandled_internal_error",
		})
	}
}
