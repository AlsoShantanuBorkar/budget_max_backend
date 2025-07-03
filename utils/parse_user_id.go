package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParseUserID(c *gin.Context) (uuid.UUID, bool) {
	userIdRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		c.Abort()
		return uuid.UUID{}, false
	}

	userId, ok := userIdRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		c.Abort()
		return uuid.UUID{}, false
	}

	return userId, true
}
