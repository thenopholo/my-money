package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestRecurrence_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    Recurrence
		want bool
	}{
		{"once é válido", Once, true},
		{"monthly é válido", Monthly, true},
		{"yearly é válido", Yearly, true},
		{"string vazia é inválido", Recurrence(""), false},
		{"tipo desconhecido é inválido", Recurrence("weekly"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.r.IsValid(); got != tt.want {
				t.Errorf("Recurrence(%q).IsValid() = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

func TestNewPlannedIncome_Validations(t *testing.T) {
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
			name:        "deve criar receita planejada válida",
			amount:      decimal.NewFromInt(5000),
			dueDay:      5,
			startDate:   &start,
			endDate:     &end,
			description: "Salário",
			frequency:   Monthly,
			wantErr:     nil,
		},
		{
			name:        "deve criar sem datas opcionais",
			amount:      decimal.NewFromInt(1000),
			dueDay:      15,
			startDate:   nil,
			endDate:     nil,
			description: "Freelance",
			frequency:   Once,
			wantErr:     nil,
		},
		{
			name:        "deve retornar erro para valor zero",
			amount:      decimal.Zero,
			dueDay:      5,
			startDate:   nil,
			endDate:     nil,
			description: "Teste",
			frequency:   Monthly,
			wantErr:     ErrInvalidAmount,
		},
		{
			name:        "deve retornar erro para valor negativo",
			amount:      decimal.NewFromInt(-100),
			dueDay:      5,
			startDate:   nil,
			endDate:     nil,
			description: "Teste",
			frequency:   Monthly,
			wantErr:     ErrInvalidAmount,
		},
		{
			name:        "deve retornar erro para frequência inválida",
			amount:      decimal.NewFromInt(1000),
			dueDay:      5,
			startDate:   nil,
			endDate:     nil,
			description: "Teste",
			frequency:   Recurrence("weekly"),
			wantErr:     ErrInvalidFrequency,
		},
		{
			name:        "deve retornar erro para due day 0",
			amount:      decimal.NewFromInt(1000),
			dueDay:      0,
			startDate:   nil,
			endDate:     nil,
			description: "Teste",
			frequency:   Monthly,
			wantErr:     ErrInvalidDueDay,
		},
		{
			name:        "deve retornar erro para due day 32",
			amount:      decimal.NewFromInt(1000),
			dueDay:      32,
			startDate:   nil,
			endDate:     nil,
			description: "Teste",
			frequency:   Monthly,
			wantErr:     ErrInvalidDueDay,
		},
		{
			name:        "deve retornar erro para end date antes de start date",
			amount:      decimal.NewFromInt(1000),
			dueDay:      5,
			startDate:   &start,
			endDate:     &endBeforeStart,
			description: "Teste",
			frequency:   Monthly,
			wantErr:     ErrEndDateBeforeStart,
		},
		{
			name:        "deve retornar erro para descrição vazia",
			amount:      decimal.NewFromInt(1000),
			dueDay:      5,
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
			pi, err := NewPlannedIncome(
				userID, accountID, categoryID,
				tt.amount, tt.dueDay, tt.startDate, tt.endDate,
				tt.description, tt.frequency, true,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewPlannedIncome() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if pi == nil {
					t.Fatal("NewPlannedIncome() retornou nil sem erro")
				}
				if !pi.Amount.Equal(tt.amount) {
					t.Errorf("Amount = %v, want %v", pi.Amount, tt.amount)
				}
				if pi.DueDay != tt.dueDay {
					t.Errorf("DueDay = %d, want %d", pi.DueDay, tt.dueDay)
				}
			}
		})
	}
}

func TestPlannedIncome_IsActiveOn(t *testing.T) {
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
		{
			name:      "ativo sem datas opcionais",
			isActive:  true,
			startDate: nil,
			endDate:   nil,
			checkDate: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pi := &PlannedIncome{
				IsActive:  tt.isActive,
				StartDate: tt.startDate,
				EndDate:   tt.endDate,
			}

			if got := pi.IsActiveOn(tt.checkDate); got != tt.want {
				t.Errorf("IsActiveOn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlannedIncome_IsDueOn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		dueDay int
		date   time.Time
		want   bool
	}{
		{
			name:   "é dia de vencimento",
			dueDay: 15,
			date:   time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			want:   true,
		},
		{
			name:   "não é dia de vencimento",
			dueDay: 15,
			date:   time.Date(2025, 6, 16, 0, 0, 0, 0, time.UTC),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pi := &PlannedIncome{DueDay: tt.dueDay}
			if got := pi.IsDueOn(tt.date); got != tt.want {
				t.Errorf("IsDueOn() = %v, want %v", got, tt.want)
			}
		})
	}
}
