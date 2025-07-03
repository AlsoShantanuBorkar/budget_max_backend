package categories

type UpdateCategoryRequest struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,min=1"`
	Type      *string `json:"type,omitempty" validate:"omitempty,oneof=expense income"`
	Icon      *string `json:"icon,omitempty"`
	IsDefault *bool   `json:"is_default,omitempty"`
}
