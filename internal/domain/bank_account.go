package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BankAccount struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	AccountType string
	BankName    string
	Balance     decimal.Decimal
	IsActive    bool
	CreatedAt   time.Time
}

func NewBankAccount( userID uuid.UUID, accountType, bankName string, balance decimal.Decimal, isActive bool) *BankAccount {
	return &BankAccount{
		ID:          uuid.New(),
		UserID:      userID,
		AccountType: accountType,
		BankName:    bankName,
		Balance:     balance,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
	}
}
