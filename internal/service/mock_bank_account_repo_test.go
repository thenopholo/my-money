package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

type mockBankAccountRepository struct {
	createFn    func(ctx context.Context, ba *domain.BankAccount) error
	getByIDFn   func(ctx context.Context, id uuid.UUID) (*domain.BankAccount, error)
	getByUserFn func(ctx context.Context, id uuid.UUID) ([]*domain.BankAccount, error)
	updateFn    func(ctx context.Context, ba *domain.BankAccount) error
	deleteFn    func(ctx context.Context, id uuid.UUID) error
}

func (m *mockBankAccountRepository) Create(ctx context.Context, ba *domain.BankAccount) error {
	return m.createFn(ctx, ba)
}

func (m *mockBankAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.BankAccount, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockBankAccountRepository) GetByUserID(ctx context.Context, id uuid.UUID) ([]*domain.BankAccount, error) {
	return m.getByUserFn(ctx, id)
}

func (m *mockBankAccountRepository) Update(ctx context.Context, ba *domain.BankAccount) error {
	return m.updateFn(ctx, ba)
}

func (m *mockBankAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}
