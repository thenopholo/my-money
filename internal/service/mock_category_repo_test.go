package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

type mockCategoryRepository struct {
	createFn    func(ctx context.Context, c *domain.Category) error
	getByIDFn   func(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	getByUserFn func(ctx context.Context, id uuid.UUID) ([]*domain.Category, error)
	updateFn    func(ctx context.Context, c *domain.Category) error
	deleteFn    func(ctx context.Context, id uuid.UUID) error
}

func (m *mockCategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	return m.createFn(ctx, c)
}

func (m *mockCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockCategoryRepository) GetByUserID(ctx context.Context, id uuid.UUID) ([]*domain.Category, error) {
	return m.getByUserFn(ctx, id)
}

func (m *mockCategoryRepository) Update(ctx context.Context, c *domain.Category) error {
	return m.updateFn(ctx, c)
}

func (m *mockCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}
