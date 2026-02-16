package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

type mockCreditCardTransactionRepository struct {
	createFn           func(ctx context.Context, t *domain.CreditCardTransaction) error
	getByIDFn          func(ctx context.Context, id uuid.UUID) (*domain.CreditCardTransaction, error)
	getByCardFn        func(ctx context.Context, cardID uuid.UUID) ([]*domain.CreditCardTransaction, error)
	getByCategoryFn    func(ctx context.Context, categoryID uuid.UUID) ([]*domain.CreditCardTransaction, error)
	getByInvoiceFn     func(ctx context.Context, invoiceID uuid.UUID) ([]*domain.CreditCardTransaction, error)
	getPendingFn       func(ctx context.Context, cardID uuid.UUID) ([]*domain.CreditCardTransaction, error)
	assignToInvFn      func(ctx context.Context, transactionID, invoiceID uuid.UUID) error
	updateFn           func(ctx context.Context, t *domain.CreditCardTransaction) error
	deleteFn           func(ctx context.Context, id uuid.UUID) error
	deleteAllByCardFn  func(ctx context.Context, cardID uuid.UUID) error
}

func (m *mockCreditCardTransactionRepository) Create(ctx context.Context, t *domain.CreditCardTransaction) error {
	return m.createFn(ctx, t)
}

func (m *mockCreditCardTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CreditCardTransaction, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockCreditCardTransactionRepository) GetByCardID(ctx context.Context, cardID uuid.UUID) ([]*domain.CreditCardTransaction, error) {
	return m.getByCardFn(ctx, cardID)
}

func (m *mockCreditCardTransactionRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*domain.CreditCardTransaction, error) {
	return m.getByCategoryFn(ctx, categoryID)
}

func (m *mockCreditCardTransactionRepository) GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*domain.CreditCardTransaction, error) {
	return m.getByInvoiceFn(ctx, invoiceID)
}

func (m *mockCreditCardTransactionRepository) GetPendingByCardID(ctx context.Context, cardID uuid.UUID) ([]*domain.CreditCardTransaction, error) {
	return m.getPendingFn(ctx, cardID)
}

func (m *mockCreditCardTransactionRepository) AssignToInvoice(ctx context.Context, transactionID, invoiceID uuid.UUID) error {
	return m.assignToInvFn(ctx, transactionID, invoiceID)
}

func (m *mockCreditCardTransactionRepository) Update(ctx context.Context, t *domain.CreditCardTransaction) error {
	return m.updateFn(ctx, t)
}

func (m *mockCreditCardTransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}

func (m *mockCreditCardTransactionRepository) DeleteAllByCardID(ctx context.Context, cardID uuid.UUID) error {
	return m.deleteAllByCardFn(ctx, cardID)
}
