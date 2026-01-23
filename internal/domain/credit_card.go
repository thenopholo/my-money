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
}

func NewCreditCard(userID uuid.UUID, name string, creditLimit decimal.Decimal, closeDay, dueDay int, isActive bool) (*CreditCard, error) {
	if name == "" {
		return nil, ErrEmptyCardName
	}

	if closeDay < 1 || closeDay > 29 {
		return nil, ErrInvalidCloseDay
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
		IsActive:    isActive,
		CreatedAt:   time.Now(),
	}, nil
}
