package transactions

type UpdateTransactionRequest struct {
	Amount     *float64 `json:"amount" validate:"omitempty,gt=0"`
	Type       *string  `json:"type" validate:"omitempty,oneof=expense income"`
	Date       *string  `json:"date" validate:"omitempty,datetime"`
	Note       *string  `json:"note"`
	CategoryID *string  `json:"category_id" validate:"omitempty,uuid4"`
	BudgetID   *string  `json:"budget_id" validate:"omitempty,uuid4"`
	Name       *string  `json:"name" validate:"omitempty,min=1"`
}
