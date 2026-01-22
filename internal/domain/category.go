package domain

import "github.com/google/uuid"

type Category struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Name         string
	CategoryType string
	Color        string
	Icon         string
}

func NewCategory(userID uuid.UUID, categoryType, name, color, icon string) *Category {
	return &Category{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         name,
		CategoryType: categoryType,
		Color:        color,
		Icon:         icon,
	}
}
