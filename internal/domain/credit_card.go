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
	UpdatedAt   time.Time
}

func NewCreditCard(userID uuid.UUID, name string, creditLimit decimal.Decimal, closeDay, dueDay int) (*CreditCard, error) {
	cc := &CreditCard{
		ID:        uuid.New(),
		UserID:    userID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := cc.SetName(name); err != nil {
		return nil, err
	}

	if err := cc.SetCreditLimit(creditLimit); err != nil {
		return nil, err
	}

	if err := cc.SetCloseDay(closeDay); err != nil {
		return nil, err
	}

	if err := cc.SetDueDay(dueDay); err != nil {
		return nil, err
	}

	return cc, nil
}

func (cc *CreditCard) SetName(name string) error {
	if name == "" {
		return ErrEmptyCardName
	}

	cc.Name = name
	return nil
}

func (cc *CreditCard) SetCreditLimit(limit decimal.Decimal) error {
	if limit.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidLimit
	}

	cc.CreditLimit = limit
	return nil
}

func (cc *CreditCard) SetCloseDay(closeDay int) error {
	if closeDay < 1 || closeDay > 28 {
		return ErrInvalidCloseDay
	}

	cc.CloseDay = closeDay
	return nil
}

func (cc *CreditCard) SetDueDay(dueDay int) error {
	if dueDay < 1 || dueDay > 31 {
		return ErrInvalidDueDay
	}

	cc.DueDay = dueDay
	return nil
}

func (cc *CreditCard) SetIsActive(isActive bool) {
	cc.IsActive = isActive
}
