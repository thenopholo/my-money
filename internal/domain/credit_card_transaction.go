package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreditCardTransaction struct {
	ID                  uuid.UUID
	CardID              uuid.UUID
	CategoryID          uuid.UUID
	InvoiceID           uuid.UUID
	Amount              decimal.Decimal
	Description         string
	Installments        int
	CurrentInstallments int
	InstallmentsValue   *decimal.Decimal
	TransactionDate     time.Time
	CreatedAt           time.Time
}

func NewCreditCardTransaction(cardID, categoryID, invoiceID uuid.UUID, amount decimal.Decimal, description string, installments, currentInstallments int, installmentsValue *decimal.Decimal) (*CreditCardTransaction, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidAmount
	}

	if description == "" {
		return nil, ErrEmptyDescription
	}

	if installmentsValue != nil && installmentsValue.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidAmount
	}

	if installments < 1 || installments > 48 {
		return nil, ErrInvalidInstallments
	}

	return &CreditCardTransaction{
		ID:                  uuid.New(),
		CardID:              cardID,
		CategoryID:          categoryID,
		InvoiceID:           invoiceID,
		Amount:              amount,
		Description:         description,
		Installments:        installments,
		CurrentInstallments: currentInstallments,
		InstallmentsValue:   installmentsValue,
		CreatedAt:           time.Now(),
		TransactionDate:     time.Now(),
	}, nil
}
