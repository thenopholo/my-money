package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Invoice struct {
	ID            uuid.UUID
	CardID        uuid.UUID
	ReferenceDate time.Time
	DueDate       time.Time
	TotalAmount   decimal.Decimal
	Status        string
	PaidAt        time.Time
}

func NewInvoice(cardID uuid.UUID, referenceDate, dueDate, paidAt time.Time, totalAmount decimal.Decimal, status string) *Invoice {
	return &Invoice{
		ID:            uuid.New(),
		CardID:        cardID,
		ReferenceDate: referenceDate,
		DueDate:       dueDate,
		TotalAmount:   totalAmount,
		Status:        status,
		PaidAt:        paidAt,
	}
}
