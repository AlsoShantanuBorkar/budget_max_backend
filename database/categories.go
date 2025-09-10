package database

import (
	"errors"

	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(cat *models.Category) error
	GetCategoryByID(id uuid.UUID, userId uuid.UUID) (*models.Category, error)
	GetUserCategories(userID uuid.UUID) ([]models.Category, error)
	UpdateCategory(id uuid.UUID, updates map[string]interface{}) error
	DeleteCategory(id uuid.UUID) error
}

type CategoryDatabaseService struct {
	database *gorm.DB
}

func NewCategoryDatabaseService(db *gorm.DB) CategoryRepository {
	return &CategoryDatabaseService{database: db}
}


func (s *CategoryDatabaseService)  CreateCategory(cat *models.Category) error {
	return s.database.Create(cat).Error
}

func (s *CategoryDatabaseService)  GetCategoryByID(id uuid.UUID, userId uuid.UUID) (*models.Category, error) {
	var cat models.Category
	err := s.database.First(&cat, "id = ? AND user_id = ?", id, userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cat, err
}

func (s *CategoryDatabaseService)  GetUserCategories(userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := s.database.Where("user_id = ?", userID).Find(&categories).Error
	return categories, err
}

func (s *CategoryDatabaseService)  UpdateCategory(id uuid.UUID, updates map[string]interface{}) error {
	return s.database.Model(&models.Category{}).Where("id = ?", id).Updates(updates).Error
}

func (s *CategoryDatabaseService)  DeleteCategory(id uuid.UUID) error {
	return s.database.Delete(&models.Category{}, "id = ?", id).Error
}
