package controllers

import (
	"net/http"

	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/services"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryControllerInterface interface {
	CreateCategory(c *gin.Context)
	GetAllCategories(c *gin.Context)
	GetCategoryByID(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

type CategoryController struct {
	service services.CategoryServiceInterface
}

func NewCategoryController(service services.CategoryServiceInterface) *CategoryController {
	return &CategoryController{
		service: service,
	}
}

func (ctrl *CategoryController) CreateCategory(c *gin.Context) {
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

	category, serviceErr := ctrl.service.CreateCategory(c, &req, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data": gin.H{
			"category": category,
		},
	})
}

func (ctrl *CategoryController) GetAllCategories(c *gin.Context) {
	userId, ok := utils.ParseUserID(c)
	if !ok {
		return
	}

	categories, serviceErr := ctrl.service.GetCategoriesByUserID(c, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categories fetched successfully",
		"data":    categories,
	})
}

func (ctrl *CategoryController) GetCategoryByID(c *gin.Context) {
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

	category, serviceErr := ctrl.service.GetCategoryByID(c, categoryID, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category fetched successfully",
		"data":    category,
	})
}

func (ctrl *CategoryController) UpdateCategory(c *gin.Context) {
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

	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
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

	updatedCategory, serviceErr := ctrl.service.UpdateCategory(c, &req, categoryId, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"data":    updatedCategory,
	})
}

func (ctrl *CategoryController) DeleteCategory(c *gin.Context) {
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

	serviceErr := ctrl.service.DeleteCategory(c, categoryId, userId)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, gin.H{"message": serviceErr.Message})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Category deleted successfully",
	})
}
