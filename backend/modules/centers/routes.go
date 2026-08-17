package centers

import (
	"educlypro/shared/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, mainDB *gorm.DB, h *Handler) {
	centers := rg.Group("/centers")
	centers.Use(middleware.AuthMiddleware(mainDB), middleware.RequireRole("super_admin"))
	{
		centers.GET("", h.List)
		centers.POST("", h.Create)
		centers.GET("/:slug", h.GetBySlug)
		centers.POST("/:slug/owner", h.AddOwner)
		centers.POST("/:slug/staff", h.AddStaff)
	}
}
