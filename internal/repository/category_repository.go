package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/repository/postgres"
)

type categoryRepository struct {
	queries *postgres.Queries
}

func NewCategoryRepository(q *postgres.Queries) *categoryRepository {
	return &categoryRepository{
		queries: q,
	}
}

func (r *categoryRepository) Create(ctx context.Context, c *domain.Category) error {
	dbCategory, err := r.queries.CreateCategory(ctx, postgres.CreateCategoryParams{
		UserID:       c.UserID,
		Name:         c.Name,
		CategoryType: string(c.CategoryType),
		Color:        stringPtrToPgText(c.Color),
		Icon:         stringPtrToPgText(c.Icon),
	})
	if err != nil {
		return err
	}

	c.ID = dbCategory.ID

	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	c, err := r.queries.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}

		return nil, err
	}

	return &domain.Category{
		ID:           c.ID,
		UserID:       c.UserID,
		Name:         c.Name,
		CategoryType: domain.CategoryType(c.CategoryType),
		Color:        pgTextToStringPtr(c.Color),
		Icon:         pgTextToStringPtr(c.Icon),
	}, nil
}

func (r *categoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Category, error) {
	dbCategories, err := r.queries.GetCategoryByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, len(dbCategories))
	for i, c := range dbCategories {
		categories[i] = &domain.Category{
			ID:           c.ID,
			UserID:       c.UserID,
			Name:         c.Name,
			CategoryType: domain.CategoryType(c.CategoryType),
			Color:        pgTextToStringPtr(c.Color),
			Icon:         pgTextToStringPtr(c.Icon),
		}
	}

	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, c *domain.Category) error {
	_, err := r.queries.UpdateCategory(ctx, postgres.UpdateCategoryParams{
		ID:           c.ID,
		Name:         c.Name,
		CategoryType: string(c.CategoryType),
		Color:        stringPtrToPgText(c.Color),
		Icon:         stringPtrToPgText(c.Icon),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrCategoryNotFound
		}

		return err
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.queries.GetCategoryByID(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrCategoryNotFound
		}

		return err
	}

	return r.queries.DeleteCategory(ctx, id)
}
