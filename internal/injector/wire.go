//go:build wireinject
// +build wireinject

package injector

import (
	"challenge-app/internal/application/service"
	"challenge-app/internal/infrastructure/repository/postgres"
	"challenge-app/internal/presentation/handler"
	"challenge-app/internal/presentation/middleware"
	"challenge-app/internal/presentation/router"
	"challenge-app/pkg/security"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
	"time"
)

// --- Value Adapters for External Inputs ---

// Adapts the JWT secret string from main.go into the []byte required by NewJWTService.
func provideJWTSecret(secret string) []byte {
	return []byte(secret)
}

// --- Provider Sets ---

// این مجموعه علاوه بر سرویس‌های امنیتی، آداپتور Secret Key را نیز برای JWT فراهم می‌کند.
var SecurityProvideSet = wire.NewSet(
	security.NewPasswordService,
	security.NewJWTService,
	provideJWTSecret, // آداپتور برای NewJWTService
)

var RepoProvideSet = wire.NewSet(
	postgres.NewUserRepository,
)

var ServiceProvideSet = wire.NewSet(
	service.NewUserService,
	service.NewAuthService,
)

var HandlerProvideSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewAuthHandler,
)

// The middleware is a provider itself, taking JWTService and returning gin.HandlerFunc.
var MiddlewareProvideSet = wire.NewSet(
	middleware.JWTAuthMiddleware,
)

// InitializeRouter wires up all dependencies and returns the Gin engine.
func InitializeRouter(db *gorm.DB, jwtSecret string, tokenExpiry time.Duration) (*gin.Engine, error) {
	wire.Build(
		
		// 2. Core Dependencies
		SecurityProvideSet,
		
		// 3. Repositories - استفاده از Provider Set تعریف شده در پکیج خودش
		RepoProvideSet, 
		
		// 4. Services & Presentation
		ServiceProvideSet,
		HandlerProvideSet,
		MiddlewareProvideSet,
		
		// 5. Final Router Setup (Router requires all handlers and middleware)
		router.SetupRouter,
	)
	return &gin.Engine{}, nil
}
