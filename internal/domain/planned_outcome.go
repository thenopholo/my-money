package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PlannedOutcome struct {
	ID         uuid.UUID
	AccountID  uuid.UUID
	UserID     uuid.UUID
	CategoryID uuid.UUID
	Amount     decimal.Decimal
	DueDay     int
	StartDate  *time.Time
	EndDate    *time.Time
	Frequency  string
	IsActive   bool
	CreatedAt  time.Time
}

func NewPlannedOutcome(userID, AccountID, categoryID uuid.UUID, amount decimal.Decimal, dueDay int, startDate, endDate *time.Time, frequency string, isActive bool) *PlannedOutcome {
	return &PlannedOutcome{
		ID:         uuid.New(),
		AccountID:  AccountID,
		UserID:     userID,
		CategoryID: categoryID,
		Amount:     amount,
		DueDay:     dueDay,
		StartDate:  startDate,
		EndDate:    endDate,
		Frequency:  frequency,
		IsActive:   isActive,
		CreatedAt:  time.Now(),
	}
}
