package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

type mockTransactionRepository struct {
	createFn             func(ctx context.Context, t *domain.Transaction) error
	getByIDFn            func(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	getByAccountFn       func(ctx context.Context, accountID uuid.UUID) ([]*domain.Transaction, error)
	getByCategoryFn      func(ctx context.Context, categoryID uuid.UUID) ([]*domain.Transaction, error)
	getByInvoiceFn       func(ctx context.Context, invoiceID uuid.UUID) ([]*domain.Transaction, error)
	getByPlannedIncomeFn func(ctx context.Context, id uuid.UUID) ([]*domain.Transaction, error)
	getByPlannedExpenseFn func(ctx context.Context, id uuid.UUID) ([]*domain.Transaction, error)
	updateFn             func(ctx context.Context, t *domain.Transaction) error
	deleteFn             func(ctx context.Context, id uuid.UUID) error
}

func (m *mockTransactionRepository) Create(ctx context.Context, t *domain.Transaction) error {
	return m.createFn(ctx, t)
}

func (m *mockTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockTransactionRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*domain.Transaction, error) {
	return m.getByAccountFn(ctx, accountID)
}

func (m *mockTransactionRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*domain.Transaction, error) {
	return m.getByCategoryFn(ctx, categoryID)
}

func (m *mockTransactionRepository) GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*domain.Transaction, error) {
	return m.getByInvoiceFn(ctx, invoiceID)
}

func (m *mockTransactionRepository) GetByPlannedIncomeID(ctx context.Context, id uuid.UUID) ([]*domain.Transaction, error) {
	return m.getByPlannedIncomeFn(ctx, id)
}

func (m *mockTransactionRepository) GetByPlannedExpenseID(ctx context.Context, id uuid.UUID) ([]*domain.Transaction, error) {
	return m.getByPlannedExpenseFn(ctx, id)
}

func (m *mockTransactionRepository) Update(ctx context.Context, t *domain.Transaction) error {
	return m.updateFn(ctx, t)
}

func (m *mockTransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}
