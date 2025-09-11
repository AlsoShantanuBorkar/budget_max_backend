package utils

import (
	"net/http"

	"github.com/AlsoShantanuBorkar/budget_max/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParseUserID(c *gin.Context) (uuid.UUID, bool) {
	userIdRaw, exists := c.Get("user_id")
	if !exists {
		appErr := errors.NewUnauthorizedError("User ID not found in context", nil, )
		c.Error(appErr)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		c.Abort()
		return uuid.UUID{}, false
	}

	userId, ok := userIdRaw.(uuid.UUID)
	if !ok {
		appErr := errors.NewUnauthorizedError("Invalid user ID type", nil, )
		c.Error(appErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		c.Abort()
		return uuid.UUID{}, false
	}

	return userId, true
}
