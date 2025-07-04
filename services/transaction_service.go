package services

import (
	"net/http"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTransaction(c *gin.Context, req *models.CreateTransactionRequest, userId uuid.UUID) (*models.Transaction, *ServiceError) {
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, NewServiceError(http.StatusBadRequest, "invalid date format")
	}

	var categoryID *uuid.UUID
	if req.CategoryIDs != "" {
		catID, err := uuid.Parse(req.CategoryIDs)
		if err != nil {
			return nil, NewServiceError(http.StatusBadRequest, "invalid category ID")
		}
		categoryID = &catID
	}

	var budgetID *uuid.UUID
	if req.BudgetID != "" {
		budID, err := uuid.Parse(req.BudgetID)
		if err != nil {
			return nil, NewServiceError(http.StatusBadRequest, "invalid budget ID")
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

	if err := database.CreateTransation(txn); err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to create transaction")
	}

	return txn, nil
}

func UpdateTransaction(c *gin.Context, req *models.UpdateTransactionRequest, txnId uuid.UUID, userId uuid.UUID) (*models.Transaction, *ServiceError) {
	// Fetch existing transaction to verify ownership
	_, err := database.GetTransactionByID(txnId, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusNotFound, "transaction not found")
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
			return nil, NewServiceError(http.StatusBadRequest, "invalid date format")
		}
		updates["date"] = parsedDate
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}
	if req.CategoryID != nil {
		parsedCategoryID, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return nil, NewServiceError(http.StatusBadRequest, "invalid category ID format")
		}
		updates["category_id"] = parsedCategoryID
	}
	if req.BudgetID != nil {
		parsedBudgetID, err := uuid.Parse(*req.BudgetID)
		if err != nil {
			return nil, NewServiceError(http.StatusBadRequest, "invalid budget ID format")
		}
		updates["budget_id"] = parsedBudgetID
	}

	// Save updated transaction
	err = database.UpdateTransaction(txnId, updates)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to update transaction")
	}

	// Fetch updated transaction
	updatedTransaction, err := database.GetTransactionByID(txnId, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to fetch updated transaction")
	}

	return updatedTransaction, nil
}

func DeleteTransaction(c *gin.Context, txnId uuid.UUID, userId uuid.UUID) *ServiceError {
	// Verify transaction exists and belongs to user
	_, err := database.GetTransactionByID(txnId, userId)
	if err != nil {
		return NewServiceError(http.StatusNotFound, "transaction not found")
	}

	err = database.DeleteTransaction(txnId)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to delete transaction")
	}

	return nil
}

func GetTransactionsByUserID(c *gin.Context, userId uuid.UUID) ([]*models.Transaction, *ServiceError) {
	txns, err := database.GetTransactionsByUser(userId)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to fetch transactions")
	}

	return txns, nil
}

func GetTransactionByID(c *gin.Context, txnID uuid.UUID, userId uuid.UUID) (*models.Transaction, *ServiceError) {
	txn, err := database.GetTransactionByID(txnID, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusNotFound, "transaction not found")
	}

	return txn, nil
}
