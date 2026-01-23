package domain

import "github.com/google/uuid"

type CategoryType string

const (
	CategoryTypeIncome CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

func (c CategoryType) IsValid() bool {
	return c == CategoryTypeIncome || c == CategoryTypeExpense
}

type Category struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Name         string
	CategoryType CategoryType
	Color        *string
	Icon         *string
}

func NewCategory(userID uuid.UUID, categoryType CategoryType, name string, color, icon *string) (*Category, error) {
	if name == "" {
		return nil, ErrEmptyCategoryName
	}

	if !categoryType.IsValid() {
		return nil, ErrInvalidCategoryType
	}

	return &Category{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         name,
		CategoryType: categoryType,
		Color:        color,
		Icon:         icon,
	}, nil
}
