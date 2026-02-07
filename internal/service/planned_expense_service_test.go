package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

func TestPlannedExpenseService_Create(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	accountID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name        string
		amount      decimal.Decimal
		dueDay      int
		description string
		frequency   domain.Recurrence
		mockRepo    func() *mockPlannedExpenseRepository
		wantErr     bool
	}{
		{
			name:        "deve criar despesa planejada com sucesso",
			amount:      decimal.NewFromInt(1500),
			dueDay:      10,
			description: "Aluguel",
			frequency:   domain.Monthly,
			mockRepo: func() *mockPlannedExpenseRepository {
				return &mockPlannedExpenseRepository{
					createFn: func(_ context.Context, _ *domain.PlannedExpense) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:        "deve retornar erro para valor zero",
			amount:      decimal.Zero,
			dueDay:      10,
			description: "Teste",
			frequency:   domain.Monthly,
			mockRepo: func() *mockPlannedExpenseRepository {
				return &mockPlannedExpenseRepository{}
			},
			wantErr: true,
		},
		{
			name:        "deve retornar erro do repo",
			amount:      decimal.NewFromInt(100),
			dueDay:      10,
			description: "Teste",
			frequency:   domain.Monthly,
			mockRepo: func() *mockPlannedExpenseRepository {
				return &mockPlannedExpenseRepository{
					createFn: func(_ context.Context, _ *domain.PlannedExpense) error {
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
			svc := NewPlannedExpenseService(tt.mockRepo())
			_, err := svc.Create(context.Background(), userID, accountID, categoryID, tt.amount, tt.dueDay, nil, nil, tt.description, tt.frequency, true)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlannedExpenseService_Update(t *testing.T) {
	t.Parallel()

	peID := uuid.New()

	tests := []struct {
		name        string
		amount      decimal.Decimal
		dueDay      int
		description string
		frequency   domain.Recurrence
		mockRepo    func() *mockPlannedExpenseRepository
		wantErr     bool
	}{
		{
			name:        "deve atualizar com sucesso",
			amount:      decimal.NewFromInt(2000),
			dueDay:      15,
			description: "Aluguel reajustado",
			frequency:   domain.Monthly,
			mockRepo: func() *mockPlannedExpenseRepository {
				return &mockPlannedExpenseRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.PlannedExpense, error) {
						return &domain.PlannedExpense{ID: peID}, nil
					},
					updateFn: func(_ context.Context, _ *domain.PlannedExpense) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:        "deve retornar erro quando não encontrada",
			amount:      decimal.NewFromInt(1000),
			dueDay:      10,
			description: "Teste",
			frequency:   domain.Monthly,
			mockRepo: func() *mockPlannedExpenseRepository {
				return &mockPlannedExpenseRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.PlannedExpense, error) {
						return nil, domain.ErrPlannedExpenseNotFound
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewPlannedExpenseService(tt.mockRepo())
			_, err := svc.Update(context.Background(), peID, tt.amount, tt.dueDay, nil, nil, tt.description, tt.frequency, true)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
