package routes

import (
	"github.com/AlsoShantanuBorkar/budget_max/controllers"
	"github.com/AlsoShantanuBorkar/budget_max/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterTransactionRoutes(rg *gin.RouterGroup) {
	transaction := rg.Group("/transaction")

	transaction.Use(middleware.AuthMiddleware())

	transaction.GET("/", controllers.GetTransactionsByUserID)
	transaction.GET("/:id", controllers.GetTransactionByID)
	transaction.POST("/", controllers.CreateTransation)
	transaction.DELETE("/:id", controllers.DeleteTransaction)
	transaction.PUT("/:id", controllers.UpdateTransaction)
}
