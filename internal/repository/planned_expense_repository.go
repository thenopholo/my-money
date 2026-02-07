package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/repository/postgres"
)

type plannedExpenseRepository struct {
	queries *postgres.Queries
}

func NewPlannedExpenseRepository(q *postgres.Queries) *plannedExpenseRepository {
	return &plannedExpenseRepository{queries: q}
}

func (r *plannedExpenseRepository) Create(ctx context.Context, pe *domain.PlannedExpense) error {
	dbExpense, err := r.queries.CreatePlannedExpense(ctx, postgres.CreatePlannedExpenseParams{
		UserID:      pe.UserID,
		AccountID:   pe.AccountID,
		CategoryID:  pe.CategoryID,
		Amount:      decimalToPgNumeric(pe.Amount),
		DueDay:      int32(pe.DueDay),
		StartDate:   timePtrToPgDate(pe.StartDate),
		EndDate:     timePtrToPgDate(pe.EndDate),
		Description: pe.Description,
		Frequency:   string(pe.Frequency),
		IsActive:    pe.IsActive,
	})
	if err != nil {
		return err
	}

	pe.ID = dbExpense.ID
	pe.CreatedAt = dbExpense.CreatedAt.Time

	return nil
}

func (r *plannedExpenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PlannedExpense, error) {
	pe, err := r.queries.GetPlannedExpenseByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPlannedExpenseNotFound
		}

		return nil, err
	}

	return &domain.PlannedExpense{
		ID:          pe.ID,
		UserID:      pe.UserID,
		AccountID:   pe.AccountID,
		CategoryID:  pe.CategoryID,
		Description: pe.Description,
		Amount:      pgNumericToDecimal(pe.Amount),
		DueDay:      int(pe.DueDay),
		Frequency:   domain.Recurrence(pe.Frequency),
		StartDate:   pgDateToTimePtr(pe.StartDate),
		EndDate:     pgDateToTimePtr(pe.EndDate),
		IsActive:    pe.IsActive,
		CreatedAt:   pe.CreatedAt.Time,
	}, nil
}

func (r *plannedExpenseRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.PlannedExpense, error) {
	dbExpenses, err := r.queries.GetPlannedExpenseByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	expenses := make([]*domain.PlannedExpense, len(dbExpenses))
	for i, pe := range dbExpenses {
		expenses[i] = &domain.PlannedExpense{
			ID:          pe.ID,
			UserID:      pe.UserID,
			AccountID:   pe.AccountID,
			CategoryID:  pe.CategoryID,
			Description: pe.Description,
			Amount:      pgNumericToDecimal(pe.Amount),
			DueDay:      int(pe.DueDay),
			Frequency:   domain.Recurrence(pe.Frequency),
			StartDate:   pgDateToTimePtr(pe.StartDate),
			EndDate:     pgDateToTimePtr(pe.EndDate),
			IsActive:    pe.IsActive,
			CreatedAt:   pe.CreatedAt.Time,
		}
	}

	return expenses, nil
}

func (r *plannedExpenseRepository) Update(ctx context.Context, pe *domain.PlannedExpense) error {
	_, err := r.queries.UpdatePlannedExpense(ctx, postgres.UpdatePlannedExpenseParams{
		ID:          pe.ID,
		Amount:      decimalToPgNumeric(pe.Amount),
		DueDay:      int32(pe.DueDay),
		StartDate:   timePtrToPgDate(pe.StartDate),
		EndDate:     timePtrToPgDate(pe.EndDate),
		Description: pe.Description,
		Frequency:   string(pe.Frequency),
		IsActive:    pe.IsActive,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrPlannedExpenseNotFound
		}

		return err
	}

	return nil
}

func (r *plannedExpenseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.queries.GetPlannedExpenseByID(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrPlannedExpenseNotFound
		}

		return err
	}

	return r.queries.DeletePlannedExpense(ctx, id)
}
