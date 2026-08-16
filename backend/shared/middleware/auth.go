package middleware

import (
	"net/http"
	"strings"
	"time"

	"educlypro/shared/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Token model mirrored here to avoid import cycle — or move to a shared/models package.
type tokenRecord struct {
	Type      string    `gorm:"column:type"`
	Token     string    `gorm:"column:token"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		if cookie, err := c.Cookie("auth_token"); err == nil {
			tokenStr = cookie
		}

		if tokenStr == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					tokenStr = parts[1]
				}
			}
		}

		if tokenStr == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "authentication required")
			c.Abort()
			return
		}

		// 1. Validate JWT signature and claims
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		// 2. Check the token still exists in the DB (revocation check)
		hashedToken := utils.HashSHA256(tokenStr)
		var record tokenRecord
		err = db.Table("tokens").
			Where("token = ? AND type = ? AND expires_at > ?", hashedToken, "access", time.Now()).
			First(&record).Error
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "token has been revoked")
			c.Abort()
			return
		}
		if claims.TokenType != "access" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token type")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
