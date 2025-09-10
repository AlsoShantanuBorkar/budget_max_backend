package database

import (
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type BudgetRepository interface {
	CreateBudget(b *models.Budget) error
	GetBudgetsByUser(userID uuid.UUID) ([]models.Budget, error)
	UpdateBudget(id uuid.UUID, updates map[string]any) error
	DeleteBudget(id uuid.UUID) error
	GetBudgetByID(budgetId uuid.UUID, userId uuid.UUID) (*models.Budget, error)
}

type BudgetService struct {	
	database *gorm.DB
}

func NewBudgetService(db *gorm.DB) *BudgetService {
	return &BudgetService{database: db}
}

func (s *BudgetService) CreateBudget(budget *models.Budget) error {
	return s.database.Create(budget).Error
}

func (s *BudgetService) GetBudgetsByUser(userID uuid.UUID) ([]models.Budget, error) {
	var budgets []models.Budget
	err := s.database.Where("user_id = ?", userID).Find(&budgets).Error
	return budgets, err
}

func (s *BudgetService) UpdateBudget(id uuid.UUID, updates map[string]any) error {
	return s.database.Model(&models.Budget{}).Where("id = ?", id).Updates(updates).Error
}

func (s *BudgetService) DeleteBudget(id uuid.UUID) error {
	return s.database.Delete(&models.Budget{}, "id = ?", id).Error
}

func (s *BudgetService) GetBudgetByID(budgetId uuid.UUID, userId uuid.UUID) (*models.Budget, error) {
	var b models.Budget
	err := s.database.First(&b, "id = ? AND user_id = ?", budgetId, userId).Error
	if err != nil {
		return nil, err
	}
	return &b, err
}
