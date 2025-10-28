package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"github.com/gin-gonic/gin"
	"challenge-app/pkg/errs" 
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
				
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "An unexpected server error occurred. Please try again later.",
					"error_code": "server_panic",
				})
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
			
			if clientErr, ok := err.(errs.ClientError); ok {
				c.AbortWithStatusJSON(clientErr.Status(), gin.H{
					"status": clientErr.Status(),
					"message": clientErr.Message(),
					"error_code": http.StatusText(clientErr.Status()), 
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