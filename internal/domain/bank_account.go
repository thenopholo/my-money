package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountType string

const (
	AccountTypeChecking AccountType = "checking"
	AccountTypeSavings  AccountType = "savings"
)

func (a AccountType) IsValid() bool {
	return a == AccountTypeChecking || a == AccountTypeSavings
}

type BankAccount struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	AccountType AccountType
	BankName    string
	Balance     decimal.Decimal
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewBankAccount(userID uuid.UUID, accountType AccountType, name, bankName string, balance decimal.Decimal) (*BankAccount, error) {
	if name == "" {
		return nil, ErrEmptyBankName
	}

	if name == "" {
		name = bankName
	}

	if !accountType.IsValid() {
		return nil, ErrInvalidAccountType
	}

	if balance.LessThan(decimal.Zero) {
		return nil, ErrNegativeBalance
	}

	return &BankAccount{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		AccountType: accountType,
		BankName:    bankName,
		Balance:     balance,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}, nil
}
