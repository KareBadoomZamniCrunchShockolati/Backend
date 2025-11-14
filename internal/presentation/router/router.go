// internal/presentation/router/router.go
package router

import (
	handler "challenge-app/internal/presentation/handler/interface"
	middleware "challenge-app/internal/presentation/middleware/interface"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow requests from localhost:3000
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "ACCEPT"},
		AllowCredentials: true,
	}))
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

		// Public challenge routes
		v1.GET("/challenges", challengeHandler.GetAllChallenges)
		v1.GET("/challenges/:id", challengeHandler.GetChallengeByID)
		v1.GET("/challenges/category/:category_id", challengeHandler.ListByCategory)
		v1.GET("/challenges/creator/:user_id", challengeHandler.ListByCreator)
		v1.GET("/challenges/:id/participants", challengeHandler.ListChallengeParticipants)
		v1.GET("/challenges/:id/comments", challengeHandler.GetAllComments)
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

		// Protected challenge routes - CRUD
		protected.POST("/challenges", challengeHandler.CreateChallenge)
		protected.PUT("/challenges/:id", challengeHandler.UpdateChallenge)
		protected.DELETE("/challenges/:id", challengeHandler.DeleteChallenge)
		protected.PUT("/challenges/:id/stop", challengeHandler.StopChallenge)

		// Protected challenge participation routes
		protected.POST("/challenges/:id/join", challengeHandler.JoinPublicChallenge)
		protected.POST("/challenges/:id/request", challengeHandler.JoinPrivateChallenge)
		protected.POST("/challenges/:id/leave", challengeHandler.LeaveChallenge)

		// Protected challenge participant management (creator only)
		protected.DELETE("/challenges/:id/participants/:participant_id", challengeHandler.RemoveParticipant)

		// Protected challenge invite routes
		protected.POST("/challenges/:id/invite", challengeHandler.InviteUserToChallenge)
		protected.POST("/challenges/invites/:invite_id/accept", challengeHandler.AcceptInvite)
		protected.POST("/challenges/invites/:invite_id/decline", challengeHandler.DeclineInvite)

		// Protected challenge request management (creator only)
		protected.POST("/challenges/requests/:request_id/accept", challengeHandler.AcceptJoinRequest)
		protected.POST("/challenges/requests/:request_id/decline", challengeHandler.DeclineJoinRequest)

		// Protected challenge comment routes
		protected.POST("/challenges/:id/comments", challengeHandler.AddComment)

		// Protected user-specific challenge routes
		protected.GET("/challenges/participating", challengeHandler.GetChallengesUserIsParticipating)
		protected.GET("/challenges/requests", challengeHandler.GetRequestsSentByUser)
		protected.GET("/challenges/invites", challengeHandler.GetInvitesSentToUser)
		protected.GET("/challenges/:id/requests", challengeHandler.GetRequestsSentToChallenge)
		protected.GET("/challenges/:id/invites", challengeHandler.GetInvitesSentFromChallenge)
	}

	return r
}
