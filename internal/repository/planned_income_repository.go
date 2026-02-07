package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/repository/postgres"
)

type plannedIncomeRepository struct {
	queries *postgres.Queries
}

func NewPlannedIncomeRepository(q *postgres.Queries) *plannedIncomeRepository {
	return &plannedIncomeRepository{queries: q}
}

func (r *plannedIncomeRepository) Create(ctx context.Context, pi *domain.PlannedIncome) error {
	dbIncome, err := r.queries.CreatePlannedIncome(ctx, postgres.CreatePlannedIncomeParams{
		UserID:      pi.UserID,
		AccountID:   pi.AccountID,
		CategoryID:  pi.CategoryID,
		Amount:      decimalToPgNumeric(pi.Amount),
		DueDay:      int32(pi.DueDay),
		StartDate:   timePtrToPgDate(pi.StartDate),
		EndDate:     timePtrToPgDate(pi.EndDate),
		Description: pi.Description,
		Frequency:   string(pi.Frequency),
		IsActive:    pi.IsActive,
	})
	if err != nil {
		return err
	}

	pi.ID = dbIncome.ID
	pi.CreatedAt = dbIncome.CreatedAt.Time

	return nil
}

func (r *plannedIncomeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PlannedIncome, error) {
	pi, err := r.queries.GetPlannedIncomeByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPlannedIncomeNotFound
		}

		return nil, err
	}

	return &domain.PlannedIncome{
		ID:          pi.ID,
		UserID:      pi.UserID,
		AccountID:   pi.AccountID,
		CategoryID:  pi.CategoryID,
		Description: pi.Description,
		Amount:      pgNumericToDecimal(pi.Amount),
		DueDay:      int(pi.DueDay),
		Frequency:   domain.Recurrence(pi.Frequency),
		StartDate:   pgDateToTimePtr(pi.StartDate),
		EndDate:     pgDateToTimePtr(pi.EndDate),
		IsActive:    pi.IsActive,
		CreatedAt:   pi.CreatedAt.Time,
	}, nil
}

func (r *plannedIncomeRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.PlannedIncome, error) {
	dbIncomes, err := r.queries.GetPlannedIncomeByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	incomes := make([]*domain.PlannedIncome, len(dbIncomes))
	for i, pi := range dbIncomes {
		incomes[i] = &domain.PlannedIncome{
			ID:          pi.ID,
			UserID:      pi.UserID,
			AccountID:   pi.AccountID,
			CategoryID:  pi.CategoryID,
			Description: pi.Description,
			Amount:      pgNumericToDecimal(pi.Amount),
			DueDay:      int(pi.DueDay),
			Frequency:   domain.Recurrence(pi.Frequency),
			StartDate:   pgDateToTimePtr(pi.StartDate),
			EndDate:     pgDateToTimePtr(pi.EndDate),
			IsActive:    pi.IsActive,
			CreatedAt:   pi.CreatedAt.Time,
		}
	}

	return incomes, nil
}

func (r *plannedIncomeRepository) Update(ctx context.Context, pi *domain.PlannedIncome) error {
	_, err := r.queries.UpdatePlannedIncome(ctx, postgres.UpdatePlannedIncomeParams{
		ID:          pi.ID,
		Amount:      decimalToPgNumeric(pi.Amount),
		DueDay:      int32(pi.DueDay),
		StartDate:   timePtrToPgDate(pi.StartDate),
		EndDate:     timePtrToPgDate(pi.EndDate),
		Description: pi.Description,
		Frequency:   string(pi.Frequency),
		IsActive:    pi.IsActive,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrPlannedIncomeNotFound
		}

		return err
	}

	return nil
}

func (r *plannedIncomeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.queries.GetPlannedIncomeByID(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrPlannedIncomeNotFound
		}

		return err
	}

	return r.queries.DeletePlannedIncome(ctx, id)
}
