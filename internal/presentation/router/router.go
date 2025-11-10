package router

import (
	handler "challenge-app/internal/presentation/handler/interface"
	middleware "challenge-app/internal/presentation/middleware/interface"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "challenge-app/docs"
	"github.com/gin-contrib/cors"
)

// SetupRouter sets up all routes, middleware, and swagger
func SetupRouter(
	userHandler handler.UserHandler,
	authHandler handler.AuthHandler,
	followHandler handler.FollowHandler,
	jwtMiddleware middleware.JWTMiddleware,
	errorMiddleware middleware.ErrorMiddleware,
) *gin.Engine {

	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(errorMiddleware.PanicRecovery()) // <- Add panic recovery
	r.Use(errorMiddleware.APIErrorTranslator()) // <- Translate client errors

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow requests from localhost:3000 or anywhere
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "ACCEPT"},
		AllowCredentials: true,
	}))

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// Public routes
	{
		v1.POST("/auth/signup", authHandler.Signup)
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/verify", authHandler.Verify)
		v1.POST("/resend-verification", authHandler.ResendVerification)
		v1.POST("/users/email/verify-change", userHandler.VerifyEmailChange)

		// Public follow routes
		v1.GET("/users/:id/followers", followHandler.GetFollowers)
		v1.GET("/users/:id/following", followHandler.GetFollowing)
		v1.GET("/users/:id/follow-stats", followHandler.GetFollowStats)
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(jwtMiddleware.Handler())
	{
		protected.GET("/users", userHandler.GetAllUsers)
		protected.GET("/users/:id", userHandler.GetUserByID)
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.PUT("/users/profile", userHandler.UpdateProfile)
		protected.POST("/users/email/change", userHandler.InitiateEmailChange)
		protected.DELETE("/users/profile", userHandler.DeleteUser)

		// Protected follow routes
		protected.POST("/follow", followHandler.Follow)
		protected.DELETE("/follow", followHandler.Unfollow)
		protected.DELETE("/followers/remove", followHandler.RemoveFollower)
		protected.GET("/follow/status/:id", followHandler.CheckFollowStatus)
	}

	return r
}
