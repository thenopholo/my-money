package domain

import "github.com/google/uuid"

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
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
	c := &Category{
		ID:     uuid.New(),
		UserID: userID,
	}

	if err := c.SetName(name); err != nil {
		return nil, err
	}

	if err := c.SetType(categoryType); err != nil {
		return nil, err
	}

	c.SetVisuals(color, icon)

	return c, nil
}

func (c *Category) SetName(name string) error {
	if name == "" {
		return ErrEmptyCategoryName
	}

	c.Name = name
	return nil
}

func (c *Category) SetType(categoryType CategoryType) error {
	if !categoryType.IsValid() {
		return ErrInvalidCategoryType
	}

	c.CategoryType = categoryType
	return nil
}

func (c *Category) SetVisuals(color, icon *string) {
	c.Color = color
	c.Icon = icon
}
