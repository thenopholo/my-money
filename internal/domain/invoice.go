package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InvoiceStatus string

const (
	InvoiceStatusOpen    InvoiceStatus = "open"
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusExpired InvoiceStatus = "expired"
)

func (s InvoiceStatus) IsPaid() bool {
	return s == InvoiceStatusPaid
}

func (s InvoiceStatus) IsExpired() bool {
	return s == InvoiceStatusExpired
}

type Invoice struct {
	ID            uuid.UUID
	CardID        uuid.UUID
	ReferenceDate time.Time
	DueDate       time.Time
	TotalAmount   decimal.Decimal
	Status        InvoiceStatus
	PaidAt        *time.Time
}

func NewInvoice(cardID uuid.UUID, referenceDate, dueDate time.Time, paidAt *time.Time, totalAmount decimal.Decimal, status InvoiceStatus) (*Invoice, error) {
	if status.IsPaid() {
		return nil, ErrInvoiceAlreadyPaid
	}

	if status.IsExpired() {
		return nil, ErrInvoiceAlreadyClosed
	}

	if totalAmount.LessThan(decimal.Zero) {
		return nil, ErrInvalidInvoice
	}

	return &Invoice{
		ID:            uuid.New(),
		CardID:        cardID,
		ReferenceDate: referenceDate,
		DueDate:       dueDate,
		TotalAmount:   totalAmount,
		Status:        status,
		PaidAt:        paidAt,
	}, nil
}
