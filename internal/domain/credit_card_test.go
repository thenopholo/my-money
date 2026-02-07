package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestNewCreditCard_Validations(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	tests := []struct {
		name        string
		userID      uuid.UUID
		cardName    string
		creditLimit decimal.Decimal
		closeDay    int
		dueDay      int
		wantErr     error
	}{
		{
			name:        "deve criar cartão com dados válidos",
			userID:      userID,
			cardName:    "Nubank",
			creditLimit: decimal.NewFromInt(5000),
			closeDay:    15,
			dueDay:      25,
			wantErr:     nil,
		},
		{
			name:        "deve retornar erro para nome vazio",
			userID:      userID,
			cardName:    "",
			creditLimit: decimal.NewFromInt(5000),
			closeDay:    15,
			dueDay:      25,
			wantErr:     ErrEmptyCardName,
		},
		{
			name:        "deve retornar erro para limite zero",
			userID:      userID,
			cardName:    "Nubank",
			creditLimit: decimal.Zero,
			closeDay:    15,
			dueDay:      25,
			wantErr:     ErrInvalidLimit,
		},
		{
			name:        "deve retornar erro para limite negativo",
			userID:      userID,
			cardName:    "Nubank",
			creditLimit: decimal.NewFromInt(-1000),
			closeDay:    15,
			dueDay:      25,
			wantErr:     ErrInvalidLimit,
		},
		{
			name:        "deve retornar erro para close day 0",
			userID:      userID,
			cardName:    "Nubank",
			creditLimit: decimal.NewFromInt(5000),
			closeDay:    0,
			dueDay:      25,
			wantErr:     ErrInvalidCloseDay,
		},
		{
			name:        "deve retornar erro para close day 29",
			userID:      userID,
			cardName:    "Nubank",
			creditLimit: decimal.NewFromInt(5000),
			closeDay:    29,
			dueDay:      25,
			wantErr:     ErrInvalidCloseDay,
		},
		{
			name:        "deve aceitar close day 1",
			userID:      userID,
			cardName:    "Nubank",
			creditLimit: decimal.NewFromInt(5000),
			closeDay:    1,
			dueDay:      25,
			wantErr:     nil,
		},
		{
			name:        "deve aceitar close day 28",
			userID:      userID,
			cardName:    "Nubank",
			creditLimit: decimal.NewFromInt(5000),
			closeDay:    28,
			dueDay:      25,
			wantErr:     nil,
		},
		{
			name:        "deve retornar erro para due day 0",
			userID:      userID,
			cardName:    "Nubank",
			creditLimit: decimal.NewFromInt(5000),
			closeDay:    15,
			dueDay:      0,
			wantErr:     ErrInvalidDueDay,
		},
		{
			name:        "deve retornar erro para due day 32",
			userID:      userID,
			cardName:    "Nubank",
			creditLimit: decimal.NewFromInt(5000),
			closeDay:    15,
			dueDay:      32,
			wantErr:     ErrInvalidDueDay,
		},
		{
			name:        "deve aceitar due day 31",
			userID:      userID,
			cardName:    "Nubank",
			creditLimit: decimal.NewFromInt(5000),
			closeDay:    15,
			dueDay:      31,
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cc, err := NewCreditCard(tt.userID, tt.cardName, tt.creditLimit, tt.closeDay, tt.dueDay)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewCreditCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if cc == nil {
					t.Fatal("NewCreditCard() retornou nil sem erro")
				}
				if !cc.IsActive {
					t.Error("IsActive deve ser true para cartão novo")
				}
				if cc.Name != tt.cardName {
					t.Errorf("Name = %q, want %q", cc.Name, tt.cardName)
				}
				if !cc.CreditLimit.Equal(tt.creditLimit) {
					t.Errorf("CreditLimit = %v, want %v", cc.CreditLimit, tt.creditLimit)
				}
			}
		})
	}
}

func TestCreditCard_SetName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"deve aceitar nome válido", "Nubank", nil},
		{"deve retornar erro para nome vazio", "", ErrEmptyCardName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cc := &CreditCard{}
			err := cc.SetName(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SetName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreditCard_SetCreditLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		limit   decimal.Decimal
		wantErr error
	}{
		{"deve aceitar limite positivo", decimal.NewFromInt(5000), nil},
		{"deve rejeitar limite zero", decimal.Zero, ErrInvalidLimit},
		{"deve rejeitar limite negativo", decimal.NewFromInt(-1000), ErrInvalidLimit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cc := &CreditCard{}
			err := cc.SetCreditLimit(tt.limit)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SetCreditLimit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreditCard_SetCloseDay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		closeDay int
		wantErr  error
	}{
		{"deve aceitar dia 1", 1, nil},
		{"deve aceitar dia 28", 28, nil},
		{"deve aceitar dia 15", 15, nil},
		{"deve rejeitar dia 0", 0, ErrInvalidCloseDay},
		{"deve rejeitar dia 29", 29, ErrInvalidCloseDay},
		{"deve rejeitar dia negativo", -1, ErrInvalidCloseDay},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cc := &CreditCard{}
			err := cc.SetCloseDay(tt.closeDay)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SetCloseDay(%d) error = %v, wantErr %v", tt.closeDay, err, tt.wantErr)
			}
		})
	}
}

func TestCreditCard_SetDueDay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		dueDay int
		wantErr error
	}{
		{"deve aceitar dia 1", 1, nil},
		{"deve aceitar dia 31", 31, nil},
		{"deve rejeitar dia 0", 0, ErrInvalidDueDay},
		{"deve rejeitar dia 32", 32, ErrInvalidDueDay},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cc := &CreditCard{}
			err := cc.SetDueDay(tt.dueDay)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SetDueDay(%d) error = %v, wantErr %v", tt.dueDay, err, tt.wantErr)
			}
		})
	}
}

func TestCreditCard_SetIsActive(t *testing.T) {
	t.Parallel()

	cc := &CreditCard{IsActive: true}
	cc.SetIsActive(false)
	if cc.IsActive {
		t.Error("SetIsActive(false) não desativou o cartão")
	}
	cc.SetIsActive(true)
	if !cc.IsActive {
		t.Error("SetIsActive(true) não ativou o cartão")
	}
}
