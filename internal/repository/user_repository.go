package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/repository/postgres"
)

type userRepository struct {
	queries *postgres.Queries
}

func NewUserRepository(q *postgres.Queries) *userRepository {
	return &userRepository{queries: q}
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) error {
	dbUser, err := r.queries.CreateUser(ctx, postgres.CreateUserParams{
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	})
	if err != nil {
		return err
	}

	u.ID = dbUser.ID
	u.CreatedAt = dbUser.CreatedAt.Time

	return nil
}

func (r *userRepository) Update(ctx context.Context, u *domain.User) error {
	dbUser, err := r.queries.UpdateUser(ctx, postgres.UpdateUserParams{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}

	u.UpdatedAt = dbUser.UpdatedAt.Time

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
  err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}, nil
}