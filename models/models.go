package models

import (
	"github.com/AlsoShantanuBorkar/budget_max/models/auth"
	"github.com/AlsoShantanuBorkar/budget_max/models/budget"
	"github.com/AlsoShantanuBorkar/budget_max/models/categories"
	"github.com/AlsoShantanuBorkar/budget_max/models/transactions"
)

type (
	// Auth models
	User                   = auth.User
	Session                = auth.Session
	RefreshToken           = auth.RefreshToken
	PasswordResetToken     = auth.PasswordResetToken
	AuthRequest            = auth.AuthRequest
	RefreshTokensRequest   = auth.RefreshTokensRequest
	TwoFactorVerifyRequest = auth.TwoFactorVerifyRequest
	TwoFactorLoginRequest  = auth.TwoFactorLoginRequest
	TwoFAClaims            = auth.TwoFAClaims

	// Category models
	Category              = categories.Category
	CreateCategoryRequest = categories.CreateCategoryRequest
	UpdateCategoryRequest = categories.UpdateCategoryRequest

	// Transaction models
	Transaction               = transactions.Transaction
	CreateTransactionRequest  = transactions.CreateTransactionRequest
	UpdateTransactionRequest  = transactions.UpdateTransactionRequest
	TransactionFiltersRequest = transactions.TransactionFiltersRequest
	DateRangeRequest          = transactions.DateRangeRequest
	AmountRangeRequest        = transactions.AmountRangeRequest

	// Budget models
	Budget              = budget.Budget
	CreateBudgetRequest = budget.CreateBudgetRequest
	UpdateBudgetRequest = budget.UpdateBudgetRequest
)
