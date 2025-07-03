package database

import (
	"errors"

	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateCategory(cat *models.Category) error {
	return DB.Create(cat).Error
}

func GetCategoryByID(id uuid.UUID, userId uuid.UUID) (*models.Category, error) {
	var cat models.Category
	err := DB.First(&cat, "id = ? AND user_id = ?", id, userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cat, err
}

func GetUserCategories(userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := DB.Where("user_id = ?", userID).Find(&categories).Error
	return categories, err
}

func UpdateCategory(id uuid.UUID, updates map[string]interface{}) error {
	return DB.Model(&models.Category{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteCategory(id uuid.UUID) error {
	return DB.Delete(&models.Category{}, "id = ?", id).Error
}
