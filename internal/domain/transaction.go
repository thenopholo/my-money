package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID              uuid.UUID
	AccountID       uuid.UUID
	CategoryID      uuid.UUID
	InvoiceID       uuid.UUID
	PlannedIncome   uuid.UUID
	PlannedOutcome  uuid.UUID
	Amount          decimal.Decimal
	TransactionType string
	Description     string
	TransactionDate time.Time
	CreatedAt       time.Time
}

func NewTransaction(accountID, categoryID, invoiceID, plannedIncome, plannedOutcome uuid.UUID, amount decimal.Decimal, transactionType, description string) *Transaction {
	return &Transaction{
		ID:              uuid.New(),
		AccountID:       accountID,
		CategoryID:      categoryID,
		InvoiceID:       invoiceID,
		PlannedIncome:   plannedIncome,
		PlannedOutcome:  plannedOutcome,
		Amount:          amount,
		TransactionType: transactionType,
		Description:     description,
		CreatedAt:       time.Now(),
		TransactionDate: time.Now(),
	}
}
