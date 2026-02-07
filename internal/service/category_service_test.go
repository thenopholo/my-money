package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

func TestCategoryService_Create(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	tests := []struct {
		name         string
		categoryType domain.CategoryType
		catName      string
		mockRepo     func() *mockCategoryRepository
		wantErr      bool
	}{
		{
			name:         "deve criar categoria com sucesso",
			categoryType: domain.CategoryTypeIncome,
			catName:      "Salário",
			mockRepo: func() *mockCategoryRepository {
				return &mockCategoryRepository{
					createFn: func(_ context.Context, _ *domain.Category) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:         "deve retornar erro para tipo inválido",
			categoryType: domain.CategoryType("invalid"),
			catName:      "Teste",
			mockRepo: func() *mockCategoryRepository {
				return &mockCategoryRepository{}
			},
			wantErr: true,
		},
		{
			name:         "deve retornar erro para nome vazio",
			categoryType: domain.CategoryTypeIncome,
			catName:      "",
			mockRepo: func() *mockCategoryRepository {
				return &mockCategoryRepository{}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewCategoryService(tt.mockRepo())
			cat, err := svc.Create(context.Background(), userID, tt.categoryType, tt.catName, nil, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && cat == nil {
				t.Error("Create() retornou nil sem erro")
			}
		})
	}
}

func TestCategoryService_Update(t *testing.T) {
	t.Parallel()

	catID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name         string
		categoryType domain.CategoryType
		catName      string
		mockRepo     func() *mockCategoryRepository
		wantErr      bool
	}{
		{
			name:         "deve atualizar com sucesso",
			categoryType: domain.CategoryTypeExpense,
			catName:      "Alimentação",
			mockRepo: func() *mockCategoryRepository {
				return &mockCategoryRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Category, error) {
						return &domain.Category{
							ID:           catID,
							UserID:       userID,
							Name:         "Comida",
							CategoryType: domain.CategoryTypeExpense,
						}, nil
					},
					updateFn: func(_ context.Context, _ *domain.Category) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:         "deve retornar erro quando não encontrada",
			categoryType: domain.CategoryTypeIncome,
			catName:      "Teste",
			mockRepo: func() *mockCategoryRepository {
				return &mockCategoryRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Category, error) {
						return nil, domain.ErrCategoryNotFound
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewCategoryService(tt.mockRepo())
			_, err := svc.Update(context.Background(), catID, tt.categoryType, tt.catName, nil, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategoryService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mockRepo func() *mockCategoryRepository
		wantErr  bool
	}{
		{
			name: "deve deletar com sucesso",
			mockRepo: func() *mockCategoryRepository {
				return &mockCategoryRepository{
					deleteFn: func(_ context.Context, _ uuid.UUID) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "deve retornar erro do repo",
			mockRepo: func() *mockCategoryRepository {
				return &mockCategoryRepository{
					deleteFn: func(_ context.Context, _ uuid.UUID) error {
						return errors.New("db error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewCategoryService(tt.mockRepo())
			err := svc.Delete(context.Background(), uuid.New())
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
