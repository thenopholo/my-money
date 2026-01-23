package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Recurrence string

const (
	Once    Recurrence = "once"
	Monthly Recurrence = "monthly"
	Yearly  Recurrence = "yearly"
)

func (r Recurrence) IsValid() bool {
	return r == Once || r == Monthly || r == Yearly
}

type PlannedIncome struct {
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

func NewPlannedIncome(userID, accountID, categoryID uuid.UUID, amount decimal.Decimal, dueDay int, startDate, endDate *time.Time, description string, frequency Recurrence, isActive bool) (*PlannedIncome, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidAmount
	}

	if !frequency.IsValid() {
		return nil, ErrInvalidFrequency
	}

	if dueDay < 1 || dueDay > 31 {
		return nil, ErrInvalidDueDay
	}

	if endDate != nil && startDate != nil && endDate.Before(*startDate) {
    return nil, ErrEndDateBeforeStart
}

	if description == "" {
		return nil, ErrEmptyDescription
	}

	return &PlannedIncome{
		ID:          uuid.New(),
		AccountID:   accountID,
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      amount,
		DueDay:      dueDay,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		Frequency:   frequency,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
	}, nil
}
