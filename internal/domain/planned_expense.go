package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PlannedExpense struct {
	ID          uuid.UUID
	AccountID   uuid.UUID
	UserID      uuid.UUID
	CategoryID  uuid.UUID
	Amount      decimal.Decimal
	DueDay      int
	StartDate   *time.Time
	EndDate     *time.Time
	Frequency   Recurrence
	Description string
	IsActive    bool
	CreatedAt   time.Time
}

func NewPlannedExpense(userID, AccountID, categoryID uuid.UUID, amount decimal.Decimal, dueDay int, startDate, endDate *time.Time, description string, frequency Recurrence, isActive bool) (*PlannedExpense, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidAmount
	}

	if !frequency.IsValid() {
		return nil, ErrInvalidFrequency
	}

	if dueDay < 1 || dueDay > 32 {
		return nil, ErrInvalidDueDay
	}

	if endDate.After(*startDate) {
		return nil, ErrEndDateBeforeStart
	}

	if description == "" {
		return nil, ErrEmptyDescription
	}

	return &PlannedExpense{
		ID:          uuid.New(),
		AccountID:   AccountID,
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      amount,
		DueDay:      dueDay,
		StartDate:   startDate,
		EndDate:     endDate,
		Frequency:   frequency,
		Description: description,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
	}, nil
}
