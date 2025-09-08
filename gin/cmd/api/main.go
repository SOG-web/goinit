// Package main bootstraps the House of Rou API server.
//
// @title           SOG-GO GIN INIT Backend API
// @version         1.0
// @description     Complete user management API with authentication, admin features, and Django equivalent functionality.
// @schemes         http https
// @host            localhost:8080
// @BasePath        /api
//
// @securityDefinitions.apikey Bearer
// @in              header
// @name            Authorization
// @description     Enter JWT Bearer token in the format: Bearer {token}
//
// @securityDefinitions.apikey Session
// @in              cookie
// @name            hor_session
package api

import (
	"log/slog"

	"github.com/SOG-web/goinit/gin/api/common/middleware"
	"github.com/SOG-web/goinit/gin/api/protocol/http/router"
	"github.com/SOG-web/goinit/gin/config"
	docs "github.com/SOG-web/goinit/gin/docs"
	userGorm "github.com/SOG-web/goinit/gin/internal/data/user/model/gorm"
	"github.com/SOG-web/goinit/gin/internal/db"
	"github.com/SOG-web/goinit/gin/internal/di"
	jwtLib "github.com/SOG-web/goinit/gin/internal/lib/jwt"
	pwresetGorm "github.com/SOG-web/goinit/gin/internal/lib/pwreset"
	"github.com/SOG-web/goinit/gin/internal/logger"
	"github.com/SOG-web/goinit/gin/internal/server"
	"gorm.io/gorm"
)


func main() {
	cfg := config.Envs

	// Logger
	lg := logger.New(cfg)
	slog.SetDefault(lg)

	// Configure swagger metadata at runtime
	docs.SwaggerInfo.BasePath = "/api"

	slog.Info("creating db", "cfg", cfg)
	var gdb *gorm.DB
	var err error
	if cfg.DBDriver == "sqlite" {
		gdb, err = db.NewSqliteDb(cfg)
	} else if cfg.DBDriver == "mysql" {
		gdb, err = db.NewMysqlDb(cfg)
	} else if cfg.DBDriver == "postgres" {
		gdb, err = db.NewPostgresDb(cfg)
	} else {
		slog.Error("unsupported db driver", "driver", cfg.DBDriver)
		panic("unsupported db driver")
	}
	if err != nil {
		slog.Error("db error", "err", err)
		return
	}
	slog.Info("db created")

	slog.Info("migrating db")
	// User models
	if err := gdb.AutoMigrate(&userGorm.UserGORM{}); err != nil {
		slog.Error("user migrate error", "err", err)
		return
	}

	// JWT and Password Reset models (only if using database implementations)
	if cfg.UseDatabaseJWT || cfg.UseDatabasePWReset {
		serviceModels := []interface{}{}
		if cfg.UseDatabaseJWT {
			serviceModels = append(serviceModels, &jwtLib.BlacklistedToken{})
		}
		if cfg.UseDatabasePWReset {
			serviceModels = append(serviceModels, &pwresetGorm.PasswordResetToken{})
		}
		if len(serviceModels) > 0 {
			if err := gdb.AutoMigrate(serviceModels...); err != nil {
				slog.Error("service models migrate error", "err", err)
				return
			}
			slog.Info("service models migrated")
		}
	}
	
	slog.Info("db migrated")

	// Initialize DI container
	if err := di.InitContainer(cfg, gdb); err != nil {
		slog.Error("failed to initialize DI container", "err", err)
		return
	}
	slog.Info("DI container initialized")

	slog.Info("creating handlers")
	slog.Info("handlers created")

	slog.Info("creating server")

	// Get JWT service from DI container
	jwtSvc := di.MustResolve[jwtLib.JWTServiceInterface](di.DIContainer)

	deps := router.Dependencies{
		SessionMW:  middleware.NewSessionMiddleware(cfg),
		PublicHost: cfg.PublicHost,
		JWTService: jwtSvc,
	}
	
	srv := server.New(cfg, deps)
	slog.Info("server created")

	slog.Info("running server")
	if err := srv.Run(); err != nil {
		slog.Error("server error", "err", err)
		return
	}
}