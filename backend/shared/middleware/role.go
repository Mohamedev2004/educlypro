package middleware

import (
	"net/http"

	"educlypro/shared/utils"

	"github.com/gin-gonic/gin"
)

// RequireRole allows only specific roles to access a route
// Usage: RequireRole("super_admin") or RequireRole("center_owner", "center_scanner")
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		for _, r := range allowedRoles {
			if r == role {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, "you don't have permission to access this resource")
		c.Abort()
	}
}
