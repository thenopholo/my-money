package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

type mockCreditCardRepository struct {
	createFn    func(ctx context.Context, cc *domain.CreditCard) error
	getByIDFn   func(ctx context.Context, id uuid.UUID) (*domain.CreditCard, error)
	getByUserFn func(ctx context.Context, id uuid.UUID) ([]*domain.CreditCard, error)
	updateFn    func(ctx context.Context, cc *domain.CreditCard) error
	deleteFn    func(ctx context.Context, id uuid.UUID) error
}

func (m *mockCreditCardRepository) Create(ctx context.Context, cc *domain.CreditCard) error {
	return m.createFn(ctx, cc)
}

func (m *mockCreditCardRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CreditCard, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockCreditCardRepository) GetByUserID(ctx context.Context, id uuid.UUID) ([]*domain.CreditCard, error) {
	return m.getByUserFn(ctx, id)
}

func (m *mockCreditCardRepository) Update(ctx context.Context, cc *domain.CreditCard) error {
	return m.updateFn(ctx, cc)
}

func (m *mockCreditCardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}
