package staff

import (
	"educlypro/shared/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, mainDB *gorm.DB, h *Handler) {
	staff := rg.Group("/staff")
	staff.Use(
		middleware.AuthMiddleware(mainDB),
		middleware.RequireRole("center_owner"),
		middleware.RequireAcademicSetup(mainDB),
	)
	{
		staff.GET("", h.List)
		staff.POST("", h.Create)
		staff.PUT("/:id", h.Update)
		staff.DELETE("/:id", h.Delete)
	}
}
