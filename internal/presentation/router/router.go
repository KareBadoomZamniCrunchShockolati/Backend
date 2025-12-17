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

// SetupRouter sets up all routes, middleware, and swagger
func SetupRouter(
	userHandler handler.UserHandler,
	authHandler handler.AuthHandler,
	followHandler handler.FollowHandler,
	challengeHandler handler.ChallengeHandler,
	userDayHandler handler.UserDayHandler,
	postHandler handler.PostHandler,
	jwtMiddleware middleware.JWTMiddleware,
	errorMiddleware middleware.ErrorMiddleware,
) *gin.Engine {

	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(errorMiddleware.PanicRecovery())
	r.Use(errorMiddleware.APIErrorTranslator())

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "ACCEPT"},
		AllowCredentials: true,
	}))

	// Root health endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "API is running",
		})
	})

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// API v1 health endpoint
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "API v1 is running",
		})
	})

	// Public routes
	{
		v1.POST("/auth/signup", authHandler.Signup)
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/verify", authHandler.Verify)
		v1.POST("/resend-verification", authHandler.ResendVerification)
		v1.POST("/users/email/verify-change", userHandler.VerifyEmailChange)

		// Public challenge routes
		v1.GET("/challenges", challengeHandler.ListDiscoverableChallenges)
		v1.GET("/challenges/public", challengeHandler.ListPublicChallenges)
		v1.GET("/challenges/category/:category_id", challengeHandler.ListByCategory)
		v1.GET("/challenges/creator/:user_id", challengeHandler.ListByCreator)
		v1.GET("/challenges/:id/comments", challengeHandler.GetAllComments)
		v1.GET("/challenges/categories", challengeHandler.GetAllCategories)

		// Public post routes
		v1.GET("/posts/:id", postHandler.GetPost)
		v1.GET("/posts/:id/comments", postHandler.GetPostComments)

		// Public follow routes
		v1.GET("/users/:id/followers", followHandler.GetFollowers)
		v1.GET("/users/:id/following", followHandler.GetFollowing)
		v1.GET("/users/:id/follow-stats", followHandler.GetFollowStats)
	}

	// Protected routes
	protected := v1.Group("")
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

		// Protected challenge routes - CRUD
		protected.POST("/challenges", challengeHandler.CreateChallenge)
		protected.GET("/challenges/:id", challengeHandler.GetChallengeByID)
		protected.PUT("/challenges/:id", challengeHandler.UpdateChallenge)
		protected.DELETE("/challenges/:id", challengeHandler.DeleteChallenge)
		protected.PUT("/challenges/:id/stop", challengeHandler.StopChallenge)

		// Protected challenge participation routes
		protected.POST("/challenges/:id/join", challengeHandler.JoinPublicChallenge)
		protected.POST("/challenges/:id/request", challengeHandler.JoinPrivateChallenge)
		protected.POST("/challenges/:id/leave", challengeHandler.LeaveChallenge)

		// Protected challenge participant management (creator only)
		protected.DELETE("/challenges/:id/participants/:participant_id", challengeHandler.RemoveParticipant)
		protected.GET("/challenges/:id/participants", challengeHandler.ListChallengeParticipants)

		// Protected challenge invite routes
		protected.POST("/challenges/:id/invite", challengeHandler.InviteUserToChallenge)
		protected.POST("/challenges/invites/:invite_id/accept", challengeHandler.AcceptInvite)
		protected.POST("/challenges/invites/:invite_id/decline", challengeHandler.DeclineInvite)

		// Protected challenge request management (creator only)
		protected.POST("/challenges/requests/:request_id/accept", challengeHandler.AcceptJoinRequest)
		protected.POST("/challenges/requests/:request_id/decline", challengeHandler.DeclineJoinRequest)

		// Protected challenge comment routes
		protected.POST("/challenges/:id/comments", challengeHandler.AddComment)
		protected.GET("/challenges/comments/:comment_id", challengeHandler.GetComment)

		// Protected challenge like routes
		protected.POST("/challenges/:id/like", challengeHandler.LikeChallenge)
		protected.DELETE("/challenges/:id/like", challengeHandler.UnlikeChallenge)
		protected.GET("/challenges/:id/likes", challengeHandler.GetChallengeLikeCount)
		protected.GET("/challenges/:id/is-liked", challengeHandler.IsUserLikedChallenge)

		// Protected user-specific challenge routes
		protected.GET("/challenges/participating", challengeHandler.GetChallengesUserIsParticipating)
		protected.GET("/challenges/requests", challengeHandler.GetRequestsSentByUser)
		protected.GET("/challenges/invites", challengeHandler.GetInvitesSentToUser)
		protected.GET("/challenges/:id/requests", challengeHandler.GetRequestsSentToChallenge)
		protected.GET("/challenges/:id/invites", challengeHandler.GetInvitesSentFromChallenge)

		// Search & Discovery
		protected.GET("/challenges/creator-username/:username", challengeHandler.ListByCreatorUsername)
		protected.GET("/challenges/category-name/:category_name", challengeHandler.ListByCategoryName)
		protected.GET("/challenges/participant-count", challengeHandler.ListByParticipantCount)
		protected.GET("/challenges/like-count", challengeHandler.ListByLikeCount)
		protected.GET("/challenges/starting-soon", challengeHandler.ListChallengesStartingSoon)
		protected.GET("/challenges/top-creators", challengeHandler.ListTopCreatorsChallenge)
		protected.GET("/creators/top", challengeHandler.ListTopCreators)
		protected.GET("/challenges/search", challengeHandler.SearchChallenges)
		protected.GET("/challenges/my/search", challengeHandler.SearchChallengesUserIsParticipating)

		protected.GET("/challenges/:id/mutual-followers", challengeHandler.GetMutualFollowersInChallenge)

		protected.POST("/challenges/:id/days", userDayHandler.SaveDayData)
		protected.PUT("/challenges/:id/days", userDayHandler.UpdateDayData)
		protected.GET("/challenges/:id/days/:date", userDayHandler.GetDayData)
		protected.DELETE("/challenges/:id/days/:date", userDayHandler.DeleteDayData)
		protected.GET("/challenges/:id/progress", userDayHandler.GetGoalProgressChart)
		protected.GET("/challenges/:id/feelings", userDayHandler.GetFeelingCounts)

		// Post CRUD operations
		protected.POST("/posts", postHandler.CreatePost)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)

		// Post feed and user posts
		protected.GET("/posts/feed", postHandler.GetFeedPosts)
		protected.GET("/posts/user/:user_id", postHandler.GetUserPosts)
		protected.GET("/posts/challenge/:challenge_id", postHandler.GetPostsByChallenge)

		// Post comment routes
		protected.POST("/posts/:id/comments", postHandler.AddPostComment)

		// ========== POLYMORPHIC LIKE ROUTES ==========
		// These work for challenges, posts, AND comments
		protected.POST("/likes", postHandler.LikeEntity)
		protected.DELETE("/likes", postHandler.UnlikeEntity)

		protected.POST("/posts/images/presign", postHandler.PresignPostImages)
		protected.POST("/users/profile/picture", userHandler.UploadProfilePicture)
		protected.POST("/challenges/:id/cover", challengeHandler.UploadChallengeCover)
		// Post like routes
		protected.POST("/posts/:id/like", postHandler.LikePost)
		protected.DELETE("/posts/:id/like", postHandler.UnlikePost)
	}

	return r
}
