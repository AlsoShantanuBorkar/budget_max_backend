package controllers

import (
	"net/http"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/errors"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/services"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionControllerInterface interface {
	CreateTransaction(c *gin.Context)
	UpdateTransaction(c *gin.Context)
	DeleteTransaction(c *gin.Context)
	GetTransactionsByUserID(c *gin.Context)
	GetTransactionByID(c *gin.Context)
	GetTransactionsByBudget(c *gin.Context)
	GetTransactionsByCategory(c *gin.Context)
	GetTransactionsByDateRange(c *gin.Context)
	GetTransactionsByType(c *gin.Context)
	GetTransactionsByAmountRange(c *gin.Context)
	GetTransactionsWithFilters(c *gin.Context)
}
type TransactionController struct {
	service services.TransactionServiceInterface
}

func NewTransactionController(service services.TransactionServiceInterface) *TransactionController {
	return &TransactionController{
		service: service,
	}
}

func (ctrl *TransactionController) CreateTransaction(c *gin.Context) {
	var req models.CreateTransactionRequest
       if err := c.ShouldBindJSON(&req); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }
       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       txn, err := ctrl.service.CreateTransaction(c, &req, userId)
       if err != nil {
	       appErr := errors.NewInternalError(err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusCreated, gin.H{
	       "message": "Transaction created successfully",
	       "data":    txn,
       })
}

func (ctrl *TransactionController) UpdateTransaction(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       var req models.UpdateTransactionRequest
       if err := c.ShouldBindJSON(&req); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       txnIdStr := c.Param("id")
       txnId, err := uuid.Parse(txnIdStr)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid transaction ID format", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       updatedTransaction, serviceErr := ctrl.service.UpdateTransaction(c, &req, txnId, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Transaction updated successfully",
	       "data":    updatedTransaction,
       })
}

func (ctrl *TransactionController) DeleteTransaction(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       txnIdStr := c.Param("id")
       txnId, err := uuid.Parse(txnIdStr)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid transaction ID format", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       serviceErr := ctrl.service.DeleteTransaction(c, txnId, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusNoContent, gin.H{
	       "message": "Transaction deleted successfully",
       })
}

func (ctrl *TransactionController) GetTransactionsByUserID(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       txns, serviceErr := ctrl.service.GetTransactionsByUserID(c, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Transactions fetched successfully",
	       "data":    txns,
       })
}

func (ctrl *TransactionController) GetTransactionByID(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       txnIDStr := c.Param("id")
       txnID, err := uuid.Parse(txnIDStr)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid transaction ID format", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       txn, serviceErr := ctrl.service.GetTransactionByID(c, txnID, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Transaction fetched successfully",
	       "data":    txn,
       })
}

func (ctrl *TransactionController) GetTransactionsByBudget(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       budgetIDStr := c.Param("budget_id")
       budgetID, err := uuid.Parse(budgetIDStr)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid budget ID format", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       txns, serviceErr := ctrl.service.GetTransactionsByBudget(c, budgetID, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Transactions fetched successfully",
	       "data":    txns,
       })
}

func (ctrl *TransactionController) GetTransactionsByCategory(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       categoryIDStr := c.Param("categoryId")
       categoryID, err := uuid.Parse(categoryIDStr)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid category ID format", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       txns, serviceErr := ctrl.service.GetTransactionsByCategory(c, categoryID, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Transactions fetched successfully",
	       "data":    txns,
       })
}

func (ctrl *TransactionController) GetTransactionsByDateRange(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       var req models.DateRangeRequest
       if err := c.ShouldBindJSON(&req); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }
       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       startDate, err := time.Parse(time.RFC3339, req.StartDate)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid start date format", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       endDate, err := time.Parse(time.RFC3339, req.EndDate)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid end date format", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       txns, serviceErr := ctrl.service.GetTransactionsByDateRange(c, startDate, endDate, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Transactions fetched successfully",
	       "data":    txns,
       })
}

func (ctrl *TransactionController) GetTransactionsByType(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       transactionType := c.Param("type")
       if transactionType != "expense" && transactionType != "income" {
	       appErr := errors.NewBadRequestError("Invalid transaction type. Must be 'expense' or 'income'", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       txns, serviceErr := ctrl.service.GetTransactionsByType(c, transactionType, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Transactions fetched successfully",
	       "data":    txns,
       })
}

func (ctrl *TransactionController) GetTransactionsByAmountRange(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       var req models.AmountRangeRequest
       if err := c.ShouldBindJSON(&req); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }
       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       txns, serviceErr := ctrl.service.GetTransactionsByAmountRange(c, req.MinAmount, req.MaxAmount, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Transactions fetched successfully",
	       "data":    txns,
       })
}

func (ctrl *TransactionController) GetTransactionsWithFilters(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       var req models.TransactionFiltersRequest
       if err := c.ShouldBindJSON(&req); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }
       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       filters := make(map[string]interface{})

       if req.BudgetID != nil {
	       budgetID, err := uuid.Parse(*req.BudgetID)
	       if err != nil {
		       appErr := errors.NewBadRequestError("Invalid budget ID format", err)
		       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
		       return
	       }
	       filters["budget_id"] = budgetID
       }

       if req.CategoryID != nil {
	       categoryID, err := uuid.Parse(*req.CategoryID)
	       if err != nil {
		       appErr := errors.NewBadRequestError("Invalid category ID format", err)
		       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
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
		       appErr := errors.NewBadRequestError("Invalid start date format", err)
		       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
		       return
	       }
	       filters["start_date"] = startDate
       }

       if req.EndDate != nil {
	       endDate, err := time.Parse(time.RFC3339, *req.EndDate)
	       if err != nil {
		       appErr := errors.NewBadRequestError("Invalid end date format", err)
		       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
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

       txns, serviceErr := ctrl.service.GetTransactionsWithFilters(c, filters, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Transactions fetched successfully",
	       "data":    txns,
       })
}
