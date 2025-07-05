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

func GetBudgetSummary(c *gin.Context) {
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

	summary, serviceErr := services.GetBudgetSummary(c, budgetID, userID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func GetWeeklySummary(c *gin.Context) {
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

	summary, serviceErr := services.GetWeeklySummary(c, userID, startDate, endDate)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func GetMonthlySummary(c *gin.Context) {
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

	summary, serviceErr := services.GetMonthlySummary(c, userID, month)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func GetYearlySummary(c *gin.Context) {
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

	summary, serviceErr := services.GetYearlySummary(c, userID, year)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func GetCategorySummary(c *gin.Context) {
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

	summary, serviceErr := services.GetCategorySummary(c, userID, categoryID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func GetCustomDateRangeSummary(c *gin.Context) {
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

	summary, serviceErr := services.GetCustomDateRangeSummary(c, userID, startDate, endDate)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func GetDailyAverageSummary(c *gin.Context) {
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

	summary, serviceErr := services.GetDailyAverageSummary(c, userID, startDate, endDate)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func GetTopCategories(c *gin.Context) {
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

	categories, serviceErr := services.GetTopCategories(c, userID, limit, transactionType)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetAllCategoriesSummary(c *gin.Context) {
	userID, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	summaries, serviceErr := services.GetAllCategoriesSummary(c, userID)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, summaries)
}
