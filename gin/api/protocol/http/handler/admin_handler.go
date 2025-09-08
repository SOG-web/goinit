package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SOG-web/gin/api/common/dto"
	userService "github.com/SOG-web/gin/internal/app/user"
	"github.com/SOG-web/gin/internal/di"
	userModel "github.com/SOG-web/gin/internal/domain/user/model"
)

type AdminHandler struct {
	userService *userService.UserService
}


func NewAdminHandlerDI() *AdminHandler {
	userSvc := di.GetUserService()
	return &AdminHandler{
		userService: userSvc,
	}
}

// GetUserStats returns user statistics (Django admin equivalent)
// @Summary Get User Statistics
// @Description Get comprehensive user statistics for admin dashboard
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.UserStatsResponse "User statistics retrieved successfully"
// @Failure 401 {object} dto.AuthErrorResponse "Unauthorized - invalid or missing token"
// @Failure 403 {object} dto.AuthErrorResponse "Forbidden - admin access required"
// @Failure 500 {object} dto.AuthErrorResponse "Internal server error"
// @Router /admin/stats [get]
func (h *AdminHandler) GetUserStats(c *gin.Context) {
	stats, err := h.userService.GetUserStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AuthErrorResponse{
			Error:      err.Error(),
			Success:    false,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserStatsResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Data:       stats,
	})
}

// SearchUsers searches for users by query (Django admin equivalent)
// @Summary Search Users
// @Description Search for users by email, username, first name, or last name
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param q query string true "Search query (email, username, first name, or last name)"
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of users per page" default(20)
// @Success 200 {object} dto.GetUsersResponse "Users found successfully"
// @Failure 400 {object} dto.AuthErrorResponse "Search query is required"
// @Failure 401 {object} dto.AuthErrorResponse "Unauthorized - invalid or missing token"
// @Failure 403 {object} dto.AuthErrorResponse "Forbidden - admin access required"
// @Failure 500 {object} dto.AuthErrorResponse "Internal server error"
// @Router /admin/users/search [get]
func (h *AdminHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
			Error:      "Search query is required",
			Success:    false,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// Get pagination parameters
	limit := 50 // Default limit
	offset := 0 // Default offset

	users, err := h.userService.SearchUsers(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AuthErrorResponse{
			Error:      err.Error(),
			Success:    false,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	usersData := make([]*dto.UserData, len(users))
	for i, user := range users {
		usersData[i] = h.userModelToDTO(user)
	}

	c.JSON(http.StatusOK, dto.GetUsersResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Data:       usersData,
		Count:      len(usersData),
	})
}

// ActivateUser activates a user account (Django admin equivalent)
// @Summary Activate User
// @Description Activate a user account (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} dto.AdminActionResponse "User activated successfully"
// @Failure 400 {object} dto.AuthErrorResponse "Invalid user ID"
// @Failure 401 {object} dto.AuthErrorResponse "Unauthorized - invalid or missing token"
// @Failure 403 {object} dto.AuthErrorResponse "Forbidden - admin access required"
// @Failure 404 {object} dto.AuthErrorResponse "User not found"
// @Failure 500 {object} dto.AuthErrorResponse "Internal server error"
// @Router /admin/users/{id}/activate [post]
func (h *AdminHandler) ActivateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
			Error:      "User ID is required",
			Success:    false,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	err := h.userService.ActivateUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
			Error:      err.Error(),
			Success:    false,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, dto.AdminActionResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "User activated successfully",
	})
}

// DeactivateUser deactivates a user account (Django admin equivalent)
// @Summary Deactivate User
// @Description Deactivate a user account (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} dto.AdminActionResponse "User deactivated successfully"
// @Failure 400 {object} dto.AuthErrorResponse "Invalid user ID"
// @Failure 401 {object} dto.AuthErrorResponse "Unauthorized - invalid or missing token"
// @Failure 403 {object} dto.AuthErrorResponse "Forbidden - admin access required"
// @Failure 404 {object} dto.AuthErrorResponse "User not found"
// @Failure 500 {object} dto.AuthErrorResponse "Internal server error"
// @Router /admin/users/{id}/deactivate [post]
func (h *AdminHandler) DeactivateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
			Error:      "User ID is required",
			Success:    false,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	err := h.userService.DeactivateUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
			Error:      err.Error(),
			Success:    false,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, dto.AdminActionResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "User deactivated successfully",
	})
}



// ForceVerifyUser forces verification of a user (Django admin equivalent)
// @Summary Force Verify User
// @Description Force verify a user account without OTP (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} dto.AdminActionResponse "User verified successfully"
// @Failure 400 {object} dto.AuthErrorResponse "Invalid user ID"
// @Failure 401 {object} dto.AuthErrorResponse "Unauthorized - invalid or missing token"
// @Failure 403 {object} dto.AuthErrorResponse "Forbidden - admin access required"
// @Failure 404 {object} dto.AuthErrorResponse "User not found"
// @Failure 500 {object} dto.AuthErrorResponse "Internal server error"
// @Router /admin/users/{id}/force-verify [post]
func (h *AdminHandler) ForceVerifyUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
			Error:      "User ID is required",
			Success:    false,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	err := h.userService.ForceVerifyUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
			Error:      err.Error(),
			Success:    false,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, dto.AdminActionResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "User verified successfully",
	})
}

// SendBulkEmail sends email to multiple users (Django equivalent)
// @Summary Send Bulk Email
// @Description Send email to multiple users at once (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.BulkEmailRequest true "Bulk email details"
// @Success 200 {object} dto.AdminActionResponse "Emails sent successfully"
// @Failure 400 {object} dto.AuthErrorResponse "Invalid request format"
// @Failure 401 {object} dto.AuthErrorResponse "Unauthorized - invalid or missing token"
// @Failure 403 {object} dto.AuthErrorResponse "Forbidden - admin access required"
// @Failure 500 {object} dto.AuthErrorResponse "Internal server error"
// @Router /admin/bulk-email [post]
func (h *AdminHandler) SendBulkEmail(c *gin.Context) {
	var req dto.BulkEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
			Error:      err.Error(),
			Success:    false,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	err := h.userService.SendBulkEmail(req.UserIDs, req.Subject, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AuthErrorResponse{
			Error:      err.Error(),
			Success:    false,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.AdminActionResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "Bulk emails sent successfully",
	})
}



// Helper function to convert user model to DTO
func (h *AdminHandler) userModelToDTO(user *userModel.User) *dto.UserData {
	return &dto.UserData{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		IsActive:    user.IsActive,
		IsVerified:  user.IsVerified,
		DateJoined:  user.DateJoined,
		LastLogin:   user.LastLogin,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
