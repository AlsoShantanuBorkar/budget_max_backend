package controllers

import (
	"net/http"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}
	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	category := models.Category{
		ID:        uuid.New(),
		UserID:    userId,
		Name:      req.Name,
		Type:      req.Type,
		Icon:      req.Icon,
		IsDefault: req.IsDefault,
	}

	err := database.CreateCategory(&category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data": gin.H{
			"category": category,
		},
	})

}

func GetAllCategories(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	categories, err := database.GetUserCategories(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error Occurred while fetching categories",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Categories fetched successfully",
		"data":    categories,
	})

}

func GetCategoryByID(c *gin.Context) {

	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Category ID format",
		})
		return
	}

	category, err := database.GetCategoryByID(categoryID, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error Fetching Category Data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category fetched successfully",
		"data":    category,
	})

}

func UpdateCategory(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	var req models.UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
		})
		return
	}

	categoryIdStr := c.Param("id")
	categoryId, err := uuid.Parse(categoryIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid category ID format",
		})
		return
	}

	// Fetch existing category
	category, err := database.GetCategoryByID(categoryId, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "category not found",
		})
		return
	}

	// Update only fields that are set
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Type != nil {
		category.Type = *req.Type
	}
	if req.Icon != nil {
		category.Icon = req.Icon
	}
	if req.IsDefault != nil {
		category.IsDefault = *req.IsDefault
	}

	// Save updated category
	err = database.UpdateCategory(categoryId, map[string]interface{}{
		"name":       category.Name,
		"type":       category.Type,
		"icon":       category.Icon,
		"is_default": category.IsDefault,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"data":    category,
	})
}

func DeleteCategory(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	categoryIdStr := c.Param("id")
	categoryId, err := uuid.Parse(categoryIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid category ID format",
		})
		return
	}

	_, err = database.GetCategoryByID(categoryId, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "category not found",
		})
		return
	}

	err = database.DeleteCategory(categoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error Occurred",
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Category deleted successfully",
	})
}
