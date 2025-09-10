package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/services"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReportsControllerInterface interface {
	GetBudgetSummary(c *gin.Context)
	GetWeeklySummary(c *gin.Context)
	GetMonthlySummary(c *gin.Context)
	GetYearlySummary(c *gin.Context)
	GetCategorySummary(c *gin.Context)
	GetCustomDateRangeSummary(c *gin.Context)
	GetDailyAverageSummary(c *gin.Context)
	GetTopCategories(c *gin.Context)
	GetAllCategoriesSummary(c *gin.Context)
}

type ReportsController struct {
	service services.ReportsServiceInterface
}

func NewReportsController(service services.ReportsServiceInterface) *ReportsController {
	return &ReportsController{
		service: service,
	}
}

func (ctrl *ReportsController) GetBudgetSummary(c *gin.Context) {
	userID, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	budgetIDStr := c.Param("budget_id")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}

	summary, serviceErr := ctrl.service.GetBudgetSummary(c, budgetID, userID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (ctrl *ReportsController) GetWeeklySummary(c *gin.Context) {
	userID, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	summary, serviceErr := ctrl.service.GetWeeklySummary(c, userID, startDate, endDate)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (ctrl *ReportsController) GetMonthlySummary(c *gin.Context) {
	userID, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	monthStr := c.Query("month")
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format. Use YYYY-MM"})
		return
	}

	summary, serviceErr := ctrl.service.GetMonthlySummary(c, userID, month)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (ctrl *ReportsController) GetYearlySummary(c *gin.Context) {
	userID, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	yearStr := c.Query("year")
	year, err := time.Parse("2006", yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format. Use YYYY"})
		return
	}

	summary, serviceErr := ctrl.service.GetYearlySummary(c, userID, year)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (ctrl *ReportsController) GetCategorySummary(c *gin.Context) {
	userID, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	categoryIDStr := c.Param("category_id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	summary, serviceErr := ctrl.service.GetCategorySummary(c, userID, categoryID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (ctrl *ReportsController) GetCustomDateRangeSummary(c *gin.Context) {
	userID, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	summary, serviceErr := ctrl.service.GetCustomDateRangeSummary(c, userID, startDate, endDate)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (ctrl *ReportsController) GetDailyAverageSummary(c *gin.Context) {
	userID, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	summary, serviceErr := ctrl.service.GetDailyAverageSummary(c, userID, startDate, endDate)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (ctrl *ReportsController) GetTopCategories(c *gin.Context) {
	userID, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	transactionType := c.Query("type")
	if transactionType != "expense" && transactionType != "income" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction type must be 'expense' or 'income'"})
		return
	}

	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	categories, serviceErr := ctrl.service.GetTopCategories(c, userID, limit, transactionType)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (ctrl *ReportsController) GetAllCategoriesSummary(c *gin.Context) {
	userID, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	summaries, serviceErr := ctrl.service.GetAllCategoriesSummary(c, userID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summaries)
}
