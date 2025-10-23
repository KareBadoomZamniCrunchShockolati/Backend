//go:build wireinject
// +build wireinject

package injector

import (
	"challenge-app/internal/application/service"
	service_interface "challenge-app/internal/application/service/interface"
	repository_interface "challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/repository/postgres/driver"
	"challenge-app/internal/infrastructure/repository/postgres"
	"challenge-app/internal/presentation/handler"
	handler_interface "challenge-app/internal/presentation/handler/interface"
	"challenge-app/internal/presentation/middleware"
	middleware_interface "challenge-app/internal/presentation/middleware/interface"
	"challenge-app/internal/presentation/router"
	"challenge-app/internal/bootstrap"
	"challenge-app/pkg/security"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type PostgresDSN string
type JWTSecret string
type TokenExpiry time.Duration

// --- Providers ---
func ProvideDSN() PostgresDSN {
	cfg := bootstrap.LoadEnv()
	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=" + cfg.SSLMode
	return PostgresDSN(dsn)
}

func ProvideJWTSecret() JWTSecret {
	cfg := bootstrap.LoadEnv()
	return JWTSecret(cfg.JWTSecretKey)
}

func ProvideTokenExpiry() TokenExpiry {
	return TokenExpiry(bootstrap.JWTTokenExpiry)
}

func ProvideJWTService(secret JWTSecret, expiry TokenExpiry) *security.JwtServiceImpl {
	return security.NewJWTService(string(secret), time.Duration(expiry))
}

func ProvidePostgresDB(dsn PostgresDSN) (*gorm.DB, error) {
	return driver.InitPostgresDB(string(dsn))
}

// --- Provider Sets ---
var SecurityProviderSet = wire.NewSet(
	security.NewPasswordService,
	ProvideJWTService,
	wire.Bind(new(security.PasswordService), new(*security.PasswordServiceImpl)),
	wire.Bind(new(security.JWTService), new(*security.JwtServiceImpl)),
)

var DatabaseProviderSet = wire.NewSet(
	ProvideDSN,
	ProvidePostgresDB,
)

var RepositoryProviderSet = wire.NewSet(
	postgres.NewUserRepository,
	wire.Bind(new(repository_interface.UserRepository), new(*postgres.UserRepository)),
)

var ServiceProviderSet = wire.NewSet(
	service.NewUserService,
	service.NewAuthService,
	wire.Bind(new(service_interface.UserServicer), new(*service.UserService)),
	wire.Bind(new(service_interface.AuthServicer), new(*service.AuthService)),
)

var HandlerProviderSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewAuthHandler,
	wire.Bind(new(handler_interface.UserHandler), new(*handler.UserHandler)),
	wire.Bind(new(handler_interface.AuthHandler), new(*handler.AuthHandler)),
)

var MiddlewareProviderSet = wire.NewSet(
	middleware.NewJWTMiddleware,
	wire.Bind(new(middleware_interface.JWTMiddleware), new(*middleware.JWTMiddleware)),
)

// --- Application ---
type Application struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func NewApplication(db *gorm.DB, router *gin.Engine) *Application {
	return &Application{
		DB:     db,
		Router: router,
	}
}

// --- Initialize Router ---
func InitializeRouter(db *gorm.DB) (*gin.Engine, error) {
	wire.Build(
		SecurityProviderSet,
		RepositoryProviderSet,
		ServiceProviderSet,
		HandlerProviderSet,
		MiddlewareProviderSet,
		router.SetupRouter,
		ProvideJWTSecret,
		ProvideTokenExpiry,
	)
	return &gin.Engine{}, nil
}

// --- Initialize Full Application ---
func InitializeApplication() (*Application, error) {
	wire.Build(
		DatabaseProviderSet,
		SecurityProviderSet,
		RepositoryProviderSet,
		ServiceProviderSet,
		HandlerProviderSet,
		MiddlewareProviderSet,
		router.SetupRouter,
		NewApplication,
		ProvideJWTSecret,
		ProvideTokenExpiry,
	)
	return &Application{}, nil
}
