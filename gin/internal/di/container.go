// Package di provides dependency injection setup for the application.
package di

import (
	"log/slog"
	"time"

	"github.com/SOG-web/goinit/gin/config"
	"github.com/SOG-web/goinit/gin/internal/app/user"
	dataRepo "github.com/SOG-web/goinit/gin/internal/data/user/repo"
	"github.com/SOG-web/goinit/gin/internal/domain/user/repo"
	"github.com/SOG-web/goinit/gin/internal/lib/email"
	jwtLib "github.com/SOG-web/goinit/gin/internal/lib/jwt"
	"github.com/SOG-web/goinit/gin/internal/lib/pwreset"
	"github.com/SOG-web/goinit/gin/internal/lib/storage"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// DIContainer is the global DI container singleton.
var DIContainer *Container

// InitContainer initializes the DI container with all dependencies.
func InitContainer(cfg config.Config, gdb *gorm.DB) error {
	slog.Info("initializing DI container")

	slog.Info("creating services")
	// Email service configuration
	emailConfig := email.EmailConfig{
		Host:     cfg.EmailHost,
		Port:     cfg.EmailPort,
		Username: cfg.EmailUsername,
		Password: cfg.EmailPassword,
		From:     cfg.EmailFrom,
	}

	// Email service (using factory pattern)
	var emailService email.EmailServiceInterface
	var err error
	emailService, err = email.NewEmailServiceFactory(
		emailConfig,
		cfg.UseLocalEmail,
		cfg.EmailLogPath, // Configurable log file path
	)
	if err != nil {
		slog.Error("failed to create email service", "err", err)
		return err
	}

	if cfg.UseLocalEmail {
		slog.Info("using local email service - emails will be logged to ./logs/emails.log")
	} else {
		slog.Info("using production email service")
	}

	// Storage initialization
	var store storage.Storage
	switch cfg.StorageBackend {
	case "s3":
		slog.Info("initializing s3 storage")
		s3Store, err := storage.NewS3Storage(storage.Config{
			Backend:           "s3",
			S3Endpoint:        cfg.S3Endpoint,
			S3Region:          cfg.S3Region,
			S3Bucket:          cfg.S3Bucket,
			S3AccessKeyID:     cfg.S3AccessKeyID,
			S3SecretAccessKey: cfg.S3SecretAccessKey,
			S3UseSSL:          cfg.S3UseSSL,
			S3ForcePathStyle:  cfg.S3ForcePathStyle,
			S3PublicBaseURL:   cfg.S3PublicBaseURL,
		})
		if err != nil {
			slog.Error("failed to init s3 storage, aborting", "err", err)
			return err
		}
		store = s3Store
	default:
		slog.Info("initializing local storage")
		store = storage.NewLocalStorage(cfg.UploadBaseDir, cfg.UploadPublicBaseURL)
	}

	// Redis configuration (only if needed)
	var redisClient *redis.Client
	if !cfg.UseDatabaseJWT || !cfg.UseDatabasePWReset {
		slog.Info("connecting to redis")
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})
		slog.Info("redis connected")
	}

	// JWT service configuration (using factory)
	jwtService := jwtLib.NewJWTServiceFactory(
		cfg.JWTSecret,
		24*time.Hour,  // Access token expiry
		720*time.Hour, // Refresh token expiry (30 days)
		redisClient,   // Redis client for token blacklisting (nil if using database)
		gdb,           // Database connection
	)
	slog.Info("jwt service created")

	// Password reset service (using factory)
	pwResetService := pwreset.NewPasswordResetServiceFactory(
		redisClient, // Redis client (nil if using database)
		gdb,         // Database connection
		time.Hour,   // TTL
	)
	
	c := New()

	// Register database
	if err := Register[*gorm.DB](c, func() *gorm.DB { return gdb }, Singleton); err != nil {
		return err
	}

	// Register email service
	if err := Provide[email.EmailServiceInterface](c, emailService); err != nil {
		return err
	}

	// Register JWT service
	if err := Provide[jwtLib.JWTServiceInterface](c, jwtService); err != nil {
		return err
	}
	
	// Register password reset service
	if err := Provide[pwreset.PasswordResetServiceInterface](c, pwResetService); err != nil {
		return err
	}

	// Register storage
	if err := Register[storage.Storage](c, func() storage.Storage { return store }, Singleton); err != nil {
		return err
	}

	// Register user repository
	if err := Register[repo.UserRepository](c, func(db *gorm.DB) repo.UserRepository {
		return dataRepo.NewGormUserRepository(db)
	}, Singleton); err != nil {
		return err
	}

	// Register user service
	if err := Register[*user.UserService](c, func(userRepo repo.UserRepository, emailSvc email.EmailServiceInterface) *user.UserService {
		return user.NewUserService(userRepo, emailSvc)
	}, Singleton); err != nil {
		return err
	}

	// TODO: Add more registrations for other services/repos as needed

	DIContainer = c
	return nil
}

// GetUserService resolves the user service from the container.
func GetUserService() *user.UserService {
	return MustResolve[*user.UserService](DIContainer)
}

// GetUserRepository resolves the user repository from the container.
func GetUserRepository() repo.UserRepository {
	return MustResolve[repo.UserRepository](DIContainer)
}

// TODO: Add getters for other services/repos
