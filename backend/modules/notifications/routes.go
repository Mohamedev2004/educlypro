package notifications

import (
	"educlypro/shared/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, h *Handler) {
	notif := rg.Group("/notifications")
	notif.Use(middleware.AuthMiddleware(db))
	{
		notif.GET("", h.List)
		notif.GET("/unread-count", h.UnreadCount)
		notif.PATCH("/:id/read", h.MarkRead)
		notif.PATCH("/read-all", h.MarkAllRead)
		notif.DELETE("/read", middleware.RequireRole("super_admin"), h.DeleteAllRead) // super_admin only
		notif.DELETE("/:id", h.Delete)
	}
}
