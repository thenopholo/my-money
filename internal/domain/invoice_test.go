package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestInvoiceStatus_Methods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		status    InvoiceStatus
		isOpen    bool
		isClosed  bool
		isPaid    bool
		isExpired bool
	}{
		{"open", InvoiceStatusOpen, true, false, false, false},
		{"closed", InvoiceStatusClosed, false, true, false, false},
		{"paid", InvoiceStatusPaid, false, false, true, false},
		{"expired", InvoiceStatusExpired, false, false, false, true},
		{"unknown", InvoiceStatus("unknown"), false, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.status.IsOpen(); got != tt.isOpen {
				t.Errorf("IsOpen() = %v, want %v", got, tt.isOpen)
			}
			if got := tt.status.IsClosed(); got != tt.isClosed {
				t.Errorf("IsClosed() = %v, want %v", got, tt.isClosed)
			}
			if got := tt.status.IsPaid(); got != tt.isPaid {
				t.Errorf("IsPaid() = %v, want %v", got, tt.isPaid)
			}
			if got := tt.status.IsExpired(); got != tt.isExpired {
				t.Errorf("IsExpired() = %v, want %v", got, tt.isExpired)
			}
		})
	}
}

func TestNewInvoice_Validations(t *testing.T) {
	t.Parallel()

	cardID := uuid.New()
	ref := time.Now()
	due := ref.AddDate(0, 1, 0)

	tests := []struct {
		name        string
		cardID      uuid.UUID
		refDate     time.Time
		dueDate     time.Time
		paidAt      *time.Time
		totalAmount decimal.Decimal
		status      InvoiceStatus
		wantErr     error
	}{
		{
			name:        "deve criar invoice open com valor positivo",
			cardID:      cardID,
			refDate:     ref,
			dueDate:     due,
			paidAt:      nil,
			totalAmount: decimal.NewFromInt(500),
			status:      InvoiceStatusOpen,
			wantErr:     nil,
		},
		{
			name:        "deve criar invoice open com valor zero",
			cardID:      cardID,
			refDate:     ref,
			dueDate:     due,
			paidAt:      nil,
			totalAmount: decimal.Zero,
			status:      InvoiceStatusOpen,
			wantErr:     nil,
		},
		{
			name:        "deve criar invoice closed",
			cardID:      cardID,
			refDate:     ref,
			dueDate:     due,
			paidAt:      nil,
			totalAmount: decimal.NewFromInt(200),
			status:      InvoiceStatusClosed,
			wantErr:     nil,
		},
		{
			name:        "deve retornar erro para status paid",
			cardID:      cardID,
			refDate:     ref,
			dueDate:     due,
			paidAt:      nil,
			totalAmount: decimal.NewFromInt(500),
			status:      InvoiceStatusPaid,
			wantErr:     ErrInvoiceAlreadyPaid,
		},
		{
			name:        "deve retornar erro para status expired",
			cardID:      cardID,
			refDate:     ref,
			dueDate:     due,
			paidAt:      nil,
			totalAmount: decimal.NewFromInt(500),
			status:      InvoiceStatusExpired,
			wantErr:     ErrInvoiceAlreadyClosed,
		},
		{
			name:        "deve retornar erro para valor negativo",
			cardID:      cardID,
			refDate:     ref,
			dueDate:     due,
			paidAt:      nil,
			totalAmount: decimal.NewFromInt(-100),
			status:      InvoiceStatusOpen,
			wantErr:     ErrInvalidInvoice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			inv, err := NewInvoice(tt.cardID, tt.refDate, tt.dueDate, tt.paidAt, tt.totalAmount, tt.status)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewInvoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if inv == nil {
					t.Fatal("NewInvoice() retornou nil sem erro")
				}
				if inv.CardID != tt.cardID {
					t.Errorf("CardID = %v, want %v", inv.CardID, tt.cardID)
				}
				if !inv.TotalAmount.Equal(tt.totalAmount) {
					t.Errorf("TotalAmount = %v, want %v", inv.TotalAmount, tt.totalAmount)
				}
				if inv.Status != tt.status {
					t.Errorf("Status = %q, want %q", inv.Status, tt.status)
				}
			}
		})
	}
}

func TestInvoice_MarkAsPaid(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name    string
		status  InvoiceStatus
		wantErr error
	}{
		{
			name:    "deve marcar como paga uma invoice open",
			status:  InvoiceStatusOpen,
			wantErr: nil,
		},
		{
			name:    "deve retornar erro para invoice já paga",
			status:  InvoiceStatusPaid,
			wantErr: ErrInvoiceAlreadyPaid,
		},
		{
			name:    "deve retornar erro para invoice expired",
			status:  InvoiceStatusExpired,
			wantErr: ErrInvoiceAlreadyClosed,
		},
		{
			name:    "deve retornar erro para invoice closed",
			status:  InvoiceStatusClosed,
			wantErr: ErrInvoiceAlreadyClosed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			inv := &Invoice{
				Status: tt.status,
			}

			err := inv.MarkAsPaid(now)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MarkAsPaid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if inv.Status != InvoiceStatusPaid {
					t.Errorf("Status = %q, want %q", inv.Status, InvoiceStatusPaid)
				}
				if inv.PaidAt == nil {
					t.Error("PaidAt não foi definido")
				}
			}
		})
	}
}
