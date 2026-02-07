package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestCategoryType_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ct   CategoryType
		want bool
	}{
		{"income é válido", CategoryTypeIncome, true},
		{"expense é válido", CategoryTypeExpense, true},
		{"string vazia é inválido", CategoryType(""), false},
		{"tipo desconhecido é inválido", CategoryType("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.ct.IsValid(); got != tt.want {
				t.Errorf("CategoryType(%q).IsValid() = %v, want %v", tt.ct, got, tt.want)
			}
		})
	}
}

func TestNewCategory_Validations(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	color := "blue"
	icon := "dollar"

	tests := []struct {
		name         string
		userID       uuid.UUID
		categoryType CategoryType
		catName      string
		color        *string
		icon         *string
		wantErr      error
	}{
		{
			name:         "deve criar categoria com dados válidos",
			userID:       userID,
			categoryType: CategoryTypeIncome,
			catName:      "Salário",
			color:        &color,
			icon:         &icon,
			wantErr:      nil,
		},
		{
			name:         "deve criar categoria sem cor e ícone",
			userID:       userID,
			categoryType: CategoryTypeExpense,
			catName:      "Alimentação",
			color:        nil,
			icon:         nil,
			wantErr:      nil,
		},
		{
			name:         "deve retornar erro para nome vazio",
			userID:       userID,
			categoryType: CategoryTypeIncome,
			catName:      "",
			color:        nil,
			icon:         nil,
			wantErr:      ErrEmptyCategoryName,
		},
		{
			name:         "deve retornar erro para tipo inválido",
			userID:       userID,
			categoryType: CategoryType("invalid"),
			catName:      "Teste",
			color:        nil,
			icon:         nil,
			wantErr:      ErrInvalidCategoryType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cat, err := NewCategory(tt.userID, tt.categoryType, tt.catName, tt.color, tt.icon)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if cat == nil {
					t.Fatal("NewCategory() retornou nil sem erro")
				}
				if cat.Name != tt.catName {
					t.Errorf("Name = %q, want %q", cat.Name, tt.catName)
				}
				if cat.CategoryType != tt.categoryType {
					t.Errorf("CategoryType = %q, want %q", cat.CategoryType, tt.categoryType)
				}
				if cat.UserID != tt.userID {
					t.Errorf("UserID = %v, want %v", cat.UserID, tt.userID)
				}
			}
		})
	}
}

func TestCategory_SetName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"deve aceitar nome válido", "Salário", nil},
		{"deve retornar erro para nome vazio", "", ErrEmptyCategoryName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Category{}
			err := c.SetName(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SetName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategory_SetType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   CategoryType
		wantErr error
	}{
		{"deve aceitar income", CategoryTypeIncome, nil},
		{"deve aceitar expense", CategoryTypeExpense, nil},
		{"deve rejeitar tipo inválido", CategoryType("unknown"), ErrInvalidCategoryType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Category{}
			err := c.SetType(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SetType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategory_SetVisuals(t *testing.T) {
	t.Parallel()

	color := "red"
	icon := "star"

	c := &Category{}
	c.SetVisuals(&color, &icon)

	if c.Color == nil || *c.Color != color {
		t.Errorf("Color = %v, want %q", c.Color, color)
	}
	if c.Icon == nil || *c.Icon != icon {
		t.Errorf("Icon = %v, want %q", c.Icon, icon)
	}

	c.SetVisuals(nil, nil)
	if c.Color != nil {
		t.Error("Color deveria ser nil")
	}
	if c.Icon != nil {
		t.Error("Icon deveria ser nil")
	}
}
