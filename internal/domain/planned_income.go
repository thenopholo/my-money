package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PlannedIncome struct {
	ID          uuid.UUID
	AccountID   uuid.UUID
	UserID      uuid.UUID
	CategoryID  uuid.UUID
	Amount      decimal.Decimal
	DueDay      int
	StartDate   *time.Time
	EndDate     *time.Time
	Frequency   string
	Description string
	IsActive    bool
	CreatedAt   time.Time
}

func NewPlannedIncome(userID, AccountID, categoryID uuid.UUID, amount decimal.Decimal, dueDay int, startDate, endDate *time.Time, description, frequency string, isActive bool) *PlannedIncome {
	return &PlannedIncome{
		ID:          uuid.New(),
		AccountID:   AccountID,
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
	}
}
