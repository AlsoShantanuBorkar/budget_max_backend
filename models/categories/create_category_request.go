package categories

type CreateCategoryRequest struct {
	Name      string  `json:"name" validate:"required,min=1"`
	Type      string  `json:"type" validate:"required,oneof=expense income"`
	Icon      *string `json:"icon,omitempty"`
	IsDefault bool    `json:"is_default"`
}
