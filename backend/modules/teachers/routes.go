package teachers

import (
	"educlypro/shared/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, mainDB *gorm.DB, h *Handler) {
	teachers := rg.Group("/teachers")
	teachers.Use(
		middleware.AuthMiddleware(mainDB),
		middleware.RequireRole("center_owner"),
		middleware.RequireAcademicSetup(mainDB),
	)
	{
		teachers.GET("", h.List)
		teachers.GET("/:slug", h.GetBySlug)
		teachers.POST("", h.Create)
		teachers.PUT("/:id", h.Update)
		teachers.PUT("/:id/subjects", h.UpdateSubjects)
		teachers.PUT("/:id/classes", h.UpdateClasses)
		teachers.DELETE("/:id", h.Delete)
	}
}
