package database

import (
	appErrors "github.com/AlsoShantanuBorkar/budget_max/errors"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BudgetDatabaseServiceInterface interface {
	CreateBudget(b *models.Budget) error
	GetBudgetsByUser(userID uuid.UUID) ([]models.Budget, error)
	UpdateBudget(id uuid.UUID, updates map[string]any) error
	DeleteBudget(id uuid.UUID) error
	GetBudgetByID(budgetId uuid.UUID, userId uuid.UUID) (*models.Budget, error)
}

type BudgetDatabaseService struct {
	database *gorm.DB
}

func NewBudgetDatabaseService(db *gorm.DB) BudgetDatabaseServiceInterface {
	return &BudgetDatabaseService{database: db}
}

func (s *BudgetDatabaseService) CreateBudget(budget *models.Budget) error {
	if err := s.database.Create(budget).Error; err != nil {
		return appErrors.NewDBError(err)
	}
	return nil
}

func (s *BudgetDatabaseService) GetBudgetsByUser(userID uuid.UUID) ([]models.Budget, error) {
	var budgets []models.Budget
	err := s.database.Where("user_id = ?", userID).Find(&budgets).Error
	if err != nil {
		return nil, appErrors.NewDBError(err)
	}
	return budgets, nil
}

func (s *BudgetDatabaseService) UpdateBudget(id uuid.UUID, updates map[string]any) error {
	if err := s.database.Model(&models.Budget{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return appErrors.NewDBError(err)
	}
	return nil
}

func (s *BudgetDatabaseService) DeleteBudget(id uuid.UUID) error {
	if err := s.database.Delete(&models.Budget{}, "id = ?", id).Error; err != nil {
		return appErrors.NewDBError(err)
	}
	return nil
}

func (s *BudgetDatabaseService) GetBudgetByID(budgetId uuid.UUID, userId uuid.UUID) (*models.Budget, error) {
	var b models.Budget
	err := s.database.First(&b, "id = ? AND user_id = ?", budgetId, userId).Error
	if err != nil {
		return nil, appErrors.NewDBError(err)
	}
	return &b, nil
}
