package routes

import (
	"github.com/AlsoShantanuBorkar/budget_max/controllers"
	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterTransactionRoutes(rg *gin.RouterGroup, ctrl controllers.TransactionControllerInterface, sessionDatabaseService database.SessionDatabaseServiceInterface) {
	transaction := rg.Group("/transaction")

	transaction.Use(middleware.AuthMiddleware(sessionDatabaseService))

	transaction.GET("/", ctrl.GetTransactionsByUserID)
	transaction.POST("/", ctrl.CreateTransaction)
	transaction.GET("/date-range", ctrl.GetTransactionsByDateRange)
	transaction.GET("/amount-range", ctrl.GetTransactionsByAmountRange)
	transaction.GET("/filters", ctrl.GetTransactionsWithFilters)

	transaction.GET("/budget/:budget_id", ctrl.GetTransactionsByBudget)
	transaction.GET("/category/:category_id", ctrl.GetTransactionsByCategory)
	transaction.GET("/type/:type", ctrl.GetTransactionsByType)
	transaction.GET("/id/:id", ctrl.GetTransactionByID)
	transaction.PUT("/id/:id", ctrl.UpdateTransaction)
	transaction.DELETE("/id/:id", ctrl.DeleteTransaction)
}
