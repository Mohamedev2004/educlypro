package logs

import (
	"educlypro/shared/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, mainDB *gorm.DB, h *Handler) {
	logs := rg.Group("/logs")
	// Audit logs span every center and every user's actions — restricted to
	// super_admin, not just hidden from other roles' UI/nav.
	logs.Use(middleware.AuthMiddleware(mainDB), middleware.RequireRole("super_admin"))
	{
		logs.GET("", h.List)
		logs.GET("/chart", h.Chart)
		logs.GET("/export", h.Export)
	}
}
