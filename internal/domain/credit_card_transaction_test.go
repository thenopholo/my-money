package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestNewCreditCardTransaction_Validations(t *testing.T) {
	t.Parallel()

	cardID := uuid.New()
	categoryID := uuid.New()
	invoiceID := uuid.New()
	installValue := decimal.NewFromInt(50)

	tests := []struct {
		name              string
		cardID            uuid.UUID
		categoryID        uuid.UUID
		invoiceID         uuid.UUID
		amount            decimal.Decimal
		description       string
		installments      int
		currentInst       int
		installmentsValue *decimal.Decimal
		wantErr           error
	}{
		{
			name:              "deve criar transação válida sem parcelas",
			cardID:            cardID,
			categoryID:        categoryID,
			invoiceID:         invoiceID,
			amount:            decimal.NewFromInt(100),
			description:       "Compra",
			installments:      1,
			currentInst:       1,
			installmentsValue: nil,
			wantErr:           nil,
		},
		{
			name:              "deve criar transação válida com parcelas",
			cardID:            cardID,
			categoryID:        categoryID,
			invoiceID:         invoiceID,
			amount:            decimal.NewFromInt(500),
			description:       "Compra parcelada",
			installments:      10,
			currentInst:       1,
			installmentsValue: &installValue,
			wantErr:           nil,
		},
		{
			name:              "deve retornar erro para valor zero",
			cardID:            cardID,
			categoryID:        categoryID,
			invoiceID:         invoiceID,
			amount:            decimal.Zero,
			description:       "Compra",
			installments:      1,
			currentInst:       1,
			installmentsValue: nil,
			wantErr:           ErrInvalidAmount,
		},
		{
			name:              "deve retornar erro para valor negativo",
			cardID:            cardID,
			categoryID:        categoryID,
			invoiceID:         invoiceID,
			amount:            decimal.NewFromInt(-100),
			description:       "Compra",
			installments:      1,
			currentInst:       1,
			installmentsValue: nil,
			wantErr:           ErrInvalidAmount,
		},
		{
			name:              "deve retornar erro para descrição vazia",
			cardID:            cardID,
			categoryID:        categoryID,
			invoiceID:         invoiceID,
			amount:            decimal.NewFromInt(100),
			description:       "",
			installments:      1,
			currentInst:       1,
			installmentsValue: nil,
			wantErr:           ErrEmptyDescription,
		},
		{
			name:              "deve retornar erro para installments 0",
			cardID:            cardID,
			categoryID:        categoryID,
			invoiceID:         invoiceID,
			amount:            decimal.NewFromInt(100),
			description:       "Compra",
			installments:      0,
			currentInst:       1,
			installmentsValue: nil,
			wantErr:           ErrInvalidInstallments,
		},
		{
			name:              "deve retornar erro para installments 49",
			cardID:            cardID,
			categoryID:        categoryID,
			invoiceID:         invoiceID,
			amount:            decimal.NewFromInt(100),
			description:       "Compra",
			installments:      49,
			currentInst:       1,
			installmentsValue: nil,
			wantErr:           ErrInvalidInstallments,
		},
		{
			name:              "deve aceitar installments 48",
			cardID:            cardID,
			categoryID:        categoryID,
			invoiceID:         invoiceID,
			amount:            decimal.NewFromInt(4800),
			description:       "Compra grande",
			installments:      48,
			currentInst:       1,
			installmentsValue: &installValue,
			wantErr:           nil,
		},
		{
			name: "deve retornar erro para installmentsValue zero",
			cardID:            cardID,
			categoryID:        categoryID,
			invoiceID:         invoiceID,
			amount:            decimal.NewFromInt(100),
			description:       "Compra",
			installments:      2,
			currentInst:       1,
			installmentsValue: func() *decimal.Decimal { v := decimal.Zero; return &v }(),
			wantErr:           ErrInvalidAmount,
		},
		{
			name: "deve retornar erro para installmentsValue negativo",
			cardID:            cardID,
			categoryID:        categoryID,
			invoiceID:         invoiceID,
			amount:            decimal.NewFromInt(100),
			description:       "Compra",
			installments:      2,
			currentInst:       1,
			installmentsValue: func() *decimal.Decimal { v := decimal.NewFromInt(-10); return &v }(),
			wantErr:           ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ccTx, err := NewCreditCardTransaction(
				tt.cardID, tt.categoryID, tt.invoiceID,
				tt.amount, tt.description, tt.installments, tt.currentInst, tt.installmentsValue,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewCreditCardTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if ccTx == nil {
					t.Fatal("NewCreditCardTransaction() retornou nil sem erro")
				}
				if !ccTx.Amount.Equal(tt.amount) {
					t.Errorf("Amount = %v, want %v", ccTx.Amount, tt.amount)
				}
				if ccTx.Description != tt.description {
					t.Errorf("Description = %q, want %q", ccTx.Description, tt.description)
				}
				if ccTx.Installments != tt.installments {
					t.Errorf("Installments = %d, want %d", ccTx.Installments, tt.installments)
				}
			}
		})
	}
}
