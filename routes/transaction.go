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
	transaction.GET("/budget/:budget_id", controllers.GetTransactionsByBudget)
	transaction.GET("/category/:category_id", controllers.GetTransactionsByCategory)
	transaction.GET("/date-range", controllers.GetTransactionsByDateRange)
	transaction.GET("/type/:type", controllers.GetTransactionsByType)
	transaction.GET("/amount-range", controllers.GetTransactionsByAmountRange)
	transaction.GET("/filters", controllers.GetTransactionsWithFilters)
}
