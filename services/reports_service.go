package services

import (
	"sort"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/database"
	"github.com/AlsoShantanuBorkar/budget_max/errors"
	"github.com/AlsoShantanuBorkar/budget_max/models/reports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReportsServiceInterface interface {
	GetBudgetSummary(c *gin.Context, budgetID uuid.UUID, userId uuid.UUID) (*reports.BudgetSummary, *ServiceError)
	GetWeeklySummary(c *gin.Context, userId uuid.UUID, startDate time.Time, endDate time.Time) (*reports.WeeklySummary, *ServiceError)
	GetMonthlySummary(c *gin.Context, userId uuid.UUID, month time.Time) (*reports.MonthlySummary, *ServiceError)
	GetYearlySummary(c *gin.Context, userId uuid.UUID, year time.Time) (*reports.YearlySummary, *ServiceError)
	GetCategorySummary(c *gin.Context, userId uuid.UUID, categoryID uuid.UUID) (*reports.CategorySummary, *ServiceError)
	GetCustomDateRangeSummary(c *gin.Context, userId uuid.UUID, startDate time.Time, endDate time.Time) (*reports.CustomDateRangeSummary, *ServiceError)
	GetDailyAverageSummary(c *gin.Context, userId uuid.UUID, startDate time.Time, endDate time.Time) (*reports.DailyAverageSummary, *ServiceError)
	GetTopCategories(c *gin.Context, userId uuid.UUID, limit int, transactionType string) ([]*reports.TopCategory, *ServiceError)
	GetAllCategoriesSummary(c *gin.Context, userId uuid.UUID) ([]*reports.CategorySummary, *ServiceError)
}

type ReportsService struct {
	transactionDatabaseService database.TransactionDatabaseServiceInterface
	categoryDatabaseService    database.CategoryDatabaseServiceInterface
	budgetDatabaseService      database.BudgetDatabaseServiceInterface
}

func NewReportsService(txnDBService database.TransactionDatabaseServiceInterface, catDBService database.CategoryDatabaseServiceInterface, budgetDBService database.BudgetDatabaseServiceInterface) ReportsServiceInterface {
	return &ReportsService{
		transactionDatabaseService: txnDBService,
		categoryDatabaseService:    catDBService,
		budgetDatabaseService:      budgetDBService,
	}
}

func (s *ReportsService) GetBudgetSummary(c *gin.Context, budgetID uuid.UUID, userId uuid.UUID) (*reports.BudgetSummary, *ServiceError) {
	budget, err := s.budgetDatabaseService.GetBudgetByID(budgetID, userId)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }
	txns, err := s.transactionDatabaseService.GetTransactionsByBudget(userId, budgetID)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.BudgetSummary{
		BudgetID:      budgetID,
		BudgetName:    budget.Name,
		BudgetAmount:  budget.Amount,
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func (s *ReportsService) GetWeeklySummary(c *gin.Context, userId uuid.UUID, startDate time.Time, endDate time.Time) (*reports.WeeklySummary, *ServiceError) {
	txns, err := s.transactionDatabaseService.GetTransactionsByDateRange(userId, startDate, endDate)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.WeeklySummary{
		StartDate:     startDate.Format("2006-01-02"),
		EndDate:       endDate.Format("2006-01-02"),
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func (s *ReportsService) GetMonthlySummary(c *gin.Context, userId uuid.UUID, month time.Time) (*reports.MonthlySummary, *ServiceError) {
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	txns, err := s.transactionDatabaseService.GetTransactionsByDateRange(userId, startOfMonth, endOfMonth)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.MonthlySummary{
		Month:         month.Format("January 2006"),
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func (s *ReportsService) GetYearlySummary(c *gin.Context, userId uuid.UUID, year time.Time) (*reports.YearlySummary, *ServiceError) {
	startOfYear := time.Date(year.Year(), 1, 1, 0, 0, 0, 0, year.Location())
	endOfYear := time.Date(year.Year(), 12, 31, 23, 59, 59, 999999999, year.Location())

	txns, err := s.transactionDatabaseService.GetTransactionsByDateRange(userId, startOfYear, endOfYear)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.YearlySummary{
		Year:          year.Format("2006"),
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func (s *ReportsService) GetCategorySummary(c *gin.Context, userId uuid.UUID, categoryID uuid.UUID) (*reports.CategorySummary, *ServiceError) {
	category, err := s.categoryDatabaseService.GetCategoryByID(categoryID, userId)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }
       if category == nil {
	       appErr := errors.NewNotFoundError("category", nil, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	txns, err := s.transactionDatabaseService.GetTransactionsByCategory(userId, categoryID)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.CategorySummary{
		CategoryID:    categoryID,
		CategoryName:  category.Name,
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func (s *ReportsService) GetCustomDateRangeSummary(c *gin.Context, userId uuid.UUID, startDate time.Time, endDate time.Time) (*reports.CustomDateRangeSummary, *ServiceError) {
	txns, err := s.transactionDatabaseService.GetTransactionsByDateRange(userId, startDate, endDate)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }
	totalExpenses := 0.0
	totalIncome := 0.0
	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}
	netBalance := totalIncome - totalExpenses
	return &reports.CustomDateRangeSummary{
		StartDate:     startDate.Format("2006-01-02"),
		EndDate:       endDate.Format("2006-01-02"),
		TotalExpenses: totalExpenses,
		TotalIncome:   totalIncome,
		NetBalance:    netBalance,
	}, nil
}

func (s *ReportsService) GetDailyAverageSummary(c *gin.Context, userId uuid.UUID, startDate time.Time, endDate time.Time) (*reports.DailyAverageSummary, *ServiceError) {
	txns, err := s.transactionDatabaseService.GetTransactionsByDateRange(userId, startDate, endDate)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	totalExpenses := 0.0
	totalIncome := 0.0
	daysDiff := int(endDate.Sub(startDate).Hours()/24) + 1

	for _, txn := range txns {
		if txn.Type == "expense" {
			totalExpenses += txn.Amount
		} else {
			totalIncome += txn.Amount
		}
	}

	dailyExpenses := totalExpenses / float64(daysDiff)
	dailyIncome := totalIncome / float64(daysDiff)
	dailyNetBalance := dailyIncome - dailyExpenses

	return &reports.DailyAverageSummary{
		Date:          startDate.Format("2006-01-02") + " to " + endDate.Format("2006-01-02"),
		TotalExpenses: dailyExpenses,
		TotalIncome:   dailyIncome,
		NetBalance:    dailyNetBalance,
	}, nil
}

func (s *ReportsService) GetTopCategories(c *gin.Context, userId uuid.UUID, limit int, transactionType string) ([]*reports.TopCategory, *ServiceError) {
	if limit <= 0 {
		limit = 5
	}

	// Get all transactions for the user
	txns, err := s.transactionDatabaseService.GetTransactionsByType(userId, transactionType)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	// Group transactions by category
	categoryTotals := make(map[string]float64)
	categoryNames := make(map[string]string)

	for _, txn := range txns {
		categoryTotals[txn.CategoryID.String()] += txn.Amount
		// We'll get category names in the next step
	}

	// Get category names
	for categoryID := range categoryTotals {
		parsedId, err := uuid.Parse(categoryID)
		if err != nil {
			continue
		}
		category, err := s.categoryDatabaseService.GetCategoryByID(parsedId, userId)
		if err != nil {
			continue
		}
		if category != nil {
			categoryNames[categoryID] = category.Name
		}
	}

	// Convert to slice for sorting
	var categories []*reports.TopCategory
	for categoryID, amount := range categoryTotals {
		categoryName, exists := categoryNames[categoryID]
		if !exists {
			categoryName = "Unknown Category"
		}
		categories = append(categories, &reports.TopCategory{
			CategoryID:   categoryID,
			CategoryName: categoryName,
			Amount:       amount,
			Type:         transactionType,
		})
	}

	// Sort by amount (descending for expenses, ascending for income)
	if transactionType == "expense" {
		sort.Slice(categories, func(i, j int) bool {
			return categories[i].Amount > categories[j].Amount
		})
	} else {
		sort.Slice(categories, func(i, j int) bool {
			return categories[i].Amount < categories[j].Amount
		})
	}

	// Limit results and add rank
	if len(categories) > limit {
		categories = categories[:limit]
	}

	for i, category := range categories {
		category.Rank = i + 1
	}

	return categories, nil
}

func (s *ReportsService) GetAllCategoriesSummary(c *gin.Context, userId uuid.UUID) ([]*reports.CategorySummary, *ServiceError) {
	categories, err := s.categoryDatabaseService.GetUserCategories(userId)
       if err != nil {
	       appErr := errors.NewInternalError(err, )
	       c.Error(appErr)
	       return nil, ServiceErrorFromAppError(appErr)
       }

	var summaries []*reports.CategorySummary

	for _, category := range categories {
		txns, err := s.transactionDatabaseService.GetTransactionsByCategory(userId, category.ID)
		if err != nil {
			continue
		}

		totalExpenses := 0.0
		totalIncome := 0.0
		for _, txn := range txns {
			if txn.Type == "expense" {
				totalExpenses += txn.Amount
			} else {
				totalIncome += txn.Amount
			}
		}

		netBalance := totalIncome - totalExpenses
		summaries = append(summaries, &reports.CategorySummary{
			CategoryID:    category.ID,
			CategoryName:  category.Name,
			TotalExpenses: totalExpenses,
			TotalIncome:   totalIncome,
			NetBalance:    netBalance,
		})
	}

	return summaries, nil
}
