package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestNewPlannedExpense_Validations(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	accountID := uuid.New()
	categoryID := uuid.New()
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	endBeforeStart := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		amount      decimal.Decimal
		dueDay      int
		startDate   *time.Time
		endDate     *time.Time
		description string
		frequency   Recurrence
		wantErr     error
	}{
		{
			name:        "deve criar despesa planejada válida",
			amount:      decimal.NewFromInt(1500),
			dueDay:      10,
			startDate:   &start,
			endDate:     &end,
			description: "Aluguel",
			frequency:   Monthly,
			wantErr:     nil,
		},
		{
			name:        "deve retornar erro para valor zero",
			amount:      decimal.Zero,
			dueDay:      10,
			startDate:   nil,
			endDate:     nil,
			description: "Teste",
			frequency:   Monthly,
			wantErr:     ErrInvalidAmount,
		},
		{
			name:        "deve retornar erro para valor negativo",
			amount:      decimal.NewFromInt(-100),
			dueDay:      10,
			startDate:   nil,
			endDate:     nil,
			description: "Teste",
			frequency:   Monthly,
			wantErr:     ErrInvalidAmount,
		},
		{
			name:        "deve retornar erro para frequência inválida",
			amount:      decimal.NewFromInt(100),
			dueDay:      10,
			startDate:   nil,
			endDate:     nil,
			description: "Teste",
			frequency:   Recurrence("weekly"),
			wantErr:     ErrInvalidFrequency,
		},
		{
			name:        "deve retornar erro para due day 0",
			amount:      decimal.NewFromInt(100),
			dueDay:      0,
			startDate:   nil,
			endDate:     nil,
			description: "Teste",
			frequency:   Monthly,
			wantErr:     ErrInvalidDueDay,
		},
		{
			name:        "deve retornar erro para due day 32",
			amount:      decimal.NewFromInt(100),
			dueDay:      32,
			startDate:   nil,
			endDate:     nil,
			description: "Teste",
			frequency:   Monthly,
			wantErr:     ErrInvalidDueDay,
		},
		{
			name:        "deve retornar erro para end date antes de start date",
			amount:      decimal.NewFromInt(100),
			dueDay:      10,
			startDate:   &start,
			endDate:     &endBeforeStart,
			description: "Teste",
			frequency:   Monthly,
			wantErr:     ErrEndDateBeforeStart,
		},
		{
			name:        "deve retornar erro para descrição vazia",
			amount:      decimal.NewFromInt(100),
			dueDay:      10,
			startDate:   nil,
			endDate:     nil,
			description: "",
			frequency:   Monthly,
			wantErr:     ErrEmptyDescription,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pe, err := NewPlannedExpense(
				userID, accountID, categoryID,
				tt.amount, tt.dueDay, tt.startDate, tt.endDate,
				tt.description, tt.frequency, true,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewPlannedExpense() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && pe == nil {
				t.Fatal("NewPlannedExpense() retornou nil sem erro")
			}
		})
	}
}

func TestPlannedExpense_IsActiveOn(t *testing.T) {
	t.Parallel()

	start := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		isActive  bool
		startDate *time.Time
		endDate   *time.Time
		checkDate time.Time
		want      bool
	}{
		{
			name:      "ativo dentro do período",
			isActive:  true,
			startDate: &start,
			endDate:   &end,
			checkDate: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			want:      true,
		},
		{
			name:      "inativo se IsActive é false",
			isActive:  false,
			startDate: &start,
			endDate:   &end,
			checkDate: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			want:      false,
		},
		{
			name:      "inativo antes do start date",
			isActive:  true,
			startDate: &start,
			endDate:   &end,
			checkDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			want:      false,
		},
		{
			name:      "inativo após end date",
			isActive:  true,
			startDate: &start,
			endDate:   &end,
			checkDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pe := &PlannedExpense{
				IsActive:  tt.isActive,
				StartDate: tt.startDate,
				EndDate:   tt.endDate,
			}

			if got := pe.IsActiveOn(tt.checkDate); got != tt.want {
				t.Errorf("IsActiveOn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlannedExpense_IsDueOn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		dueDay int
		date   time.Time
		want   bool
	}{
		{
			name:   "é dia de vencimento",
			dueDay: 10,
			date:   time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC),
			want:   true,
		},
		{
			name:   "não é dia de vencimento",
			dueDay: 10,
			date:   time.Date(2025, 6, 11, 0, 0, 0, 0, time.UTC),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pe := &PlannedExpense{DueDay: tt.dueDay}
			if got := pe.IsDueOn(tt.date); got != tt.want {
				t.Errorf("IsDueOn() = %v, want %v", got, tt.want)
			}
		})
	}
}
