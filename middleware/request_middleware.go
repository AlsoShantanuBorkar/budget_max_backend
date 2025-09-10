package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		requestId := uuid.New().String()
		c.Set("request_id", requestId)
		c.Writer.Header().Set("X-Request-ID", requestId)
		c.Next()
	}

}