package notifications

import (
	"educlypro/shared/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Notifications are a super_admin / center_owner-only feature — no other
// role has an inbox.
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, h *Handler) {
	notif := rg.Group("/notifications")
	notif.Use(middleware.AuthMiddleware(db), middleware.RequireRole("super_admin", "center_owner"))
	{
		notif.GET("", h.List)
		notif.GET("/unread-count", h.UnreadCount)
		notif.PATCH("/:id/read", h.MarkRead)
		notif.PATCH("/read-all", h.MarkAllRead)
		notif.DELETE("/read", h.DeleteAllRead)
		notif.DELETE("/:id", h.Delete)
	}
}
