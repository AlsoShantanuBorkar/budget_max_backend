package services

import (
	"net/http"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateBudget(c *gin.Context, req *models.CreateBudgetRequest, userId uuid.UUID) (*models.Budget, *ServiceError) {
	budget := &models.Budget{
		ID:        uuid.New(),
		UserID:    userId,
		Type:      req.Type,
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}

	if err := database.CreateBudget(budget); err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to create budget")
	}

	return budget, nil
}

func UpdateBudget(c *gin.Context, req *models.UpdateBudgetRequest, budgetId uuid.UUID, userId uuid.UUID) (*models.Budget, *ServiceError) {
	// Fetch existing budget to verify ownership
	_, err := database.GetBudgetByID(budgetId, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusNotFound, "budget not found")
	}

	updates := make(map[string]any)
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = *req.EndDate
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}

	// Save updated budget
	err = database.UpdateBudget(budgetId, updates)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to update budget")
	}

	// Fetch updated budget
	updatedBudget, err := database.GetBudgetByID(budgetId, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to fetch updated budget")
	}

	return updatedBudget, nil
}

func DeleteBudget(c *gin.Context, budgetId uuid.UUID, userId uuid.UUID) *ServiceError {
	// Verify budget exists and belongs to user
	_, err := database.GetBudgetByID(budgetId, userId)
	if err != nil {
		return NewServiceError(http.StatusNotFound, "budget not found")
	}

	err = database.DeleteBudget(budgetId)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, "failed to delete budget")
	}

	return nil
}

func GetBudgetsByUserID(c *gin.Context, userId uuid.UUID) ([]models.Budget, *ServiceError) {
	budgets, err := database.GetBudgetsByUser(userId)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, "failed to fetch budgets")
	}

	return budgets, nil
}

func GetBudgetByID(c *gin.Context, budgetID uuid.UUID, userId uuid.UUID) (*models.Budget, *ServiceError) {
	budget, err := database.GetBudgetByID(budgetID, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusNotFound, "budget not found")
	}

	return budget, nil
}
