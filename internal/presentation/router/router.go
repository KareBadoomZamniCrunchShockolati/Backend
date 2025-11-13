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
	followHandler handler.FollowHandler,
	jwtMiddleware middleware.JWTMiddleware,
	challengeHandler handler.ChallengeHandler,
) *gin.Engine {
	r := gin.Default()
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
		v1.GET("/users/:user_id/followers", followHandler.GetFollowers)
		v1.GET("/users/:user_id/following", followHandler.GetFollowing)
		v1.GET("/users/:user_id/follow-stats", followHandler.GetFollowStats)

		v1.GET("/challenges", challengeHandler.GetAllChallenges)
		v1.GET("/challenges/:id", challengeHandler.GetChallengeByID)
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(jwtMiddleware.Handler())
	{
		protected.GET("/users", userHandler.GetAllUsers)
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.PUT("/users/profile", userHandler.UpdateProfile)
		protected.POST("/users/email/change", userHandler.InitiateEmailChange)
		protected.DELETE("/users/profile", userHandler.DeleteUser)

		// Protected follow routes
		protected.POST("/follow", followHandler.Follow)
		protected.DELETE("/follow", followHandler.Unfollow)
		protected.DELETE("/followers/remove", followHandler.RemoveFollower)
		protected.GET("/follow/status/:user_id", followHandler.CheckFollowStatus)

		protected.POST("/challenges", challengeHandler.CreateChallenge)
		protected.PUT("/challenges/:id", challengeHandler.UpdateChallenge)
		protected.DELETE("/challenges/:id", challengeHandler.DeleteChallenge)
	}

	return r
}
