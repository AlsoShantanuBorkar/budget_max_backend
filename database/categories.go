package database

import (
	"errors"

	appErrors "github.com/AlsoShantanuBorkar/budget_max/errors"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryDatabaseServiceInterface interface {
	CreateCategory(cat *models.Category) error
	GetCategoryByID(id uuid.UUID, userId uuid.UUID) (*models.Category, error)
	GetUserCategories(userID uuid.UUID) ([]models.Category, error)
	UpdateCategory(id uuid.UUID, updates map[string]interface{}) error
	DeleteCategory(id uuid.UUID) error
}

type CategoryDatabaseService struct {
	database *gorm.DB
}

func (s *CategoryDatabaseService) CreateCategory(cat *models.Category) error {
	if err := s.database.Create(cat).Error; err != nil {
		return appErrors.NewDBError(err)
	}
	return nil
}

func NewCategoryDatabaseService(db *gorm.DB) CategoryDatabaseServiceInterface {
	return &CategoryDatabaseService{database: db}
}

func (s *CategoryDatabaseService) GetCategoryByID(id uuid.UUID, userId uuid.UUID) (*models.Category, error) {
	var cat models.Category
	err := s.database.First(&cat, "id = ? AND user_id = ?", id, userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, appErrors.NewDBError(err)
	}
	return &cat, nil
}

func (s *CategoryDatabaseService) GetUserCategories(userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := s.database.Where("user_id = ?", userID).Find(&categories).Error
	if err != nil {
		return nil, appErrors.NewDBError(err)
	}
	return categories, nil
}

func (s *CategoryDatabaseService) UpdateCategory(id uuid.UUID, updates map[string]interface{}) error {
	if err := s.database.Model(&models.Category{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return appErrors.NewDBError(err)
	}
	return nil
}

func (s *CategoryDatabaseService) DeleteCategory(id uuid.UUID) error {
	if err := s.database.Delete(&models.Category{}, "id = ?", id).Error; err != nil {
		return appErrors.NewDBError(err)
	}
	return nil
}
