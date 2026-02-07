package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestTransactionType_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tt   TransactionType
		want bool
	}{
		{"income é válido", TransactionTypeIncome, true},
		{"expense é válido", TransactionTypeExpense, true},
		{"string vazia é inválido", TransactionType(""), false},
		{"tipo desconhecido é inválido", TransactionType("transfer"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.tt.IsValid(); got != tt.want {
				t.Errorf("TransactionType(%q).IsValid() = %v, want %v", tt.tt, got, tt.want)
			}
		})
	}
}

func TestNewTransaction_Validations(t *testing.T) {
	t.Parallel()

	accountID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	futureDate := now.AddDate(0, 0, 5)

	tests := []struct {
		name            string
		accountID       uuid.UUID
		categoryID      uuid.UUID
		amount          decimal.Decimal
		transactionType TransactionType
		description     string
		date            time.Time
		wantErr         error
	}{
		{
			name:            "deve criar transação income válida",
			accountID:       accountID,
			categoryID:      categoryID,
			amount:          decimal.NewFromInt(500),
			transactionType: TransactionTypeIncome,
			description:     "Salário",
			date:            yesterday,
			wantErr:         nil,
		},
		{
			name:            "deve criar transação expense válida",
			accountID:       accountID,
			categoryID:      categoryID,
			amount:          decimal.NewFromInt(100),
			transactionType: TransactionTypeExpense,
			description:     "Mercado",
			date:            yesterday,
			wantErr:         nil,
		},
		{
			name:            "deve usar descrição padrão quando vazia",
			accountID:       accountID,
			categoryID:      categoryID,
			amount:          decimal.NewFromInt(100),
			transactionType: TransactionTypeIncome,
			description:     "",
			date:            yesterday,
			wantErr:         nil,
		},
		{
			name:            "deve retornar erro para tipo inválido",
			accountID:       accountID,
			categoryID:      categoryID,
			amount:          decimal.NewFromInt(100),
			transactionType: TransactionType("invalid"),
			description:     "Teste",
			date:            yesterday,
			wantErr:         ErrInvalidTransactionType,
		},
		{
			name:            "deve retornar erro para valor zero",
			accountID:       accountID,
			categoryID:      categoryID,
			amount:          decimal.Zero,
			transactionType: TransactionTypeIncome,
			description:     "Teste",
			date:            yesterday,
			wantErr:         ErrInvalidAmount,
		},
		{
			name:            "deve retornar erro para valor negativo",
			accountID:       accountID,
			categoryID:      categoryID,
			amount:          decimal.NewFromInt(-100),
			transactionType: TransactionTypeIncome,
			description:     "Teste",
			date:            yesterday,
			wantErr:         ErrInvalidAmount,
		},
		{
			name:            "deve retornar erro para data futura",
			accountID:       accountID,
			categoryID:      categoryID,
			amount:          decimal.NewFromInt(100),
			transactionType: TransactionTypeIncome,
			description:     "Teste",
			date:            futureDate,
			wantErr:         ErrTransactionInFuture,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tx, err := NewTransaction(
				tt.accountID, tt.categoryID, nil, nil, nil,
				tt.amount, tt.transactionType, tt.description, tt.date,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if tx == nil {
					t.Fatal("NewTransaction() retornou nil sem erro")
				}
				if tt.description == "" && tx.Description != "Transacao" {
					t.Errorf("Description = %q, want 'Transacao'", tx.Description)
				}
				if !tx.Amount.Equal(tt.amount) {
					t.Errorf("Amount = %v, want %v", tx.Amount, tt.amount)
				}
			}
		})
	}
}
