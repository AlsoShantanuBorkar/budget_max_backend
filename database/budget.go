package database

import (
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
)

func CreateBudget(b *models.Budget) error {
	return DB.Create(b).Error
}

func GetBudgetsByUser(userID uuid.UUID) ([]models.Budget, error) {
	var budgets []models.Budget
	err := DB.Where("user_id = ?", userID).Find(&budgets).Error
	return budgets, err
}

func UpdateBudget(id uuid.UUID, updates map[string]any) error {
	return DB.Model(&models.Budget{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteBudget(id uuid.UUID) error {
	return DB.Delete(&models.Budget{}, "id = ?", id).Error
}

func GetBudgetByID(budgetId uuid.UUID, userId uuid.UUID) (*models.Budget, error) {
	var b models.Budget
	err := DB.First(&b, "id = ? AND user_id = ?", budgetId, userId).Error
	if err != nil {
		return nil, err
	}
	return &b, err
}
