package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

func TestPlannedIncomeService_Create(t *testing.T) {
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
		mockRepo    func() *mockPlannedIncomeRepository
		wantErr     bool
	}{
		{
			name:        "deve criar receita planejada com sucesso",
			amount:      decimal.NewFromInt(5000),
			dueDay:      5,
			description: "Salário",
			frequency:   domain.Monthly,
			mockRepo: func() *mockPlannedIncomeRepository {
				return &mockPlannedIncomeRepository{
					createFn: func(_ context.Context, _ *domain.PlannedIncome) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:        "deve retornar erro para valor zero",
			amount:      decimal.Zero,
			dueDay:      5,
			description: "Teste",
			frequency:   domain.Monthly,
			mockRepo: func() *mockPlannedIncomeRepository {
				return &mockPlannedIncomeRepository{}
			},
			wantErr: true,
		},
		{
			name:        "deve retornar erro do repo",
			amount:      decimal.NewFromInt(1000),
			dueDay:      5,
			description: "Teste",
			frequency:   domain.Monthly,
			mockRepo: func() *mockPlannedIncomeRepository {
				return &mockPlannedIncomeRepository{
					createFn: func(_ context.Context, _ *domain.PlannedIncome) error {
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
			svc := NewPlannedIncomeService(tt.mockRepo())
			_, err := svc.Create(context.Background(), userID, accountID, categoryID, tt.amount, tt.dueDay, nil, nil, tt.description, tt.frequency, true)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlannedIncomeService_Update(t *testing.T) {
	t.Parallel()

	piID := uuid.New()

	tests := []struct {
		name        string
		amount      decimal.Decimal
		dueDay      int
		description string
		frequency   domain.Recurrence
		mockRepo    func() *mockPlannedIncomeRepository
		wantErr     bool
	}{
		{
			name:        "deve atualizar com sucesso",
			amount:      decimal.NewFromInt(6000),
			dueDay:      10,
			description: "Salário reajustado",
			frequency:   domain.Monthly,
			mockRepo: func() *mockPlannedIncomeRepository {
				return &mockPlannedIncomeRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.PlannedIncome, error) {
						return &domain.PlannedIncome{ID: piID}, nil
					},
					updateFn: func(_ context.Context, _ *domain.PlannedIncome) error {
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
			mockRepo: func() *mockPlannedIncomeRepository {
				return &mockPlannedIncomeRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.PlannedIncome, error) {
						return &domain.PlannedIncome{ID: piID}, nil
					},
				}
			},
			wantErr: true,
		},
		{
			name:        "deve retornar erro para descrição vazia",
			amount:      decimal.NewFromInt(1000),
			dueDay:      10,
			description: "",
			frequency:   domain.Monthly,
			mockRepo: func() *mockPlannedIncomeRepository {
				return &mockPlannedIncomeRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.PlannedIncome, error) {
						return &domain.PlannedIncome{ID: piID}, nil
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewPlannedIncomeService(tt.mockRepo())
			_, err := svc.Update(context.Background(), piID, tt.amount, tt.dueDay, nil, nil, tt.description, tt.frequency, true)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
