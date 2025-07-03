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
	Transaction              = transactions.Transaction
	CreateTransactionRequest = transactions.CreateTransactionRequest
	UpdateTransactionRequest = transactions.UpdateTransactionRequest

	// Budget models
	Budget              = budget.Budget
	CreateBudgetRequest = budget.CreateBudgetRequest
	UpdateBudgetRequest = budget.UpdateBudgetRequest
)

//  		"refresh": "df9fbc98-3acf-4469-9c76-f59b4f292b9f",
