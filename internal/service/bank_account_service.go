package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

type BankAccountService struct {
	bankAccountRepo BankAccountRepository
}

func NewBankAccountService(bankAccountRepo BankAccountRepository) *BankAccountService {
	return &BankAccountService{bankAccountRepo: bankAccountRepo}
}

func (bas *BankAccountService) Create(ctx context.Context, userID uuid.UUID, accountType domain.AccountType, name, bankName string, balance decimal.Decimal) (*domain.BankAccount, error) {
	newAccount, err := domain.NewBankAccount(userID, accountType, name, bankName, balance)
	if err != nil {
		return nil, err
	}

	newBankAccount, err := bas.bankAccountRepo.Create(ctx, newAccount)
	if err != nil {
		return nil, err
	}

	return newBankAccount, nil
}

func (bas *BankAccountService) GetByID(ctx context.Context, id uuid.UUID) (*domain.BankAccount, error) {
	account, err := bas.bankAccountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (bas *BankAccountService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.BankAccount, error) {
	accounts, err := bas.bankAccountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (bas *BankAccountService) Update(ctx context.Context, id uuid.UUID, name string, bankName string, balance decimal.Decimal, isActive bool) (*domain.BankAccount, error) {
	account, err := bas.bankAccountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if bankName != "" {
		return nil, domain.ErrEmptyBankName
	}
	if balance.LessThan(decimal.Zero) {
		return nil, domain.ErrNegativeBalance
	}
	if name == "" {
		name = bankName
	}

	account.Name = name
	account.BankName = bankName
	account.Balance = balance
	account.IsActive = isActive

	if err := bas.bankAccountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (bas *BankAccountService) Delete(ctx context.Context, id uuid.UUID) error {
	return bas.bankAccountRepo.Delete(ctx, id)
}
