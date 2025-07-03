package routes

import (
	"github.com/AlsoShantanuBorkar/budget_max/controllers"
	"github.com/AlsoShantanuBorkar/budget_max/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterBudgetRoutes(rg *gin.RouterGroup) {
	budget := rg.Group("/budget")

	budget.Use(middleware.AuthMiddleware())

	budget.GET("/", controllers.GetBudgetsByUserID)
	budget.GET("/:id", controllers.GetBudgetByID)
	budget.POST("/", controllers.CreateBudget)
	budget.DELETE("/:id", controllers.DeleteBudget)
	budget.PUT("/:id", controllers.UpdateBudget)
}
