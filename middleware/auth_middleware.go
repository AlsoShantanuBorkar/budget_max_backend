package middleware

import (
	"strings"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(sessionService database.SessionDatabaseServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			appErr := errors.NewUnauthorizedError("Unauthorized", nil, )
			c.Error(appErr)			
			c.AbortWithStatusJSON(appErr.Code, gin.H{"message": appErr.Message})
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := uuid.Parse(tokenStr)

		if err != nil {
			appErr := errors.NewUnauthorizedError("Unauthorized", err, )
			c.Error(appErr)
			c.AbortWithStatusJSON(appErr.Code, gin.H{"message": appErr.Message})
			return
		}

		session, err := sessionService.GetSessionByToken(token)

		if err != nil {
			appErr := errors.NewUnauthorizedError("Unauthorized", err, )
			c.Error(appErr)
			c.AbortWithStatusJSON(appErr.Code, gin.H{"message": appErr.Message})
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			appErr := errors.NewUnauthorizedError("Session Expired", nil, )
			c.Error(appErr)
			c.AbortWithStatusJSON(appErr.Code, gin.H{"message": appErr.Message})
			return
		}

		c.Set("user_id", session.UserID)
		c.Next()
	}
}
