package routes

import (
	"github.com/AlsoShantanuBorkar/budget_max/controllers"
	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterReportsRoutes(rg *gin.RouterGroup, ctrl controllers.ReportsControllerInterface, sessionDatabaseService database.SessionDatabaseServiceInterface) {
	reportsGroup := rg.Group("/reports")
	reportsGroup.Use(middleware.AuthMiddleware(sessionDatabaseService))

	// Budget reports
	reportsGroup.GET("/budget/:budget_id", ctrl.GetBudgetSummary)

	// Time-based reports
	reportsGroup.GET("/weekly", ctrl.GetWeeklySummary)
	reportsGroup.GET("/monthly", ctrl.GetMonthlySummary)
	reportsGroup.GET("/yearly", ctrl.GetYearlySummary)
	reportsGroup.GET("/custom-range", ctrl.GetCustomDateRangeSummary)
	reportsGroup.GET("/daily-average", ctrl.GetDailyAverageSummary)

	// Category reports
	reportsGroup.GET("/category/:category_id", ctrl.GetCategorySummary)
	reportsGroup.GET("/categories", ctrl.GetAllCategoriesSummary)
	reportsGroup.GET("/top-categories", ctrl.GetTopCategories)
}
