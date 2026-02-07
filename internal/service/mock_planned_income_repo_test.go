package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

type mockPlannedIncomeRepository struct {
	createFn    func(ctx context.Context, pi *domain.PlannedIncome) error
	getByIDFn   func(ctx context.Context, id uuid.UUID) (*domain.PlannedIncome, error)
	getByUserFn func(ctx context.Context, id uuid.UUID) ([]*domain.PlannedIncome, error)
	updateFn    func(ctx context.Context, pi *domain.PlannedIncome) error
	deleteFn    func(ctx context.Context, id uuid.UUID) error
}

func (m *mockPlannedIncomeRepository) Create(ctx context.Context, pi *domain.PlannedIncome) error {
	return m.createFn(ctx, pi)
}

func (m *mockPlannedIncomeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PlannedIncome, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockPlannedIncomeRepository) GetByUserID(ctx context.Context, id uuid.UUID) ([]*domain.PlannedIncome, error) {
	return m.getByUserFn(ctx, id)
}

func (m *mockPlannedIncomeRepository) Update(ctx context.Context, pi *domain.PlannedIncome) error {
	return m.updateFn(ctx, pi)
}

func (m *mockPlannedIncomeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}
