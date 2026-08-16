package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestTimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		ctx := context.WithValue(c.Request.Context(), "startTime", start)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
