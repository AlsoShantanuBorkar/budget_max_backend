package middleware

import (
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger() gin.HandlerFunc{
	return func(c *gin.Context){
		start :=time.Now()
		c.Next()
		latency:=time.Since(start)
		status :=c.Writer.Status()

		if len(c.Errors)>0{
			for _, e := range c.Errors {
				utils.AppLogger.Error("Request failed",
					zap.Int("status", status),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("error", e.Error()),
					zap.Duration("latency", latency),
					zap.String("client_ip", c.ClientIP()),
				)
			}
		}else {
			utils.AppLogger.Info("Request completed",
				zap.Int("status", status),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Duration("latency", latency),
				zap.String("client_ip", c.ClientIP()),
			)
		}
		
	}
}
