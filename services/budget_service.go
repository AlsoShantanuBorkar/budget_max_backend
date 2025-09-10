package services

import (
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/errors"
	"github.com/AlsoShantanuBorkar/budget_max/models"
	"github.com/AlsoShantanuBorkar/budget_max/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BudgetServiceInterface interface {
	CreateBudget(c *gin.Context, req *models.CreateBudgetRequest, userId uuid.UUID) (*models.Budget, *ServiceError)
	UpdateBudget(c *gin.Context, req *models.UpdateBudgetRequest, budgetId uuid.UUID, userId uuid.UUID) (*models.Budget, *ServiceError)
	DeleteBudget(c *gin.Context, budgetId uuid.UUID, userId uuid.UUID) *ServiceError
	GetBudgetsByUserID(c *gin.Context, userId uuid.UUID) ([]models.Budget, *ServiceError)
	GetBudgetByID(c *gin.Context, budgetID uuid.UUID, userId uuid.UUID) (*models.Budget, *ServiceError)
}

type BudgetService struct {
	databaseService database.BudgetDatabaseServiceInterface
}

func NewBudgetService(dbService database.BudgetDatabaseServiceInterface) BudgetServiceInterface {
	return &BudgetService{databaseService: dbService}
}

func (s *BudgetService) CreateBudget(c *gin.Context, req *models.CreateBudgetRequest, userId uuid.UUID) (*models.Budget, *ServiceError) {
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

       if err := s.databaseService.CreateBudget(budget); err != nil {
	       utils.GetLogger().Error().Err(err).Msg("failed to create budget")
	       return nil, ServiceErrorFromAppError(errors.NewDBError(err))
       }

       return budget, nil
}

func (s *BudgetService) UpdateBudget(c *gin.Context, req *models.UpdateBudgetRequest, budgetId uuid.UUID, userId uuid.UUID) (*models.Budget, *ServiceError) {
       // Fetch existing budget to verify ownership
       _, err := s.databaseService.GetBudgetByID(budgetId, userId)
       if err != nil {
	       utils.GetLogger().Error().Err(err).Msg("budget not found for update")
	       return nil, ServiceErrorFromAppError(errors.NewNotFoundError("budget", err))
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
       err = s.databaseService.UpdateBudget(budgetId, updates)
       if err != nil {
	       utils.GetLogger().Error().Err(err).Msg("failed to update budget")
	       return nil, ServiceErrorFromAppError(errors.NewDBError(err))
       }

       // Fetch updated budget
       updatedBudget, err := s.databaseService.GetBudgetByID(budgetId, userId)
       if err != nil {
	       utils.GetLogger().Error().Err(err).Msg("failed to fetch updated budget")
	       return nil, ServiceErrorFromAppError(errors.NewDBError(err))
       }

       return updatedBudget, nil
}

func (s *BudgetService) DeleteBudget(c *gin.Context, budgetId uuid.UUID, userId uuid.UUID) *ServiceError {
       // Verify budget exists and belongs to user
       _, err := s.databaseService.GetBudgetByID(budgetId, userId)
       if err != nil {
	       utils.GetLogger().Error().Err(err).Msg("budget not found for delete")
	       return ServiceErrorFromAppError(errors.NewNotFoundError("budget", err))
       }

       err = s.databaseService.DeleteBudget(budgetId)
       if err != nil {
	       utils.GetLogger().Error().Err(err).Msg("failed to delete budget")
	       return ServiceErrorFromAppError(errors.NewDBError(err))
       }

       return nil
}

func (s *BudgetService) GetBudgetsByUserID(c *gin.Context, userId uuid.UUID) ([]models.Budget, *ServiceError) {
       budgets, err := s.databaseService.GetBudgetsByUser(userId)
       if err != nil {
	       utils.GetLogger().Error().Err(err).Msg("failed to fetch budgets")
	       return nil, ServiceErrorFromAppError(errors.NewDBError(err))
       }

       return budgets, nil
}

func (s *BudgetService) GetBudgetByID(c *gin.Context, budgetID uuid.UUID, userId uuid.UUID) (*models.Budget, *ServiceError) {
       budget, err := s.databaseService.GetBudgetByID(budgetID, userId)
       if err != nil {
	       utils.GetLogger().Error().Err(err).Msg("budget not found by id")
	       return nil, ServiceErrorFromAppError(errors.NewNotFoundError("budget", err))
       }

       return budget, nil
}
