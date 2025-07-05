package routes

import (
	"github.com/AlsoShantanuBorkar/budget_max/controllers"
	"github.com/AlsoShantanuBorkar/budget_max/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterReportsRoutes(rg *gin.RouterGroup) {
	reportsGroup := rg.Group("/reports")
	reportsGroup.Use(middleware.AuthMiddleware())

	// Budget reports
	reportsGroup.GET("/budget/:budget_id", controllers.GetBudgetSummary)

	// Time-based reports
	reportsGroup.GET("/weekly", controllers.GetWeeklySummary)
	reportsGroup.GET("/monthly", controllers.GetMonthlySummary)
	reportsGroup.GET("/yearly", controllers.GetYearlySummary)
	reportsGroup.GET("/custom-range", controllers.GetCustomDateRangeSummary)
	reportsGroup.GET("/daily-average", controllers.GetDailyAverageSummary)

	// Category reports
	reportsGroup.GET("/category/:category_id", controllers.GetCategorySummary)
	reportsGroup.GET("/categories", controllers.GetAllCategoriesSummary)
	reportsGroup.GET("/top-categories", controllers.GetTopCategories)
}
