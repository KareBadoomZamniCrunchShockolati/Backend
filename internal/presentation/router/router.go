package router

import (
	handlers "challenge-app/internal/presentation/handler"
	"challenge-app/internal/presentation/middleware"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "challenge-app/docs"
)

// SetupRouter now accepts the UserHandler and the JWT Middleware function.
func SetupRouter(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	jwtMiddleware *middleware.JWTMiddleware,
) *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/api/v1")
	// Public routes (Login, Signup)
	{
		v1.POST("/auth/signup", authHandler.Signup)
		v1.POST("/auth/login", authHandler.Login)
	}

	// Protected Group: All routes here require a valid JWT token.
	protected := r.Group("/api/v1")
	protected.Use(jwtMiddleware.Handler())

	{
		protected.GET("/users", userHandler.GetAllUsers)
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.PUT("/users/profile", userHandler.UpdateProfile)
		protected.DELETE("/users/profile", userHandler.DeleteUser)
	}

	return r
}
