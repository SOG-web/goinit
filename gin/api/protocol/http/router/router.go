package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"sog.com/goinit/gin/api/protocol/http/handler"
	"sog.com/goinit/gin/api/protocol/http/routes"
	"sog.com/goinit/gin/internal/di"
	jwtLib "sog.com/goinit/gin/internal/lib/jwt"
)

type Dependencies struct {
	SessionMW            gin.HandlerFunc
	PublicHost           string
	JWTService           jwtLib.JWTServiceInterface
}

func New(deps Dependencies) *gin.Engine {
	r := gin.Default()

	// jwtSvc := di.MustResolve[jwt.JWTServiceInterface](di.DIContainer)
	redisClient, _ := di.Resolve[*redis.Client](di.DIContainer)

	// Limit multipart memory to 16 MiB (tunable)
	r.MaxMultipartMemory = 16 << 20

	// Session middleware
	if deps.SessionMW != nil {
		r.Use(deps.SessionMW)
	}

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	if redisClient != nil {
		r.GET("/health", handler.HealthWithRedis(redisClient))
	} else {
		r.GET("/health", handler.Health)
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/doc.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

	// Serve static uploads (profile images, etc.)
	r.Static("/uploads", "./uploads")

	// Setup all routes
	if deps.JWTService != nil {
		setupAllRoutes(r, deps.JWTService, deps.PublicHost)
	}

	return r
}

// SetupAllRoutes sets up all API routes
func setupAllRoutes(router *gin.Engine, jwtSvc jwtLib.JWTServiceInterface, publicHost string) {
	// Authentication routes
	routes.SetupAuthRoutes(router, jwtSvc)

	// User management routes
	routes.SetupUserRoutes(router, jwtSvc)

	// Password reset routes
	routes.SetupPasswordResetRoutes(router, publicHost)

	// Admin routes
	routes.SetupAdminRoutes(router, jwtSvc)

	// Real-time routes (SSE and WebSocket)
	routes.SetupSSERoutes(router)
	routes.SetupWSRoutes(router)

	
}