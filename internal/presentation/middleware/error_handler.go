package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"github.com/gin-gonic/gin"
	"challenge-app/internal/domain/exception"
	"errors"
)

type ErrorMiddleware struct {
}

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

                switch {
                case errors.As(err, &exception.InternalServerException{}):
                    var ie *exception.InternalServerException
                    _ = errors.As(err, &ie)
                    c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                        "status":     http.StatusInternalServerError,
                        "message":    nil,
                        "error_code": ie.Code(),
                    })
                case errors.As(err, &exception.ConflictException{}):
                    var ce *exception.ConflictException
                    _ = errors.As(err, &ce)
                    c.AbortWithStatusJSON(http.StatusConflict, gin.H{
                        "status":     http.StatusConflict,
                        "message":    ce.Error(),
                        "error_code": ce.Code(),
                    })
                case errors.As(err, &exception.ForbiddenException{}):
                    var fe *exception.ForbiddenException
                    _ = errors.As(err, &fe)
                    c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                        "status":     http.StatusForbidden,
                        "message":    fe.Error(),
                        "error_code": fe.Code(),
                    })
                case errors.As(err, &exception.NotFoundException{}):
                    var ne *exception.NotFoundException
                    _ = errors.As(err, &ne)
                    c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
                        "status":     http.StatusNotFound,
                        "message":    ne.Error(),
                        "error_code": ne.Code(),
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

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var clientErr exception.ClientError
			if errors.As(err, &clientErr) {
				status := clientErr.HTTPStatus()
				c.AbortWithStatusJSON(status, gin.H{
					"status":  status,
					"message": clientErr.Error(),
					"error_code": clientErr.Code(),
					"details": clientErr.Details(),
				})
				return 
			}
			
			log.Printf("UNHANDLED STANDARD ERROR: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"message": "An unexpected internal error occurred.",
				"error_code": "unhandled_internal_error",
			})
		}
	}
}
