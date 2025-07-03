package transactions

import (
	"github.com/google/uuid"
)

type CreateTransactionRequest struct {
	Amount      float64    `json:"amount" validate:"required,gt=0"`
	Type        string     `json:"type" validate:"required,oneof=expense income"`
	Name        string     `json:"name" validate:"required,min=1"`
	Date        string     `json:"date" validate:"required,datetime"`
	Note        string     `json:"note" validate:"omitempty"`
	CategoryIDs *uuid.UUID `json:"category_ids" validate:"omitempty,uuid4"` // optional
}
