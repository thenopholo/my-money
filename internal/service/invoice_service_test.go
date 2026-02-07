package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

func TestInvoiceService_CloseMonthInvoice(t *testing.T) {
	t.Parallel()

	cardID := uuid.New()
	refDate := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		mockCard    func() *mockCreditCardRepository
		mockCCTx    func() *mockCreditCardTransactionRepository
		mockInvoice func() *mockInvoiceRepository
		wantErr     error
	}{
		{
			name: "deve fechar fatura com sucesso",
			mockCard: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCard, error) {
						return &domain.CreditCard{
							ID:       cardID,
							IsActive: true,
							DueDay:   25,
						}, nil
					},
				}
			},
			mockCCTx: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{
					getPendingFn: func(_ context.Context, _ uuid.UUID) ([]*domain.CreditCardTransaction, error) {
						return []*domain.CreditCardTransaction{
							{
								ID:              uuid.New(),
								CardID:          cardID,
								Amount:          decimal.NewFromInt(100),
								TransactionDate: refDate,
							},
							{
								ID:              uuid.New(),
								CardID:          cardID,
								Amount:          decimal.NewFromInt(200),
								TransactionDate: refDate,
							},
						}, nil
					},
					assignToInvFn: func(_ context.Context, _, _ uuid.UUID) error {
						return nil
					},
				}
			},
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{
					createFn: func(_ context.Context, _ *domain.Invoice) error {
						return nil
					},
				}
			},
			wantErr: nil,
		},
		{
			name: "deve retornar erro para cartão inativo",
			mockCard: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCard, error) {
						return &domain.CreditCard{
							ID:       cardID,
							IsActive: false,
						}, nil
					},
				}
			},
			mockCCTx: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{}
			},
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{}
			},
			wantErr: domain.ErrCreditCardInactive,
		},
		{
			name: "deve retornar erro sem transações pendentes",
			mockCard: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCard, error) {
						return &domain.CreditCard{
							ID:       cardID,
							IsActive: true,
							DueDay:   25,
						}, nil
					},
				}
			},
			mockCCTx: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{
					getPendingFn: func(_ context.Context, _ uuid.UUID) ([]*domain.CreditCardTransaction, error) {
						return []*domain.CreditCardTransaction{}, nil
					},
				}
			},
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{}
			},
			wantErr: domain.ErrNoPendingTransactions,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewInvoiceService(tt.mockInvoice(), tt.mockCard(), tt.mockCCTx())
			invoice, err := svc.CloseMonthInvoice(context.Background(), cardID, refDate)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CloseMonthInvoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if invoice == nil {
					t.Error("CloseMonthInvoice() retornou nil sem erro")
				}
				expected := decimal.NewFromInt(300)
				if !invoice.TotalAmount.Equal(expected) {
					t.Errorf("TotalAmount = %v, want %v", invoice.TotalAmount, expected)
				}
			}
		})
	}
}

func TestInvoiceService_GetByID(t *testing.T) {
	t.Parallel()

	invID := uuid.New()

	tests := []struct {
		name        string
		mockInvoice func() *mockInvoiceRepository
		wantErr     error
	}{
		{
			name: "deve retornar fatura com sucesso",
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Invoice, error) {
						return &domain.Invoice{ID: invID}, nil
					},
				}
			},
			wantErr: nil,
		},
		{
			name: "deve retornar erro quando não encontrada",
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Invoice, error) {
						return nil, domain.ErrInvoiceNotFound
					},
				}
			},
			wantErr: domain.ErrInvoiceNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewInvoiceService(tt.mockInvoice(), &mockCreditCardRepository{}, &mockCreditCardTransactionRepository{})
			_, err := svc.GetByID(context.Background(), invID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
