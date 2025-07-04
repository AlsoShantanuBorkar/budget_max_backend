package services

import (
	"net/http"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCategory(c *gin.Context, req *models.CreateCategoryRequest, userId uuid.UUID) (*models.Category, *ServiceError) {
	category := &models.Category{
		ID:        uuid.New(),
		UserID:    userId,
		Name:      req.Name,
		Type:      req.Type,
		Icon:      req.Icon,
		IsDefault: req.IsDefault,
	}

	if err := database.CreateCategory(category); err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to create category")
	}

	return category, nil
}

func UpdateCategory(c *gin.Context, req *models.UpdateCategoryRequest, categoryId uuid.UUID, userId uuid.UUID) (*models.Category, *ServiceError) {
	// Fetch existing category to verify ownership
	_, err := database.GetCategoryByID(categoryId, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusNotFound, "category not found")
	}

	updates := make(map[string]any)
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Icon != nil {
		updates["icon"] = req.Icon
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
	}

	// Save updated category
	err = database.UpdateCategory(categoryId, updates)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to update category")
	}

	// Fetch updated category
	updatedCategory, err := database.GetCategoryByID(categoryId, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to fetch updated category")
	}

	return updatedCategory, nil
}

func DeleteCategory(c *gin.Context, categoryId uuid.UUID, userId uuid.UUID) *ServiceError {
	// Verify category exists and belongs to user
	_, err := database.GetCategoryByID(categoryId, userId)
	if err != nil {
		return NewServiceError(http.StatusNotFound, "category not found")
	}

	err = database.DeleteCategory(categoryId)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to delete category")
	}

	return nil
}

func GetCategoriesByUserID(c *gin.Context, userId uuid.UUID) ([]models.Category, *ServiceError) {
	categories, err := database.GetUserCategories(userId)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to fetch categories")
	}

	return categories, nil
}

func GetCategoryByID(c *gin.Context, categoryID uuid.UUID, userId uuid.UUID) (*models.Category, *ServiceError) {
	category, err := database.GetCategoryByID(categoryID, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusNotFound, "category not found")
	}

	return category, nil
}
