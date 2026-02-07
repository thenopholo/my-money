package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

type mockPlannedExpenseRepository struct {
	createFn    func(ctx context.Context, pe *domain.PlannedExpense) error
	getByIDFn   func(ctx context.Context, id uuid.UUID) (*domain.PlannedExpense, error)
	getByUserFn func(ctx context.Context, id uuid.UUID) ([]*domain.PlannedExpense, error)
	updateFn    func(ctx context.Context, pe *domain.PlannedExpense) error
	deleteFn    func(ctx context.Context, id uuid.UUID) error
}

func (m *mockPlannedExpenseRepository) Create(ctx context.Context, pe *domain.PlannedExpense) error {
	return m.createFn(ctx, pe)
}

func (m *mockPlannedExpenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PlannedExpense, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockPlannedExpenseRepository) GetByUserID(ctx context.Context, id uuid.UUID) ([]*domain.PlannedExpense, error) {
	return m.getByUserFn(ctx, id)
}

func (m *mockPlannedExpenseRepository) Update(ctx context.Context, pe *domain.PlannedExpense) error {
	return m.updateFn(ctx, pe)
}

func (m *mockPlannedExpenseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}
