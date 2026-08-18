package routes

import (
	"educlypro/config"
	"educlypro/modules/academic"
	"educlypro/modules/auth"
	"educlypro/modules/centers"
	"educlypro/modules/logs"
	"educlypro/modules/notifications"
	"educlypro/modules/notifications/delivery"
	"educlypro/modules/staff"
	"educlypro/modules/subcenters"
	"educlypro/modules/teachers"
	"educlypro/shared/database"
	"educlypro/shared/middleware"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, db *gorm.DB, publisher message.Publisher) {

	// ==========================================
	// GLOBAL MIDDLEWARE
	// ==========================================
	r.Use(middleware.CORS(config.Cfg.Server.AllowedOrigins))
	r.Use(middleware.RateLimiterMiddleware(rate.Limit(5), 10, publisher))
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RequestTimingMiddleware())

	// ==========================================
	// API ROUTES
	// ==========================================
	api := r.Group("/api")

	// ------------------------------------------
	// VERSION 1 (Current)
	// ------------------------------------------
	v1 := api.Group("/v1")

	auth.RegisterRoutes(v1, db, publisher)

	// Notifications Wiring (HTTP only)
	notifRepo := notifications.NewRepository(db)
	userResolver := auth.NewUserResolver(db)
	dispatchers := []delivery.Dispatcher{
		delivery.NewInAppDispatcher(db),
		delivery.NewEmailDispatcher(),
	}
	notifService := notifications.NewService(notifRepo, userResolver, dispatchers, publisher)
	notifHandler := notifications.NewHandler(notifService)

	notifications.RegisterRoutes(v1, db, notifHandler)

	// Logs Wiring (HTTP only) - reads from LogsDB, auth via MainDB
	logRepo := logs.NewRepository(database.LogsDB)
	logService := logs.NewService(logRepo)
	logHandler := logs.NewHandler(logService)
	logs.RegisterRoutes(v1, db, logHandler)

	// Staff Wiring (HTTP only) - center_owner managing their own center's staff
	staffRepo := staff.NewRepository(db)
	staffService := staff.NewService(staffRepo, auth.NewRepository(db), publisher)
	staffHandler := staff.NewHandler(staffService)
	staff.RegisterRoutes(v1, db, staffHandler)

	// Sub-centers Wiring (HTTP only) - center_owner managing their own
	// center's sub-centers; also consumed by the centers module below so a
	// super_admin can populate a sub-center picker for an arbitrary center.
	subCentersRepo := subcenters.NewRepository(db)
	subCentersService := subcenters.NewService(subCentersRepo)
	subCentersHandler := subcenters.NewHandler(subCentersService)
	subcenters.RegisterRoutes(v1, db, subCentersHandler)

	// Centers Wiring (HTTP only) - super_admin managing all centers
	centersRepo := centers.NewRepository(db)
	centersService := centers.NewService(centersRepo, staffService, subCentersService)
	centersHandler := centers.NewHandler(centersService)
	centers.RegisterRoutes(v1, db, centersHandler)

	// Academic Wiring (HTTP only) - center_owner managing their own center's
	// grade/major/subject structure.
	academicRepo := academic.NewRepository(db)
	academicService := academic.NewService(academicRepo)
	academicHandler := academic.NewHandler(academicService)
	academic.RegisterRoutes(v1, db, academicHandler)

	// Teachers Wiring (HTTP only) - center_owner managing their own center's
	// teachers (not auth.User accounts — plain records, many-to-many with
	// academic.Subject).
	teachersRepo := teachers.NewRepository(db)
	teachersService := teachers.NewService(teachersRepo)
	teachersHandler := teachers.NewHandler(teachersService)
	teachers.RegisterRoutes(v1, db, teachersHandler)
	// courses.RegisterCourseRoutes(v1, db, publisher)
	// students.RegisterStudentRoutes(v1, db, publisher)
	// parents.RegisterParentRoutes(v1, db, publisher)

	// ------------------------------------------
	// VERSION 2 (Future Example)
	// ------------------------------------------
	// When you eventually build V2, you group it here.
	// Modules that haven't changed can just be mapped to both!
	//
	// v2 := api.Group("/v2")
	// auth.RegisterRoutes(v2, db, publisher) // Auth didn't change
	// students_v2.RegisterStudentRoutes(v2, db, publisher) // Students V2 changed!
}
