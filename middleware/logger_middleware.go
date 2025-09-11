package middleware

import (
	"fmt"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/errors"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
       return func(c *gin.Context) {
	       start := time.Now()
	       c.Next()
	       latency := time.Since(start)
	       requestID := c.GetString("request_id")
	       logger := utils.GetLogger()

	       if len(c.Errors) > 0 {
			fmt.Println("Error Middleware Triggered")
		       for _, e := range c.Errors {
				err := e.Err
		       appErr, ok := err.(*errors.AppError)

		       if !ok {
			       appErr = errors.NewInternalError(err, )
		       }

		       logger.Error().			   	
			       Str("request_id", requestID).
			       Int("status", appErr.Code).
			       Dur("latency", latency).
			       Str("method", c.Request.Method).
			       Str("path", c.Request.URL.Path).
			   		
			       Err(appErr.Err).
			       Msg(appErr.Message)

		       c.AbortWithStatusJSON(appErr.Code, gin.H{
			       "message": appErr.Message,
		       })
			}
			return
	       }
	       logger.Info().
		       Str("request_id", requestID).
		       Int("status", c.Writer.Status()).
		       Dur("latency", latency).
		       Str("method", c.Request.Method).
		       Str("path", c.Request.URL.Path).
		       Msg("Request completed")
       }
}
