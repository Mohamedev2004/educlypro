package auth

import (
	"net/http"
	"time"

	"educlypro/config"
	"educlypro/shared/utils"

	"github.com/ThreeDotsLabs/watermill/message" // Added Watermill import
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

const (
	accessCookieName  = "auth_token"
	refreshCookieName = "refresh_token"
	sessionCookieName = "session_exists"
)

// Updated to accept the Publisher and pass it to the Service
func NewHandler(db *gorm.DB, publisher message.Publisher) *Handler {
	service := NewService(db, publisher)
	return &Handler{service: service}
}

// POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	// Passed c.Request.Context()
	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			utils.ErrorResponseWithCode(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password.")
			return
		}
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "login_failed", "Login failed. Please try again.")
		return
	}

	setAuthCookies(c, resp)

	utils.SuccessResponse(c, http.StatusOK, "login successful", publicAuthResponse(resp))
}

// POST /auth/refresh
func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshCookieName)
	if err != nil || refreshToken == "" {
		clearAuthCookies(c)
		utils.ErrorResponseWithCode(c, http.StatusUnauthorized, "refresh_token_required", "Your session has expired. Please sign in again.")
		return
	}

	// Passed c.Request.Context()
	resp, err := h.service.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		clearAuthCookies(c)
		utils.ErrorResponseWithCode(c, http.StatusUnauthorized, "invalid_refresh_token", "Your session has expired. Please sign in again.")
		return
	}

	setAuthCookies(c, resp)
	utils.SuccessResponse(c, http.StatusOK, "token refreshed", gin.H{"token": resp.Token})
}

// POST /auth/logout
func (h *Handler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponseWithCode(c, http.StatusUnauthorized, "unauthorized", "You are not authorized.")
		return
	}

	// Passed c.Request.Context()
	if err := h.service.Logout(c.Request.Context(), userID.(uint)); err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "logout_failed", "Logout failed. Please try again.")
		return
	}

	clearAuthCookies(c)

	utils.SuccessResponse(c, http.StatusOK, "logged out successfully", nil)
}

func setAuthCookies(c *gin.Context, resp *AuthResponse) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(accessCookieName,
		resp.Token,
		int(utils.GetAccessTokenExpiry()/time.Second),
		"/",
		config.Cfg.Cookie.Domain,
		config.Cfg.Cookie.Secure,
		config.Cfg.Cookie.HttpOnly,
	)
	c.SetCookie(refreshCookieName,
		resp.RefreshToken,
		int(utils.GetRefreshTokenExpiry()/time.Second),
		"/",
		config.Cfg.Cookie.Domain,
		config.Cfg.Cookie.Secure,
		config.Cfg.Cookie.HttpOnly,
	)
	c.SetCookie(sessionCookieName,
		"true",
		int(utils.GetRefreshTokenExpiry()/time.Second),
		"/",
		config.Cfg.Cookie.Domain,
		config.Cfg.Cookie.Secure,
		false,
	) // Non-HttpOnly for client-side session detection
}

func clearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(accessCookieName,
		"",
		-1,
		"/",
		config.Cfg.Cookie.Domain,
		config.Cfg.Cookie.Secure,
		config.Cfg.Cookie.HttpOnly,
	)
	c.SetCookie(refreshCookieName,
		"",
		-1,
		"/",
		config.Cfg.Cookie.Domain,
		config.Cfg.Cookie.Secure,
		config.Cfg.Cookie.HttpOnly,
	)
	c.SetCookie(sessionCookieName,
		"",
		-1,
		"/",
		config.Cfg.Cookie.Domain,
		config.Cfg.Cookie.Secure,
		false,
	) // still client-readable
}

func publicAuthResponse(resp *AuthResponse) *AuthResponse {
	return &AuthResponse{
		User:  resp.User,
		Token: resp.Token,
	}
}

// GET /auth/me  (protected)
func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponseWithCode(c, http.StatusUnauthorized, "unauthorized", "You are not authorized.")
		return
	}

	// Passed c.Request.Context()
	resp, err := h.service.Me(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusNotFound, "user_not_found", "User not found.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user fetched", gin.H{"user": resp})
}

// PUT /auth/profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponseWithCode(c, http.StatusUnauthorized, "unauthorized", "You are not authorized.")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	// Passed c.Request.Context()
	resp, err := h.service.UpdateProfile(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		if err.Error() == "email already in use" {
			utils.ValidationErrorResponse(c, http.StatusConflict, "email_taken", "Email is already in use.", map[string]string{"email": "validation.taken"})
			return
		}
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "profile_update_failed", "Profile update failed. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "profile updated", resp)
}

// PUT /auth/password
func (h *Handler) UpdatePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponseWithCode(c, http.StatusUnauthorized, "unauthorized", "You are not authorized.")
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	// Passed c.Request.Context()
	if err := h.service.UpdatePassword(c.Request.Context(), userID.(uint), &req); err != nil {
		if err.Error() == "invalid current password" {
			utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_current_password", "Current password is incorrect.", map[string]string{"currentPassword": "auth.current_password_invalid"})
			return
		}
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "password_update_failed", "Password update failed. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "password updated", nil)
}

// POST /auth/forgot-password
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	// Passed c.Request.Context()
	if err := h.service.ForgotPassword(c.Request.Context(), &req); err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "forgot_password_failed", "Failed to process forgot password request.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "If your email exists, a reset link has been sent.", nil)
}

// POST /auth/reset-password
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	// Passed c.Request.Context()
	if err := h.service.ResetPassword(c.Request.Context(), &req); err != nil {
		if err.Error() == "invalid or expired token" {
			utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_or_expired_token", "This reset link is invalid or has expired.")
			return
		}
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "reset_password_failed", "Password reset failed. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password has been reset successfully.", nil)
}
