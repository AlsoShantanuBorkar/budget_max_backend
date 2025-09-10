package routes

import (
	"github.com/AlsoShantanuBorkar/budget_max/controllers"
	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUnprotectedAuthRoutes(rg *gin.RouterGroup, ctrl controllers.AuthControllerInterface) {
	auth := rg.Group("/auth")

	auth.POST("/signup", ctrl.Signup)
	auth.POST("/login", ctrl.Login)
	auth.POST("/refresh", ctrl.RefreshToken)
	auth.POST("/2fa/login", ctrl.LoginWith2FA)
}

func RegisterAuthRoutes(rg *gin.RouterGroup, ctrl controllers.AuthControllerInterface, sessionDatabaseService database.SessionDatabaseServiceInterface) {
	auth := rg.Group("/auth")

	auth.Use(middleware.AuthMiddleware(sessionDatabaseService))
	auth.POST("/logout", ctrl.Logout)
	auth.POST("/2fa/setup", ctrl.Generate2FA)
	auth.POST("/2fa/verify", ctrl.Verify2FA)
	auth.PUT("/2fa/disable", ctrl.Disable2FA)
}
