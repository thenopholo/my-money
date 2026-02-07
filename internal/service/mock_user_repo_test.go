package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

type mockUserRepository struct {
	createFn     func(ctx context.Context, user *domain.User) error
	getByIDFn    func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	getByEmailFn func(ctx context.Context, email string) (*domain.User, error)
	updateFn     func(ctx context.Context, user *domain.User) error
	deleteFn     func(ctx context.Context, id uuid.UUID) error
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.createFn(ctx, user)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.getByEmailFn(ctx, email)
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	return m.updateFn(ctx, user)
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFn(ctx, id)
}
