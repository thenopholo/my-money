package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionTypeIncome TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

func (t TransactionType) IsValid() bool {
	return t == TransactionTypeIncome || t == TransactionTypeExpense
}

type Transaction struct {
	ID              uuid.UUID
	AccountID       uuid.UUID
	CategoryID      uuid.UUID
	InvoiceID       *uuid.UUID
	PlannedIncome   *uuid.UUID
	PlannedExpense  *uuid.UUID
	Amount          decimal.Decimal
	TransactionType TransactionType
	Description     string
	TransactionDate time.Time
	CreatedAt       time.Time
}

func NewTransaction(accountID, categoryID uuid.UUID, invoiceID, plannedIncome, plannedExpense *uuid.UUID, amount decimal.Decimal, transactionType TransactionType, description string, transactionDate time.Time) (*Transaction, error) {
	if !transactionType.IsValid() {
		return nil, ErrInvalidTransactionType
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidAmount
	}

	if transactionDate.After(time.Now().AddDate(0, 0, 1)) {
		return nil, ErrTransactionInFuture
	}

	if description == "" {
		description = "Transação"
	}

	return &Transaction{
		ID:              uuid.New(),
		AccountID:       accountID,
		CategoryID:      categoryID,
		InvoiceID:       invoiceID,
		PlannedIncome:   plannedIncome,
		PlannedExpense:  plannedExpense,
		Amount:          amount,
		TransactionType: transactionType,
		Description:     description,
		TransactionDate: transactionDate,
		CreatedAt:       time.Now(),
	}, nil
}
