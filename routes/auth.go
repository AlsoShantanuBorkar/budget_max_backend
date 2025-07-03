package routes

import (
	"github.com/AlsoShantanuBorkar/budget_max/controllers"
	"github.com/AlsoShantanuBorkar/budget_max/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")

	auth.POST("/signup", controllers.Signup)
	auth.POST("/login", controllers.Login)
	auth.POST("/refresh", controllers.RefreshToken)
	auth.POST("/2fa/login", controllers.LoginWith2FA)

	// Protected route - needs session token
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/logout", controllers.Logout)
	auth.POST("/2fa/setup", controllers.Generate2FA)
	auth.POST("/2fa/verify", controllers.Verify2FA)
	auth.PUT("/2fa/disable", controllers.Disable2FA)
}
