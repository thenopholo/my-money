package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreditCard struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	CreditLimit decimal.Decimal
	CloseDay    int
	DueDay      int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt  time.Time
}

func NewCreditCard(userID uuid.UUID, name string, creditLimit decimal.Decimal, closeDay, dueDay int) (*CreditCard, error) {
	if name == "" {
		return nil, ErrEmptyCardName
	}

	if closeDay < 1 || closeDay > 28 {
		return nil, ErrInvalidCloseDay
	}

	if dueDay < 1 || dueDay > 31 {
		return nil, ErrInvalidDueDay
	}

	if creditLimit.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidLimit
	}

	return &CreditCard{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		CreditLimit: creditLimit,
		CloseDay:    closeDay,
		DueDay:      dueDay,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}
