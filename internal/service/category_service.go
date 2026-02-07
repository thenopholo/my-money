package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

type CategoryService struct {
	categoryRepo CategoryRepository
}

func NewCategoryService(categoryRepo CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (cs *CategoryService) Create(ctx context.Context, userID uuid.UUID, categoryType domain.CategoryType, name string, color, icon *string) (*domain.Category, error) {
	category, err := domain.NewCategory(userID, categoryType, name, color, icon)
	if err != nil {
		return nil, err
	}

	if err := cs.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (cs *CategoryService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return cs.categoryRepo.GetByID(ctx, id)
}

func (cs *CategoryService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Category, error) {
	return cs.categoryRepo.GetByUserID(ctx, userID)
}

func (cs *CategoryService) Update(ctx context.Context, id uuid.UUID, categoryType domain.CategoryType, name string, color, icon *string) (*domain.Category, error) {
	category, err := cs.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := category.SetName(name); err != nil {
		return nil, err
	}
	if err := category.SetType(categoryType); err != nil {
		return nil, err
	}
	category.SetVisuals(color, icon)

	if err := cs.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (cs *CategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	return cs.categoryRepo.Delete(ctx, id)
}
