package database

import (
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/errors"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionDatabaseServiceInterface interface {
	CreateTransaction(txn *models.Transaction) error
	GetTransactionsByUser(userID uuid.UUID) ([]*models.Transaction, error)
	UpdateTransaction(id uuid.UUID, updates map[string]any) error
	DeleteTransaction(id uuid.UUID) error
	GetTransactionByID(txnId uuid.UUID, userId uuid.UUID) (*models.Transaction, error)
	GetTransactionsByBudget(userID uuid.UUID, budgetID uuid.UUID) ([]*models.Transaction, error)
	GetTransactionsByCategory(userID uuid.UUID, categoryID uuid.UUID) ([]*models.Transaction, error)
	GetTransactionsByDateRange(userID uuid.UUID, startDate, endDate time.Time) ([]*models.Transaction, error)
	GetTransactionsByType(userID uuid.UUID, transactionType string) ([]*models.Transaction, error)
	GetTransactionsByAmountRange(userID uuid.UUID, minAmount, maxAmount float64) ([]*models.Transaction, error)
	GetTransactionsWithFilters(userID uuid.UUID, filters map[string]interface{}) ([]*models.Transaction, error)
}

type TransactionDatabaseService struct {
	database *gorm.DB
}

func NewTransactionDatabaseService(db *gorm.DB) TransactionDatabaseServiceInterface {
	return &TransactionDatabaseService{database: db}
}

func (s *TransactionDatabaseService) CreateTransaction(txn *models.Transaction) error {
	if err := s.database.Create(txn).Error; err != nil {
		return errors.NewDBError(err)
	}
	return nil
}

func (s *TransactionDatabaseService) GetTransactionsByUser(userID uuid.UUID) ([]*models.Transaction, error) {
	var txns []*models.Transaction
	err := s.database.Where("user_id = ?", userID).Find(&txns).Error
	if err != nil {
		return nil, errors.NewDBError(err)
	}
	return txns, nil
}

func (s *TransactionDatabaseService) UpdateTransaction(id uuid.UUID, updates map[string]any) error {
	if err := s.database.Model(&models.Transaction{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return errors.NewDBError(err)
	}
	return nil
}

func (s *TransactionDatabaseService) DeleteTransaction(id uuid.UUID) error {
	if err := s.database.Delete(&models.Transaction{}, "id = ?", id).Error; err != nil {
		return errors.NewDBError(err)
	}
	return nil
}

func (s *TransactionDatabaseService) GetTransactionByID(txnId uuid.UUID, userId uuid.UUID) (*models.Transaction, error) {
	var txn models.Transaction
	err := s.database.First(&txn, "id = ? AND user_id = ?", txnId, userId).Error
	if err != nil {
		return nil, errors.NewDBError(err)
	}
	return &txn, nil
}

// GetTransactionsByBudget returns all transactions for a specific budget
func (s *TransactionDatabaseService) GetTransactionsByBudget(userID uuid.UUID, budgetID uuid.UUID) ([]*models.Transaction, error) {
	var txns []*models.Transaction
	err := s.database.Where("user_id = ? AND budget_id = ?", userID, budgetID).Find(&txns).Error
	if err != nil {
		return nil, errors.NewDBError(err)
	}
	return txns, nil
}

// GetTransactionsByCategory returns all transactions for a specific category
func (s *TransactionDatabaseService) GetTransactionsByCategory(userID uuid.UUID, categoryID uuid.UUID) ([]*models.Transaction, error) {
	var txns []*models.Transaction
	err := s.database.Where("user_id = ? AND category_id = ?", userID, categoryID).Find(&txns).Error
	if err != nil {
		return nil, errors.NewDBError(err)
	}
	return txns, nil
}

// GetTransactionsByDateRange returns all transactions within a date range
func (s *TransactionDatabaseService) GetTransactionsByDateRange(userID uuid.UUID, startDate, endDate time.Time) ([]*models.Transaction, error) {
	var txns []*models.Transaction
	err := s.database.Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).Find(&txns).Error
	if err != nil {
		return nil, errors.NewDBError(err)
	}
	return txns, nil
}

// GetTransactionsByType returns all transactions of a specific type (expense/income)
func (s *TransactionDatabaseService) GetTransactionsByType(userID uuid.UUID, transactionType string) ([]*models.Transaction, error) {
	var txns []*models.Transaction
	err := s.database.Where("user_id = ? AND type = ?", userID, transactionType).Find(&txns).Error
	if err != nil {
		return nil, errors.NewDBError(err)
	}
	return txns, nil
}

// GetTransactionsByAmountRange returns all transactions within an amount range
func (s *TransactionDatabaseService) GetTransactionsByAmountRange(userID uuid.UUID, minAmount, maxAmount float64) ([]*models.Transaction, error) {
	var txns []*models.Transaction
	err := s.database.Where("user_id = ? AND amount >= ? AND amount <= ?", userID, minAmount, maxAmount).Find(&txns).Error
	if err != nil {
		return nil, errors.NewDBError(err)
	}
	return txns, nil
}

// GetTransactionsWithFilters returns transactions with multiple optional filters
func (s *TransactionDatabaseService) GetTransactionsWithFilters(userID uuid.UUID, filters map[string]interface{}) ([]*models.Transaction, error) {
	query := s.database.Where("user_id = ?", userID)

	if budgetID, ok := filters["budget_id"].(uuid.UUID); ok {
		query = query.Where("budget_id = ?", budgetID)
	}

	if categoryID, ok := filters["category_id"].(uuid.UUID); ok {
		query = query.Where("category_id = ?", categoryID)
	}

	if transactionType, ok := filters["type"].(string); ok {
		query = query.Where("type = ?", transactionType)
	}

	if startDate, ok := filters["start_date"].(time.Time); ok {
		query = query.Where("date >= ?", startDate)
	}

	if endDate, ok := filters["end_date"].(time.Time); ok {
		query = query.Where("date <= ?", endDate)
	}

	if minAmount, ok := filters["min_amount"].(float64); ok {
		query = query.Where("amount >= ?", minAmount)
	}

	if maxAmount, ok := filters["max_amount"].(float64); ok {
		query = query.Where("amount <= ?", maxAmount)
	}

	var txns []*models.Transaction
	err := query.Find(&txns).Error
	if err != nil {
		return nil, errors.NewDBError(err)
	}
	return txns, nil
}
