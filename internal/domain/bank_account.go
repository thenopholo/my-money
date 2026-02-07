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
	if !accountType.IsValid() {
		return nil, ErrInvalidAccountType
	}

	ba := &BankAccount{
		ID:          uuid.New(),
		UserID:      userID,
		AccountType: accountType,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := ba.SetBankName(bankName); err != nil {
		return nil, err
	}

	ba.SetName(name)

	if err := ba.SetBalance(balance); err != nil {
		return nil, err
	}

	return ba, nil
}

func (ba *BankAccount) SetBankName(name string) error {
	if name == "" {
		return ErrEmptyBankName
	}

	ba.BankName = name
	return nil
}

func (ba *BankAccount) SetName(name string) {
	if name == "" {
		name = ba.BankName
	}

	ba.Name = name
}

func (ba *BankAccount) SetBalance(balance decimal.Decimal) error {
	if balance.LessThan(decimal.Zero) && !ba.AllowsOverdraft() {
		return ErrNegativeBalance
	}

	ba.Balance = balance
	return nil
}

func (ba *BankAccount) SetIsActive(isActive bool) {
	ba.IsActive = isActive
}

func (ba *BankAccount) AllowsOverdraft() bool {
	return ba.AccountType == AccountTypeChecking
}

func (ba *BankAccount) ApplyIncome(amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	return ba.SetBalance(ba.Balance.Add(amount))
}

func (ba *BankAccount) ApplyExpense(amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	return ba.SetBalance(ba.Balance.Sub(amount))
}
