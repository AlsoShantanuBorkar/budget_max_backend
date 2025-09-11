package controllers

import (
	"net/http"

	"github.com/AlsoShantanuBorkar/budget_max/errors"
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
	       appErr := errors.NewBadRequestError("Invalid Request", err, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }
       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       category, serviceErr := ctrl.service.CreateCategory(c, &req, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
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
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       categories, serviceErr := ctrl.service.GetCategoriesByUserID(c, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
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
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       categoryIDStr := c.Param("id")
       categoryID, err := uuid.Parse(categoryIDStr)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid Category ID format", err, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       category, serviceErr := ctrl.service.GetCategoryByID(c, categoryID, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
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
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       var req models.UpdateCategoryRequest

       if err := c.ShouldBindJSON(&req); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       categoryIdStr := c.Param("id")
       categoryId, err := uuid.Parse(categoryIdStr)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid category ID format", err, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       updatedCategory, serviceErr := ctrl.service.UpdateCategory(c, &req, categoryId, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
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
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       categoryIdStr := c.Param("id")
       categoryId, err := uuid.Parse(categoryIdStr)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid category ID format", err, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       serviceErr := ctrl.service.DeleteCategory(c, categoryId, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr, )
		   c.Error(appErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusNoContent, gin.H{
	       "message": "Category deleted successfully",
       })
}
