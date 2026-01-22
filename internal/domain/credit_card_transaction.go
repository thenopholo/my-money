package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreditCardTransaction struct {
	ID              uuid.UUID
	CardID          uuid.UUID
	CategoryID      uuid.UUID
	InvoiceID       uuid.UUID
	Amount          decimal.Decimal
	Description     string
	Installments    int
	TransactionDate time.Time
	CreatedAt       time.Time
}

func NewCreditCardTransaction(cardID, categoryID, invoiceID uuid.UUID, amount decimal.Decimal, description string, installments int) *CreditCardTransaction {
	return &CreditCardTransaction{
		ID:              uuid.New(),
		CardID:          cardID,
		CategoryID:      categoryID,
		InvoiceID:       invoiceID,
		Amount:          amount,
		Description:     description,
		Installments:    installments,
		CreatedAt:       time.Now(),
		TransactionDate: time.Now(),
	}
}
