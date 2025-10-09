package router

import (
	"challenge-app/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRouter now accepts the UserHandler and the JWT Middleware function.
func SetupRouter(
    userHandler *handlers.UserHandler, 
    authHandler *handlers.AuthHandler,
    jwtMiddleware gin.HandlerFunc, 
) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	// Public routes (Login, Signup)
	{
		v1.POST("/auth/signup", authHandler.Signup)
		v1.POST("/auth/login", authHandler.Login)
	}

	// Protected Group: All routes here require a valid JWT token.
	protected := r.Group("/api/v1")
	protected.Use(jwtMiddleware) 

	{
		protected.GET("/users", userHandler.GetAllUsers) 
		protected.GET("/users/profile", userHandler.GetProfile) 
		protected.PUT("/users/profile", userHandler.UpdateProfile) 
		protected.DELETE("/users/profile", userHandler.DeleteUser)
	}

	return r
}