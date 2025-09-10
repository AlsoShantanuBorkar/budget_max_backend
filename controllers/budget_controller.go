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

type BudgetControllerInterface interface {
	CreateBudget(c *gin.Context)
	UpdateBudget(c *gin.Context)
	DeleteBudget(c *gin.Context)
	GetBudgetsByUserID(c *gin.Context)
	GetBudgetByID(c *gin.Context)
}

type BudgetController struct {
	service services.BudgetServiceInterface
}

func NewBudgetController(service services.BudgetServiceInterface) *BudgetController {
	return &BudgetController{
		service: service,
	}
}

func (ctrl *BudgetController) CreateBudget(c *gin.Context) {
	var req models.CreateBudgetRequest

       if err := c.ShouldBindJSON(&req); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       budget, serviceErr := ctrl.service.CreateBudget(c, &req, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusCreated, gin.H{
	       "message": "Budget created successfully",
	       "data": gin.H{
		       "budget": budget,
	       },
       })
}

func (ctrl *BudgetController) UpdateBudget(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       var req models.UpdateBudgetRequest

       if err := c.ShouldBindJSON(&req); err != nil {
	       appErr := errors.NewBadRequestError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       if err := utils.GetValidator().Struct(req); err != nil {
	       appErr := errors.NewValidationError("Invalid Request", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       budgetIdStr := c.Param("id")
       budgetId, err := uuid.Parse(budgetIdStr)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid budget ID format", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       updatedBudget, serviceErr := ctrl.service.UpdateBudget(c, &req, budgetId, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Budget updated successfully",
	       "data":    updatedBudget,
       })
}

func (ctrl *BudgetController) DeleteBudget(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       budgetIdStr := c.Param("id")
       budgetId, err := uuid.Parse(budgetIdStr)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid budget ID format", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       serviceErr := ctrl.service.DeleteBudget(c, budgetId, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusNoContent, gin.H{
	       "message": "Budget deleted successfully",
       })
}

func (ctrl *BudgetController) GetBudgetsByUserID(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       budgets, serviceErr := ctrl.service.GetBudgetsByUserID(c, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Budgets fetched successfully",
	       "data":    budgets,
       })
}

func (ctrl *BudgetController) GetBudgetByID(c *gin.Context) {
       userId, ok := utils.ParseUserID(c)
       if !ok {
	       appErr := errors.NewUnauthorizedError("Invalid user ID", nil)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       budgetIDStr := c.Param("id")
       budgetID, err := uuid.Parse(budgetIDStr)
       if err != nil {
	       appErr := errors.NewBadRequestError("Invalid budget ID format", err)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       budget, serviceErr := ctrl.service.GetBudgetByID(c, budgetID, userId)
       if serviceErr != nil {
	       appErr := errors.NewInternalError(serviceErr)
	       c.JSON(appErr.Code, gin.H{"message": appErr.Message})
	       return
       }

       c.JSON(http.StatusOK, gin.H{
	       "message": "Budget fetched successfully",
	       "data":    budget,
       })
}
