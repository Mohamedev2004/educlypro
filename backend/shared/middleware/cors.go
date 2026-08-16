package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(origins []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: origins,

		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},

		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization", "X-Request-ID",
		},

		ExposeHeaders: []string{
			"Content-Length",
		},

		AllowCredentials: true,

		// cache preflight request (better performance)
		MaxAge: 12 * 60 * 60, // 12 hours (in seconds)
	})
}
