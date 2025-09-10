package routes

import (
	"github.com/AlsoShantanuBorkar/budget_max/controllers"
	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterBudgetRoutes(rg *gin.RouterGroup, ctrl controllers.BudgetControllerInterface, sessionDatabaseService database.SessionDatabaseServiceInterface) {
	budget := rg.Group("/budget")

	budget.Use(middleware.AuthMiddleware(sessionDatabaseService))

	budget.GET("/", ctrl.GetBudgetsByUserID)
	budget.GET("/:id", ctrl.GetBudgetByID)
	budget.POST("/", ctrl.CreateBudget)
	budget.DELETE("/:id", ctrl.DeleteBudget)
	budget.PUT("/:id", ctrl.UpdateBudget)
}
