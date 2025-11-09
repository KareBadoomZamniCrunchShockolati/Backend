//go:build wireinject
// +build wireinject

package injector

import (
	"challenge-app/internal/application/service"
	service_interface "challenge-app/internal/application/service/interface"
	"challenge-app/internal/bootstrap"
	repository_interface "challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/repository/postgres"
	"challenge-app/internal/infrastructure/repository/postgres/driver"
	redisRepo "challenge-app/internal/infrastructure/repository/redis"
	"challenge-app/internal/presentation/handler"
	handler_interface "challenge-app/internal/presentation/handler/interface"
	"challenge-app/internal/presentation/middleware"
	middleware_interface "challenge-app/internal/presentation/middleware/interface"
	"challenge-app/internal/presentation/router"
	"challenge-app/pkg/email"
	"challenge-app/pkg/security"
	"context"

	"github.com/go-playground/validator/v10"

	"challenge-app/pkg/validation"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type PostgresDSN string
type JWTSecret string
type TokenExpiry time.Duration
type RedisClient *redis.Client

// --- Providers ---
func ProvideDSN() PostgresDSN {
	cfg := bootstrap.LoadEnv()
	dsn := "host=" + cfg.Database.Host +
		" user=" + cfg.Database.User +
		" password=" + cfg.Database.Password +
		" dbname=" + cfg.Database.Name +
		" port=" + cfg.Database.Port +
		" sslmode=" + cfg.Database.SSLMode
	return PostgresDSN(dsn)
}

func ProvidePostgresDB(dsn PostgresDSN) (*gorm.DB, error) {
	return driver.InitPostgresDB(string(dsn))
}

func ProvideRedisClient() (*redis.Client, error) {
	cfg := bootstrap.LoadEnv()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return client, nil
}

func ProvideEmailService(cfg *bootstrap.Env) *email.EmailServiceImpl {
	return email.NewEmailService(cfg)
}

func ProvideJWTService(cfg *bootstrap.Env) *security.JwtServiceImpl {
	return security.NewJWTService(cfg)
}

func ProvideValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("password_policy", validation.PasswordValidationFunc)
	return v
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

var RedisProviderSet = wire.NewSet(
	ProvideRedisClient,
	redisRepo.NewVerificationRepository,
	wire.Bind(new(repository_interface.VerificationRepository), new(*redisRepo.VerificationRepository)),
)

var EmailProviderSet = wire.NewSet(
	bootstrap.LoadEnv,
	ProvideEmailService,
	wire.Bind(new(email.EmailService), new(*email.EmailServiceImpl)),
)

var RepositoryProviderSet = wire.NewSet(
	postgres.NewUserRepository,
	postgres.NewChallengeRepository,
	wire.Bind(new(repository_interface.UserRepository), new(*postgres.UserRepository)),
	wire.Bind(new(repository_interface.ChallengeRepository), new(*postgres.ChallengeRepository)),
)

var ServiceProviderSet = wire.NewSet(
	service.NewUserService,
	service.NewAuthService,
	service.NewChallengeService,
	wire.Bind(new(service_interface.UserServicer), new(*service.UserService)),
	wire.Bind(new(service_interface.AuthServicer), new(*service.AuthService)),
	wire.Bind(new(service_interface.ChallengeServicer), new(*service.ChallengeService)),
)

var HandlerProviderSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewAuthHandler,
	handler.NewChallengeHandler,
	wire.Bind(new(handler_interface.UserHandler), new(*handler.UserHandler)),
	wire.Bind(new(handler_interface.AuthHandler), new(*handler.AuthHandler)),
	wire.Bind(new(handler_interface.ChallengeHandler), new(*handler.ChallengeHandler)),
)

var MiddlewareProviderSet = wire.NewSet(
	middleware.NewJWTMiddleware,
	middleware.NewErrorProvider,
	wire.Bind(new(middleware_interface.ErrorMiddleware), new(*middleware.ErrorMiddleware)),
	wire.Bind(new(middleware_interface.JWTMiddleware), new(*middleware.JWTMiddleware)),
)
var FollowProviderSet = wire.NewSet(
	postgres.NewFollowRepository,
	service.NewFollowService,
	handler.NewFollowHandler,
	wire.Bind(new(repository_interface.FollowRepository), new(*postgres.FollowRepository)),
	wire.Bind(new(service_interface.FollowServicer), new(*service.FollowService)),
	wire.Bind(new(handler_interface.FollowHandler), new(*handler.FollowHandlerImpl)), // Updated this line
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
func InitializeRouter(db *gorm.DB, validator *validator.Validate) (*gin.Engine, error) {
	wire.Build(
		SecurityProviderSet,
		RepositoryProviderSet,
		RedisProviderSet,
		EmailProviderSet,
		ServiceProviderSet,
		HandlerProviderSet,
		FollowProviderSet, // ADD THIS LINE
		MiddlewareProviderSet,
		router.SetupRouter,
	)
	return &gin.Engine{}, nil
}

// --- Initialize Full Application ---
func InitializeApplication() (*Application, error) {
	wire.Build(
		DatabaseProviderSet,
		SecurityProviderSet,
		RepositoryProviderSet,
		RedisProviderSet,
		EmailProviderSet,
		ServiceProviderSet,
		HandlerProviderSet,
		FollowProviderSet, // ADD THIS LINE
		MiddlewareProviderSet,
		router.SetupRouter,
		NewApplication,
	)
	return &Application{}, nil
}
