package controllers

import (
	"net/http"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/services"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTransation(c *gin.Context) {
	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}
	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}

	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	txn, err := services.CreateTransaction(c, &req, userId)
	if err != nil {
		c.JSON(err.Code, gin.H{"message": err.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaction created successfully",
		"data":    txn,
	})
}

func UpdateTransaction(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	var req models.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	txnIdStr := c.Param("id")
	txnId, err := uuid.Parse(txnIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid transaction ID format",
		})
		return
	}

	updatedTransaction, serviceErr := services.UpdateTransaction(c, &req, txnId, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction updated successfully",
		"data":    updatedTransaction,
	})
}

func DeleteTransaction(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	txnIdStr := c.Param("id")
	txnId, err := uuid.Parse(txnIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid transaction ID format",
		})
		return
	}

	serviceErr := services.DeleteTransaction(c, txnId, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Transaction deleted successfully",
	})
}

func GetTransactionsByUserID(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	txns, serviceErr := services.GetTransactionsByUserID(c, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transactions fetched successfully",
		"data":    txns,
	})
}

func GetTransactionByID(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	txnIDStr := c.Param("id")
	txnID, err := uuid.Parse(txnIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid transaction ID format",
		})
		return
	}

	txn, serviceErr := services.GetTransactionByID(c, txnID, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction fetched successfully",
		"data":    txn,
	})
}

func GetTransactionsByBudget(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	budgetIDStr := c.Param("budgetId")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid budget ID format",
		})
		return
	}

	txns, serviceErr := services.GetTransactionsByBudget(c, budgetID, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transactions fetched successfully",
		"data":    txns,
	})
}

func GetTransactionsByCategory(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	categoryIDStr := c.Param("categoryId")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid category ID format",
		})
		return
	}

	txns, serviceErr := services.GetTransactionsByCategory(c, categoryID, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transactions fetched successfully",
		"data":    txns,
	})
}

func GetTransactionsByDateRange(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	var req models.DateRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}
	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid end date format"})
		return
	}

	txns, serviceErr := services.GetTransactionsByDateRange(c, startDate, endDate, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transactions fetched successfully",
		"data":    txns,
	})
}

func GetTransactionsByType(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	transactionType := c.Param("type")
	if transactionType != "expense" && transactionType != "income" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid transaction type. Must be 'expense' or 'income'",
		})
		return
	}

	txns, serviceErr := services.GetTransactionsByType(c, transactionType, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transactions fetched successfully",
		"data":    txns,
	})
}

func GetTransactionsByAmountRange(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	var req models.AmountRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}
	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}

	txns, serviceErr := services.GetTransactionsByAmountRange(c, req.MinAmount, req.MaxAmount, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transactions fetched successfully",
		"data":    txns,
	})
}

func GetTransactionsWithFilters(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	var req models.TransactionFiltersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}
	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}

	filters := make(map[string]interface{})

	if req.BudgetID != nil {
		budgetID, err := uuid.Parse(*req.BudgetID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid budget ID format"})
			return
		}
		filters["budget_id"] = budgetID
	}

	if req.CategoryID != nil {
		categoryID, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid category ID format"})
			return
		}
		filters["category_id"] = categoryID
	}

	if req.Type != nil {
		filters["type"] = *req.Type
	}

	if req.StartDate != nil {
		startDate, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid start date format"})
			return
		}
		filters["start_date"] = startDate
	}

	if req.EndDate != nil {
		endDate, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid end date format"})
			return
		}
		filters["end_date"] = endDate
	}

	if req.MinAmount != nil {
		filters["min_amount"] = *req.MinAmount
	}

	if req.MaxAmount != nil {
		filters["max_amount"] = *req.MaxAmount
	}

	txns, serviceErr := services.GetTransactionsWithFilters(c, filters, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transactions fetched successfully",
		"data":    txns,
	})
}
