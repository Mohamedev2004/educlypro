package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Try to grab the ID from the incoming React frontend request
		reqID := c.GetHeader("X-Request-ID")

		// 2. If it's missing (e.g., Postman or a third-party webhook), generate one
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// 3. Set it in the Gin context (useful for standard HTTP logging)
		c.Set("X-Request-ID", reqID)

		// 4. Send it back in the response header so the client knows it
		c.Header("X-Request-ID", reqID)

		// 5. CRITICAL STEP FOR WATERMILL:
		// Gin's c.Set() does NOT automatically push values into c.Request.Context().
		// We must manually wrap the standard Go context so your Service layer can read it!
		ctx := context.WithValue(c.Request.Context(), "X-Request-ID", reqID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
