package router

import (
	handler "challenge-app/internal/presentation/handler/interface"
	middleware "challenge-app/internal/presentation/middleware/interface"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "challenge-app/docs"
)

func SetupRouter(
	userHandler handler.UserHandler,
	authHandler handler.AuthHandler,
	jwtMiddleware middleware.JWTMiddleware,
) *gin.Engine {
	r := gin.Default()
  r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/api/v1")
	// Public routes (Login, Signup, Verification)
	{
		v1.POST("/auth/signup", authHandler.Signup)
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/verify", authHandler.Verify)
		v1.POST("/resend-verification", authHandler.ResendVerification)
		v1.POST("/users/email/verify-change", userHandler.VerifyEmailChange)
	}

	// Protected Group: All routes here require a valid JWT token.
	protected := r.Group("/api/v1")
	protected.Use(jwtMiddleware.Handler())

	{
		protected.GET("/users", userHandler.GetAllUsers)
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.PUT("/users/profile", userHandler.UpdateProfile)
		protected.POST("/users/email/change", userHandler.InitiateEmailChange)
		protected.DELETE("/users/profile", userHandler.DeleteUser)
	}

	return r
}
