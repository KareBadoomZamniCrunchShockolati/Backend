package router

import (
	"challenge-app/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router and defines all API endpoints.
func SetupRouter(authHandler *handlers.AuthHandler) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		// --- Public Auth Routes (No Auth Required) ---
		v1.POST("/auth/signup", authHandler.Signup)
		v1.POST("/auth/login", authHandler.Login)
		// --- User CRUD Routes (TEMPORARY: Passing ID in URL) ---
		
        // Read: Get a user profile by ID
		v1.GET("/users/:id", authHandler.GetProfile) 
        
		v1.GET("/users", authHandler.GetAllUsers) 
        // Update: Update a user profile by ID
		v1.PUT("/users/:id", authHandler.UpdateProfile) 
        
        // Delete: Delete a user by ID
		v1.DELETE("/users/:id", authHandler.DeleteUser) 
	}

	return r
}