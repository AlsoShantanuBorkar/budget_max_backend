package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(sessionService database.SessionDatabaseServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := uuid.Parse(tokenStr)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		session, err := sessionService.GetSessionByToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session Expired"})
			return
		}

		c.Set("user_id", session.UserID)
		c.Next()
	}
}
