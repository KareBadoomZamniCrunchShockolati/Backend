//go:build wireinject
// +build wireinject

package injector

import (
	"challenge-app/internal/application/service"
	service_interface "challenge-app/internal/application/service/interface"
	repository_interface "challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/repository/postgres"
	"challenge-app/internal/presentation/handler"
	handler_interface "challenge-app/internal/presentation/handler/interface"
	"challenge-app/internal/presentation/middleware"
	middleware_interface "challenge-app/internal/presentation/middleware/interface"
	"challenge-app/internal/presentation/router"
	"challenge-app/pkg/security"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// --- Provider Sets ---

// Security & JWT
var SecurityProvideSet = wire.NewSet(
	security.NewPasswordService,
	security.NewJWTService,
	wire.Bind(new(security.PasswordService), new(*security.PasswordServiceImpl)),
	wire.Bind(new(security.JWTService), new(*security.JwtServiceImpl)),
)

// Repositories
var RepoProvideSet = wire.NewSet(
	postgres.NewUserRepository,
	wire.Bind(new(repository_interface.UserRepository), new(*postgres.UserRepository)),
)

// Services
var ServiceProvideSet = wire.NewSet(
	service.NewUserService,
	service.NewAuthService,
	wire.Bind(new(service_interface.AuthServicer), new(*service.AuthService)),
	wire.Bind(new(service_interface.UserServicer), new(*service.UserService)),
)

// Handlers
var HandlerProvideSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewAuthHandler,
	wire.Bind(new(handler_interface.UserHandler), new(*handler.UserHandler)),
	wire.Bind(new(handler_interface.AuthHandler), new(*handler.AuthHandler)),
)

// Middleware
var MiddlewareProvideSet = wire.NewSet(
	middleware.NewJWTMiddleware,
	wire.Bind(new(middleware_interface.JWTMiddleware), new(*middleware.JWTMiddleware)),
)

// Environment & Constants
// var ConfigProvideSet = wire.NewSet(
// 	bootstrap.LoadEnv,
// )

// InitializeRouter wires up all dependencies and returns the Gin engine.
func InitializeRouter(db *gorm.DB, jwtSecret string, tokenExpiry time.Duration) (*gin.Engine, error) {
	wire.Build(
		SecurityProvideSet,   // JWT / Password
		RepoProvideSet,       // Repos
		ServiceProvideSet,    // Services
		HandlerProvideSet,    // Handlers
		MiddlewareProvideSet, // Middleware
		router.SetupRouter,   // Final router
	)
	return &gin.Engine{}, nil
}
