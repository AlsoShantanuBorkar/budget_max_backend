package services

import (
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/errors"

	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionServiceInterface interface {
	CreateTransaction(c *gin.Context, req *models.CreateTransactionRequest, userId uuid.UUID) (*models.Transaction, *ServiceError)
	UpdateTransaction(c *gin.Context, req *models.UpdateTransactionRequest, txnId uuid.UUID, userId uuid.UUID) (*models.Transaction, *ServiceError)
	DeleteTransaction(c *gin.Context, txnId uuid.UUID, userId uuid.UUID) *ServiceError
	GetTransactionsByUserID(c *gin.Context, userId uuid.UUID) ([]*models.Transaction, *ServiceError)
	GetTransactionByID(c *gin.Context, txnID uuid.UUID, userId uuid.UUID) (*models.Transaction, *ServiceError)
	GetTransactionsByBudget(c *gin.Context, budgetID uuid.UUID, userId uuid.UUID) ([]*models.Transaction, *ServiceError)
	GetTransactionsByCategory(c *gin.Context, categoryID uuid.UUID, userId uuid.UUID) ([]*models.Transaction, *ServiceError)
	GetTransactionsByDateRange(c *gin.Context, startDate, endDate time.Time, userId uuid.UUID) ([]*models.Transaction, *ServiceError)
	GetTransactionsByType(c *gin.Context, transactionType string, userId uuid.UUID) ([]*models.Transaction, *ServiceError)
	GetTransactionsByAmountRange(c *gin.Context, minAmount, maxAmount float64, userId uuid.UUID) ([]*models.Transaction, *ServiceError)
	GetTransactionsWithFilters(c *gin.Context, filters map[string]interface{}, userId uuid.UUID) ([]*models.Transaction, *ServiceError)
}

type TransactionService struct {
	transactionDatabase database.TransactionDatabaseServiceInterface
}

func NewTransactionService(dbService database.TransactionDatabaseServiceInterface) *TransactionService {
	return &TransactionService{
		transactionDatabase: dbService,
	}
}

func (s *TransactionService) CreateTransaction(c *gin.Context, req *models.CreateTransactionRequest, userId uuid.UUID) (*models.Transaction, *ServiceError) {
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		appErr := errors.NewBadRequestError("invalid date format", err,)
		c.Error(appErr)
		return nil, ServiceErrorFromAppError(appErr)
	}

	var categoryID *uuid.UUID
	if req.CategoryIDs != "" {
		catID, err := uuid.Parse(req.CategoryIDs)
		if err != nil {
			appErr := errors.NewBadRequestError("invalid category ID", err,)
			c.Error(appErr)
			return nil, ServiceErrorFromAppError(appErr)
		}
		categoryID = &catID
	}

	var budgetID *uuid.UUID
	if req.BudgetID != "" {
		budID, err := uuid.Parse(req.BudgetID)
		if err != nil {
			appErr := errors.NewBadRequestError("invalid budget ID", err,)
			c.Error(appErr)
			return nil, ServiceErrorFromAppError(appErr)
		}
		budgetID = &budID
	}

	txn := &models.Transaction{
		ID:         uuid.New(),
		UserID:     userId,
		Amount:     req.Amount,
		Type:       req.Type,
		Name:       req.Name,
		Note:       req.Note,
		Date:       date,
		CategoryID: categoryID,
		BudgetID:   budgetID,
		CreatedAt:  time.Now(),
	}

	if err := s.transactionDatabase.CreateTransaction(txn); err != nil {
		appErr := errors.NewDBError(err)
		c.Error(appErr)
		return nil, ServiceErrorFromAppError(appErr)
	}

	return txn, nil
}

func (s *TransactionService) UpdateTransaction(c *gin.Context, req *models.UpdateTransactionRequest, txnId uuid.UUID, userId uuid.UUID) (*models.Transaction, *ServiceError) {
	// Fetch existing transaction to verify ownership
	_, err := s.transactionDatabase.GetTransactionByID(txnId, userId)
	if err != nil {
		appErr := errors.NewNotFoundError("transaction", err,)
		c.Error(appErr)
		return nil, ServiceErrorFromAppError(appErr)
	}

	updates := make(map[string]any)
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Date != nil {
		parsedDate, err := time.Parse(time.RFC3339, *req.Date)
		if err != nil {
			appErr := errors.NewBadRequestError("invalid date format", err)
			c.Error(appErr)
			return nil, ServiceErrorFromAppError(appErr)
		}
		updates["date"] = parsedDate
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}
	if req.CategoryID != nil {
		parsedCategoryID, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			appErr := errors.NewBadRequestError("invalid category ID format", err)
			c.Error(appErr)
			return nil, ServiceErrorFromAppError(appErr)
		}
		updates["category_id"] = parsedCategoryID
	}
	if req.BudgetID != nil {
		parsedBudgetID, err := uuid.Parse(*req.BudgetID)
		if err != nil {
			appErr := errors.NewBadRequestError("invalid budget ID format", err)
			c.Error(appErr)
			return nil, ServiceErrorFromAppError(appErr)
		}
		updates["budget_id"] = parsedBudgetID
	}

	// Save updated transaction
	err = s.transactionDatabase.UpdateTransaction(txnId, updates)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	// Fetch updated transaction
	updatedTransaction, err := s.transactionDatabase.GetTransactionByID(txnId, userId)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	return updatedTransaction, nil
}

func (s *TransactionService) DeleteTransaction(c *gin.Context, txnId uuid.UUID, userId uuid.UUID) *ServiceError {
	// Verify transaction exists and belongs to user
	_, err := s.transactionDatabase.GetTransactionByID(txnId, userId)
	if err != nil {
		appErr := errors.NewNotFoundError("transaction", err)
		c.Error(appErr)
		return ServiceErrorFromAppError(appErr)
	}

	err = s.transactionDatabase.DeleteTransaction(txnId)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return ServiceErrorFromAppError(appErr)
       }

	return nil
}

func (s *TransactionService) GetTransactionsByUserID(c *gin.Context, userId uuid.UUID) ([]*models.Transaction, *ServiceError) {
	txns, err := s.transactionDatabase.GetTransactionsByUser(userId)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	return txns, nil
}

func (s *TransactionService) GetTransactionByID(c *gin.Context, txnID uuid.UUID, userId uuid.UUID) (*models.Transaction, *ServiceError) {
	txn, err := s.transactionDatabase.GetTransactionByID(txnID, userId)
	if err != nil {
		appErr := errors.NewNotFoundError("transaction", err)
		c.Error(appErr)
		return nil, ServiceErrorFromAppError(appErr)
	}

	return txn, nil
}

func (s *TransactionService) GetTransactionsByBudget(c *gin.Context, budgetID uuid.UUID, userId uuid.UUID) ([]*models.Transaction, *ServiceError) {
	txns, err := s.transactionDatabase.GetTransactionsByBudget(userId, budgetID)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	return txns, nil
}

func (s *TransactionService) GetTransactionsByCategory(c *gin.Context, categoryID uuid.UUID, userId uuid.UUID) ([]*models.Transaction, *ServiceError) {
       txns, err := s.transactionDatabase.GetTransactionsByCategory(userId, categoryID)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	return txns, nil
}

func (s *TransactionService) GetTransactionsByDateRange(c *gin.Context, startDate, endDate time.Time, userId uuid.UUID) ([]*models.Transaction, *ServiceError) {
       txns, err := s.transactionDatabase.GetTransactionsByDateRange(userId, startDate, endDate)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	return txns, nil
}

func (s *TransactionService) GetTransactionsByType(c *gin.Context, transactionType string, userId uuid.UUID) ([]*models.Transaction, *ServiceError) {
	if transactionType != "expense" && transactionType != "income" {
		appErr := errors.NewBadRequestError("invalid transaction type", nil)
		c.Error(appErr)
		return nil, ServiceErrorFromAppError(appErr)
	}

       txns, err := s.transactionDatabase.GetTransactionsByType(userId, transactionType)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	return txns, nil
}

func (s *TransactionService) GetTransactionsByAmountRange(c *gin.Context, minAmount, maxAmount float64, userId uuid.UUID) ([]*models.Transaction, *ServiceError) {
	if minAmount < 0 || maxAmount < 0 || minAmount > maxAmount {
		appErr := errors.NewBadRequestError("invalid amount range", nil,)
		c.Error(appErr)
		return nil, ServiceErrorFromAppError(appErr)
	}

       txns, err := s.transactionDatabase.GetTransactionsByAmountRange(userId, minAmount, maxAmount)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	return txns, nil
}

func (s *TransactionService) GetTransactionsWithFilters(c *gin.Context, filters map[string]interface{}, userId uuid.UUID) ([]*models.Transaction, *ServiceError) {
       txns, err := s.transactionDatabase.GetTransactionsWithFilters(userId, filters)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	return txns, nil
}
