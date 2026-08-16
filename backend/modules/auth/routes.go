package auth

import (
	"educlypro/shared/middleware"

	"github.com/ThreeDotsLabs/watermill/message" // Import Watermill
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Add the publisher to the function signature
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, publisher message.Publisher) {

	// Pass the publisher down to your handler (and eventually your service)
	handler := NewHandler(db, publisher)

	auth := router.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.Refresh)
		auth.POST("/logout", middleware.AuthMiddleware(db), handler.Logout)
		auth.GET("/me", middleware.AuthMiddleware(db), handler.Me)
		auth.PUT("/profile", middleware.AuthMiddleware(db), handler.UpdateProfile)
		auth.PUT("/password", middleware.AuthMiddleware(db), handler.UpdatePassword)
		auth.POST("/forgot-password", handler.ForgotPassword)
		auth.POST("/reset-password", handler.ResetPassword)
	}
}
