package transactions

type TransactionFiltersRequest struct {
	BudgetID   *string  `json:"budget_id" validate:"omitempty,uuid4"`
	CategoryID *string  `json:"category_id" validate:"omitempty,uuid4"`
	Type       *string  `json:"type" validate:"omitempty,oneof=expense income"`
	StartDate  *string  `json:"start_date" validate:"omitempty,datetime"`
	EndDate    *string  `json:"end_date" validate:"omitempty,datetime"`
	MinAmount  *float64 `json:"min_amount" validate:"omitempty,gte=0"`
	MaxAmount  *float64 `json:"max_amount" validate:"omitempty,gte=0"`
}

type DateRangeRequest struct {
	StartDate string `json:"start_date" validate:"required,datetime"`
	EndDate   string `json:"end_date" validate:"required,datetime"`
}

type AmountRangeRequest struct {
	MinAmount float64 `json:"min_amount" validate:"required,gte=0"`
	MaxAmount float64 `json:"max_amount" validate:"required,gte=0"`
}
