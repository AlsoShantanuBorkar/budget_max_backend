package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetBudgetSummary godoc
// @Summary Get budget summary
// @Description Get summary of a specific budget including total expenses, income and net balance
// @Tags reports
// @Accept json
// @Produce json
// @Param budget_id path string true "Budget ID"
// @Success 200 {object} reports.BudgetSummary
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reports/budget/{budget_id} [get]
func GetBudgetSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	budgetIDStr := c.Param("budget_id")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}

	summary, err := services.GetBudgetSummary(c, budgetID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetWeeklySummary godoc
// @Summary Get weekly summary
// @Description Get summary of transactions for a specific week
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} reports.WeeklySummary
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reports/weekly [get]
func GetWeeklySummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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

	summary, err := services.GetWeeklySummary(c, userID.(uuid.UUID), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetMonthlySummary godoc
// @Summary Get monthly summary
// @Description Get summary of transactions for a specific month
// @Tags reports
// @Accept json
// @Produce json
// @Param month query string true "Month (YYYY-MM)"
// @Success 200 {object} reports.MonthlySummary
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reports/monthly [get]
func GetMonthlySummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	monthStr := c.Query("month")
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format. Use YYYY-MM"})
		return
	}

	summary, err := services.GetMonthlySummary(c, userID.(uuid.UUID), month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetYearlySummary godoc
// @Summary Get yearly summary
// @Description Get summary of transactions for a specific year
// @Tags reports
// @Accept json
// @Produce json
// @Param year query string true "Year (YYYY)"
// @Success 200 {object} reports.YearlySummary
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reports/yearly [get]
func GetYearlySummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	yearStr := c.Query("year")
	year, err := time.Parse("2006", yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format. Use YYYY"})
		return
	}

	summary, err := services.GetYearlySummary(c, userID.(uuid.UUID), year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetCategorySummary godoc
// @Summary Get category summary
// @Description Get summary of transactions for a specific category
// @Tags reports
// @Accept json
// @Produce json
// @Param category_id path string true "Category ID"
// @Success 200 {object} reports.CategorySummary
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reports/category/{category_id} [get]
func GetCategorySummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	categoryIDStr := c.Param("category_id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	summary, err := services.GetCategorySummary(c, userID.(uuid.UUID), categoryID)
	if err != nil {
		if err.Error() == "Category not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetCustomDateRangeSummary godoc
// @Summary Get custom date range summary
// @Description Get summary of transactions for a custom date range
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} reports.CustomDateRangeSummary
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reports/custom-range [get]
func GetCustomDateRangeSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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

	summary, err := services.GetCustomDateRangeSummary(c, userID.(uuid.UUID), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetDailyAverageSummary godoc
// @Summary Get daily average summary
// @Description Get daily average of transactions for a date range
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} reports.DailyAverageSummary
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reports/daily-average [get]
func GetDailyAverageSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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

	summary, err := services.GetDailyAverageSummary(c, userID.(uuid.UUID), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetTopCategories godoc
// @Summary Get top categories
// @Description Get top categories by transaction amount
// @Tags reports
// @Accept json
// @Produce json
// @Param type query string true "Transaction type (expense or income)"
// @Param limit query int false "Number of top categories to return (default: 5)"
// @Success 200 {array} reports.TopCategory
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reports/top-categories [get]
func GetTopCategories(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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

	categories, err := services.GetTopCategories(c, userID.(uuid.UUID), limit, transactionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetAllCategoriesSummary godoc
// @Summary Get all categories summary
// @Description Get summary of transactions for all categories
// @Tags reports
// @Accept json
// @Produce json
// @Success 200 {array} reports.CategorySummary
// @Failure 401 {object} map[string]interface{}
// @Router /reports/categories [get]
func GetAllCategoriesSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	summaries, err := services.GetAllCategoriesSummary(c, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summaries)
}
