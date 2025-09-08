package routes

import (
	"github.com/gin-gonic/gin"

	"sog.com/goinit/gin/api/common/middleware"
	"sog.com/goinit/gin/api/protocol/http/handler"
	"sog.com/goinit/gin/internal/lib/jwt"
)

// SetupAuthRoutes sets up all authentication routes (Django's authentication/api/urls.py equivalent)
func SetupAuthRoutes(router *gin.Engine, jwtSvc jwt.JWTServiceInterface) {
	authHandler := handler.NewAuthHandlerDI()

	// Authentication API routes group
	auth := router.Group("/api/auth")
	{
		// User registration (POST /api/auth/register/)
		auth.POST("/register/", authHandler.UserRegister)

		// User login (POST /api/auth/login/)
		auth.POST("/login/", authHandler.UserLogin)

		// User logout (GET /api/auth/logout/) - requires authentication
		auth.GET("/logout/", middleware.RequireAuth(jwtSvc), authHandler.UserLogout)

		// OTP verification (POST /api/auth/verify/)
		auth.POST("/verify/", authHandler.VerifyOTP)

		// Delete account (DELETE /api/auth/delete/) - requires authentication
		auth.DELETE("/delete/", middleware.RequireAuth(jwtSvc), authHandler.DeleteAccount)

		// Change password (PUT /api/auth/change-password/) - requires authentication
		auth.PUT("/change-password/", middleware.RequireAuth(jwtSvc), authHandler.ChangePassword)

		// Resend OTP (PUT /api/auth/resend-otp/:id/)
		auth.PUT("/resend-otp/:id/", authHandler.ResendOTP)
	}
}

// SetupUserRoutes sets up user management routes
func SetupUserRoutes(router *gin.Engine, jwtSvc jwt.JWTServiceInterface) {
	userHandler := handler.NewUserHandlerDI()

	// User management API routes group
	user := router.Group("/api/user")
	{
		// Get current user profile (GET /api/user/profile/) - requires authentication
		user.GET("/profile/", middleware.RequireAuth(jwtSvc), userHandler.GetUserProfile)

		// Update current user profile (PUT /api/user/profile/) - requires authentication
		user.PUT("/profile/", middleware.RequireAuth(jwtSvc), userHandler.UpdateUserProfile)

		// Upload/Update profile image (POST /api/user/profile/image/) - requires authentication
		user.POST("/profile/image/", middleware.RequireAuth(jwtSvc), userHandler.UploadProfileImage)

		// Admin routes - requires staff privileges
		admin := user.Group("/admin")
		admin.Use(middleware.RequireAuth(jwtSvc))
		admin.Use(middleware.RequireAdmin())
		{
			// Get all users (GET /api/user/admin/users/) - admin only
			admin.GET("/users/", userHandler.GetAllUsers)

			// Get verified users (GET /api/user/admin/verified/) - admin only
			admin.GET("/verified/", userHandler.GetVerifiedUsers)

			// Get unverified users (GET /api/user/admin/unverified/) - admin only
			admin.GET("/unverified/", userHandler.GetUnverifiedUsers)

			// Get user by ID (GET /api/user/admin/:id/) - admin only
			admin.GET("/:id/", userHandler.GetUserByID)
		}
	}
}

// SetupPasswordResetRoutes registers password reset request and confirm endpoints
func SetupPasswordResetRoutes(router *gin.Engine, publicHost string) {
	
	prh := handler.NewPasswordResetHandlerDI(publicHost)
	auth := router.Group("/api/auth")
	{
		auth.POST("/password-reset/request/", prh.RequestPasswordReset)
		auth.POST("/password-reset/confirm/", prh.ConfirmPasswordReset)
	}
}

// SetupAdminRoutes sets up admin-specific routes (Django admin equivalent)
func SetupAdminRoutes(router *gin.Engine, jwtSvc jwt.JWTServiceInterface) {
	adminHandler := handler.NewAdminHandlerDI()

	// Admin API routes group - requires staff privileges
	admin := router.Group("/api/admin")
	admin.Use(middleware.RequireAuth(jwtSvc))
	admin.Use(middleware.RequireAdmin())
	{
		// User management endpoints
		admin.GET("/stats/", adminHandler.GetUserStats)
		admin.GET("/search/", adminHandler.SearchUsers)

		// User actions
		admin.PUT("/users/:id/activate/", adminHandler.ActivateUser)
		admin.PUT("/users/:id/deactivate/", adminHandler.DeactivateUser)
		admin.PUT("/users/:id/force-verify/", adminHandler.ForceVerifyUser)

		// Bulk operations
		admin.POST("/bulk-email/", adminHandler.SendBulkEmail)
	}
}
