package subcenters

import (
	"educlypro/shared/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, mainDB *gorm.DB, h *Handler) {
	subcenters := rg.Group("/subcenters")
	subcenters.Use(middleware.AuthMiddleware(mainDB), middleware.RequireRole("center_owner"))
	{
		subcenters.GET("", h.List)
		subcenters.POST("", h.Create)
		subcenters.PUT("/:id", h.Update)
		subcenters.DELETE("/:id", h.Delete)
	}
}
