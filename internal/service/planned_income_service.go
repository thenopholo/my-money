package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

type PlannedIncomeService struct {
	repo PlannedIncomeRepository
}

func NewPlannedIncomeService(repo PlannedIncomeRepository) *PlannedIncomeService {
	return &PlannedIncomeService{repo: repo}
}

func (s *PlannedIncomeService) Create(
	ctx context.Context,
	userID, accountID, categoryID uuid.UUID,
	amount decimal.Decimal,
	dueDay int,
	startDate, endDate *time.Time,
	description string,
	frequency domain.Recurrence,
	isActive bool,
) (*domain.PlannedIncome, error) {
	planned, err := domain.NewPlannedIncome(userID, accountID, categoryID, amount, dueDay, startDate, endDate, description, frequency, isActive)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, planned); err != nil {
		return nil, err
	}

	return planned, nil
}

func (s *PlannedIncomeService) GetByID(ctx context.Context, id uuid.UUID) (*domain.PlannedIncome, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PlannedIncomeService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.PlannedIncome, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *PlannedIncomeService) Update(
	ctx context.Context,
	id uuid.UUID,
	amount decimal.Decimal,
	dueDay int,
	startDate, endDate *time.Time,
	description string,
	frequency domain.Recurrence,
	isActive bool,
) (*domain.PlannedIncome, error) {
	planned, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, domain.ErrInvalidAmount
	}
	if dueDay < 1 || dueDay > 31 {
		return nil, domain.ErrInvalidDueDay
	}
	if !frequency.IsValid() {
		return nil, domain.ErrInvalidFrequency
	}
	if description == "" {
		return nil, domain.ErrEmptyDescription
	}
	if endDate != nil && startDate != nil && endDate.Before(*startDate) {
		return nil, domain.ErrEndDateBeforeStart
	}

	planned.Amount = amount
	planned.DueDay = dueDay
	planned.StartDate = startDate
	planned.EndDate = endDate
	planned.Description = description
	planned.Frequency = frequency
	planned.IsActive = isActive

	if err := s.repo.Update(ctx, planned); err != nil {
		return nil, err
	}

	return planned, nil
}

func (s *PlannedIncomeService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
