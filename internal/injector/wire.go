//go:build wireinject
// +build wireinject

package injector

import (
	"context"
	"time"

	"challenge-app/internal/application/service"
	service_interface "challenge-app/internal/application/service/interface"
	workers "challenge-app/internal/application/service/workers"
	"challenge-app/internal/bootstrap"
	domainLoc "challenge-app/internal/domain/localization"
	repository_interface "challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/localization"
	"challenge-app/internal/infrastructure/realtime/ws"
	"challenge-app/internal/infrastructure/repository/postgres"
	"challenge-app/internal/infrastructure/repository/postgres/driver"
	redisRepo "challenge-app/internal/infrastructure/repository/redis"
	infra_storage "challenge-app/internal/infrastructure/storage"
	"challenge-app/internal/presentation/handler"
	handler_interface "challenge-app/internal/presentation/handler/interface"
	"challenge-app/internal/presentation/middleware"
	middleware_interface "challenge-app/internal/presentation/middleware/interface"
	"challenge-app/internal/presentation/router"
	"challenge-app/pkg/email"
	"challenge-app/pkg/security"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type PostgresDSN string

// Providers
func ProvideEnv() *bootstrap.Env {
	return bootstrap.LoadEnv()
}

func ProvideDSN(cfg *bootstrap.Env) PostgresDSN {
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

func ProvideRedisClient(cfg *bootstrap.Env) (*redis.Client, error) {
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
	return validator.New()
}

func ProvideS3Storage(cfg *bootstrap.Env) (*infra_storage.S3Storage, error) {
	return infra_storage.NewS3Storage(
		cfg.Storage.Endpoint,
		cfg.Storage.Region,
		cfg.Storage.Bucket,
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		cfg.Storage.PublicURL,
	)
}

func ProvideTranslator() domainLoc.Translator {
	return localization.NewTranslationService()
}

func GetChallengeCompletionWorker(
	completionRepo repository_interface.ChallengeCompletionRepository,
	challengeRepo repository_interface.ChallengeRepository,
	participantRepo repository_interface.ChallengeParticipantRepository,
	userDayRepo repository_interface.UserDayRepository,
	userRepo repository_interface.UserRepository,
	categoryRepo repository_interface.CategoryRepository,
) *workers.ChallengeCompletionWorker {
	// Run every minute by default
	return workers.NewChallengeCompletionWorker(
		completionRepo,
		challengeRepo,
		participantRepo,
		userDayRepo,
		userRepo,
		categoryRepo,
		1*time.Minute, // Interval
	)
}

// Provider Sets
var EnvProviderSet = wire.NewSet(
	ProvideEnv,
)

var LocalizationProviderSet = wire.NewSet(
	ProvideTranslator,
)

var RealtimeProviderSet = wire.NewSet(
	ws.NewHub,
	ws.NewWSNotifier,
)

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
	redisRepo.NewTempUploadRepository,
	wire.Bind(new(repository_interface.TempUploadRepository), new(*redisRepo.TempUploadRepository)),
	wire.Bind(new(repository_interface.VerificationRepository), new(*redisRepo.VerificationRepository)),
)

var EmailProviderSet = wire.NewSet(
	ProvideEmailService,
	wire.Bind(new(email.EmailService), new(*email.EmailServiceImpl)),
)

var StorageProviderSet = wire.NewSet(
	ProvideS3Storage,
	wire.Bind(new(infra_storage.ObjectStorage), new(*infra_storage.S3Storage)),
)

var RepositoryProviderSet = wire.NewSet(
	postgres.NewUserRepository,
	postgres.NewChallengeRepository,
	postgres.NewCategoryRepository,
	postgres.NewChallengeParticipationRepository,
	postgres.NewChallengeParticipantRepository,
	postgres.NewCommentRepository,
	postgres.NewLikeRepository,
	postgres.NewUserDayRepository,
	postgres.NewFollowRepository,
	postgres.NewPostRepository,
	postgres.NewChallengeCompletionRepository,
	postgres.NewNotificationRepo,
	wire.Bind(new(repository_interface.ChallengeCompletionRepository), new(*postgres.ChallengeCompletionRepository)),
	wire.Bind(new(repository_interface.LikeRepository), new(*postgres.LikeRepository)),
	wire.Bind(new(repository_interface.UserDayRepository), new(*postgres.UserDayRepository)),
	wire.Bind(new(repository_interface.ChallengeInviteRepository), new(*postgres.ChallengeParticipationRepository)),
	wire.Bind(new(repository_interface.ChallengeJoinRequestRepository), new(*postgres.ChallengeParticipationRepository)),
	wire.Bind(new(repository_interface.CommentRepository), new(*postgres.CommentRepository)),
	wire.Bind(new(repository_interface.ChallengeParticipantRepository), new(*postgres.ChallengeParticipantRepository)),
	wire.Bind(new(repository_interface.UserRepository), new(*postgres.UserRepository)),
	wire.Bind(new(repository_interface.ChallengeRepository), new(*postgres.ChallengeRepository)),
	wire.Bind(new(repository_interface.CategoryRepository), new(*postgres.CategoryRepository)),
	wire.Bind(new(repository_interface.FollowRepository), new(*postgres.FollowRepository)),
	wire.Bind(new(repository_interface.PostRepository), new(*postgres.PostRepository)),
)

var ServiceProviderSet = wire.NewSet(
	service.NewUserService,
	service.NewAuthService,
	service.NewChallengeService,
	service.NewFollowService,
	service.NewPostService,
	service.NewUserDayService,
	service.NewTempUploadCleaner,
	service.NewChallengeCompletionService,
	service.NewNotificationService,
	wire.Bind(new(service_interface.ChallengeServicer), new(*service.ChallengeService)),
	wire.Bind(new(service_interface.UserServicer), new(*service.UserService)),
	wire.Bind(new(service_interface.AuthServicer), new(*service.AuthService)),
	wire.Bind(new(service_interface.FollowServicer), new(*service.FollowService)),
	wire.Bind(new(service_interface.PostServicer), new(*service.PostService)),
	wire.Bind(new(service_interface.UserDayServicer), new(*service.UserDayService)),
	wire.Bind(new(service_interface.ChallengeCompletionServicer), new(*service.ChallengeCompletionService)),
	wire.Bind(new(service_interface.NotificationService), new(*service.NotificationServiceImpl)),
)

var HandlerProviderSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewAuthHandler,
	handler.NewChallengeHandler,
	handler.NewFollowHandler,
	handler.NewPostHandler,
	handler.NewUserDayHandler,
	handler.NewChallengeCompletionHandler,
	handler.NewNotificationHandler,
	handler.NewWSNotificationHandler,
	handler.NewDebugHandler,
	wire.Bind(new(handler_interface.DebugHandler), new(*handler.DebugHandlerImpl)),
	wire.Bind(new(handler_interface.UserHandler), new(*handler.UserHandler)),
	wire.Bind(new(handler_interface.AuthHandler), new(*handler.AuthHandler)),
	wire.Bind(new(handler_interface.ChallengeHandler), new(*handler.ChallengeHandler)),
	wire.Bind(new(handler_interface.FollowHandler), new(*handler.FollowHandlerImpl)),
	wire.Bind(new(handler_interface.PostHandler), new(*handler.PostHandler)),
	wire.Bind(new(handler_interface.UserDayHandler), new(*handler.UserDayHandler)),
	wire.Bind(new(handler_interface.ChallengeCompletionHandler), new(*handler.ChallengeCompletionHandler)),
	wire.Bind(new(handler_interface.NotificationHandler), new(*handler.NotificationHandlerImpl)),
	wire.Bind(new(handler_interface.WSNotificationHandler), new(*handler.WSNotificationHandlerImpl)),
)

var MiddlewareProviderSet = wire.NewSet(
	middleware.NewJWTMiddleware,
	middleware.NewLocalizationMiddleware,
	middleware.NewErrorMiddleware,
	wire.Bind(new(middleware_interface.ErrorMiddleware), new(*middleware.ErrorMiddleware)),
	wire.Bind(new(middleware_interface.JWTMiddleware), new(*middleware.JWTMiddleware)),
	wire.Bind(new(middleware_interface.LocalizationMiddleware), new(*middleware.LocalizationMiddleware)),
)

var WorkerProviderSet = wire.NewSet(
	GetChallengeCompletionWorker,
)

// Application
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

type ApplicationContainer struct {
	App    *Application
	Worker *workers.ChallengeCompletionWorker
}

func NewApplicationContainer(app *Application, worker *workers.ChallengeCompletionWorker) *ApplicationContainer {
	return &ApplicationContainer{
		App:    app,
		Worker: worker,
	}
}

// Initialize Router
func InitializeRouter(db *gorm.DB, v *validator.Validate) (*gin.Engine, error) {
	wire.Build(
		EnvProviderSet,
		LocalizationProviderSet,
		StorageProviderSet,
		SecurityProviderSet,
		RepositoryProviderSet,
		RealtimeProviderSet,
		RedisProviderSet,
		EmailProviderSet,
		ServiceProviderSet,
		HandlerProviderSet,
		MiddlewareProviderSet,
		router.SetupRouter,
	)
	return &gin.Engine{}, nil
}

// Initialize Full Application
func InitializeApplication() (*Application, error) {
	wire.Build(
		EnvProviderSet,
		LocalizationProviderSet,
		StorageProviderSet,
		DatabaseProviderSet,
		SecurityProviderSet,
		RepositoryProviderSet,
		RedisProviderSet,
		EmailProviderSet,
		RealtimeProviderSet,
		ServiceProviderSet,
		HandlerProviderSet,
		MiddlewareProviderSet,
		router.SetupRouter,
		NewApplication,
	)
	return &Application{}, nil
}

func InitializeApplicationWithWorker() (*ApplicationContainer, error) {
	wire.Build(
		EnvProviderSet,
		LocalizationProviderSet,
		StorageProviderSet,
		DatabaseProviderSet,
		SecurityProviderSet,
		RepositoryProviderSet,
		RedisProviderSet,
		EmailProviderSet,
		RealtimeProviderSet,
		ServiceProviderSet,
		HandlerProviderSet,
		MiddlewareProviderSet,
		WorkerProviderSet,
		router.SetupRouter,
		NewApplication,
		NewApplicationContainer,
	)
	return &ApplicationContainer{}, nil
}
