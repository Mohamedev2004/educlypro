package academic

import (
	"educlypro/shared/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, mainDB *gorm.DB, h *Handler) {
	academic := rg.Group("/academic")
	academic.Use(middleware.AuthMiddleware(mainDB), middleware.RequireRole("center_owner"))
	{
		academic.GET("/grades", h.Tree)
		academic.POST("/grades", h.AddGrade)
		academic.DELETE("/grades/:id", h.RemoveGrade)
		academic.POST("/grades/:id/majors", h.AddMajor)
		academic.DELETE("/majors/:id", h.RemoveMajor)
		academic.POST("/majors/:id/subjects", h.AddSubject)
		academic.DELETE("/subjects/:id", h.RemoveSubject)
	}
}
