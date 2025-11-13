package dto

type CreateCategoryDTO struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,max=255"`
}

type UpdateCategoryDTO struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}

type CategoryResponseDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}