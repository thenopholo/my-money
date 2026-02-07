package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

type mockInvoiceRepository struct {
	createFn    func(ctx context.Context, in *domain.Invoice) error
	getByIDFn   func(ctx context.Context, id uuid.UUID) (*domain.Invoice, error)
	getByCardFn func(ctx context.Context, cardID uuid.UUID) ([]*domain.Invoice, error)
	updateFn    func(ctx context.Context, in *domain.Invoice) error
	deleteFn    func(ctx context.Context, id uuid.UUID) error
}

func (m *mockInvoiceRepository) Create(ctx context.Context, in *domain.Invoice) error {
	return m.createFn(ctx, in)
}

func (m *mockInvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockInvoiceRepository) GetByCardID(ctx context.Context, cardID uuid.UUID) ([]*domain.Invoice, error) {
	return m.getByCardFn(ctx, cardID)
}

func (m *mockInvoiceRepository) Update(ctx context.Context, in *domain.Invoice) error {
	return m.updateFn(ctx, in)
}

func (m *mockInvoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}
