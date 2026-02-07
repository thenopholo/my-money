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

func TestCreditCardTransactionService_Create(t *testing.T) {
	t.Parallel()

	cardID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name         string
		amount       decimal.Decimal
		description  string
		installments int
		mockCard     func() *mockCreditCardRepository
		mockCat      func() *mockCategoryRepository
		mockCCTx     func() *mockCreditCardTransactionRepository
		mockInvoice  func() *mockInvoiceRepository
		wantErr      bool
	}{
		{
			name:         "deve criar transação de cartão com sucesso",
			amount:       decimal.NewFromInt(100),
			description:  "Compra",
			installments: 1,
			mockCard: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCard, error) {
						return &domain.CreditCard{
							ID:          cardID,
							IsActive:    true,
							CreditLimit: decimal.NewFromInt(5000),
						}, nil
					},
				}
			},
			mockCat: func() *mockCategoryRepository {
				return &mockCategoryRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Category, error) {
						return &domain.Category{CategoryType: domain.CategoryTypeExpense}, nil
					},
				}
			},
			mockCCTx: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{
					getPendingFn: func(_ context.Context, _ uuid.UUID) ([]*domain.CreditCardTransaction, error) {
						return []*domain.CreditCardTransaction{}, nil
					},
					createFn: func(_ context.Context, _ *domain.CreditCardTransaction) error {
						return nil
					},
				}
			},
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{
					getByCardFn: func(_ context.Context, _ uuid.UUID) ([]*domain.Invoice, error) {
						return []*domain.Invoice{}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:         "deve retornar erro para cartão inativo",
			amount:       decimal.NewFromInt(100),
			description:  "Compra",
			installments: 1,
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
			mockCat: func() *mockCategoryRepository {
				return &mockCategoryRepository{}
			},
			mockCCTx: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{}
			},
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{}
			},
			wantErr: true,
		},
		{
			name:         "deve retornar erro para categoria income",
			amount:       decimal.NewFromInt(100),
			description:  "Compra",
			installments: 1,
			mockCard: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCard, error) {
						return &domain.CreditCard{
							ID:          cardID,
							IsActive:    true,
							CreditLimit: decimal.NewFromInt(5000),
						}, nil
					},
				}
			},
			mockCat: func() *mockCategoryRepository {
				return &mockCategoryRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Category, error) {
						return &domain.Category{CategoryType: domain.CategoryTypeIncome}, nil
					},
				}
			},
			mockCCTx: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{}
			},
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{}
			},
			wantErr: true,
		},
		{
			name:         "deve retornar erro para limite excedido",
			amount:       decimal.NewFromInt(6000),
			description:  "Compra grande",
			installments: 1,
			mockCard: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCard, error) {
						return &domain.CreditCard{
							ID:          cardID,
							IsActive:    true,
							CreditLimit: decimal.NewFromInt(5000),
						}, nil
					},
				}
			},
			mockCat: func() *mockCategoryRepository {
				return &mockCategoryRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Category, error) {
						return &domain.Category{CategoryType: domain.CategoryTypeExpense}, nil
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
				return &mockInvoiceRepository{
					getByCardFn: func(_ context.Context, _ uuid.UUID) ([]*domain.Invoice, error) {
						return []*domain.Invoice{}, nil
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewCreditCardTransactionService(tt.mockCCTx(), tt.mockCard(), tt.mockCat(), tt.mockInvoice())
			tx, err := svc.Create(context.Background(), cardID, categoryID, tt.amount, tt.description, tt.installments, now)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tx == nil {
				t.Error("Create() retornou nil sem erro")
			}
		})
	}
}

func TestCreditCardTransactionService_AssignToInvoice(t *testing.T) {
	t.Parallel()

	txID := uuid.New()
	invoiceID := uuid.New()
	cardID := uuid.New()

	tests := []struct {
		name        string
		mockCCTx    func() *mockCreditCardTransactionRepository
		mockInvoice func() *mockInvoiceRepository
		wantErr     error
	}{
		{
			name: "deve atribuir transação à fatura com sucesso",
			mockCCTx: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCardTransaction, error) {
						return &domain.CreditCardTransaction{ID: txID, CardID: cardID}, nil
					},
					assignToInvFn: func(_ context.Context, _, _ uuid.UUID) error {
						return nil
					},
				}
			},
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Invoice, error) {
						return &domain.Invoice{
							ID:     invoiceID,
							CardID: cardID,
							Status: domain.InvoiceStatusOpen,
						}, nil
					},
				}
			},
			wantErr: nil,
		},
		{
			name: "deve retornar erro para fatura não aberta",
			mockCCTx: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCardTransaction, error) {
						return &domain.CreditCardTransaction{ID: txID, CardID: cardID}, nil
					},
				}
			},
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Invoice, error) {
						return &domain.Invoice{
							ID:     invoiceID,
							CardID: cardID,
							Status: domain.InvoiceStatusPaid,
						}, nil
					},
				}
			},
			wantErr: domain.ErrInvoiceNotOpen,
		},
		{
			name: "deve retornar erro para card mismatch",
			mockCCTx: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCardTransaction, error) {
						return &domain.CreditCardTransaction{ID: txID, CardID: cardID}, nil
					},
				}
			},
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Invoice, error) {
						return &domain.Invoice{
							ID:     invoiceID,
							CardID: uuid.New(), // outro cartão
							Status: domain.InvoiceStatusOpen,
						}, nil
					},
				}
			},
			wantErr: domain.ErrInvoiceCardMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewCreditCardTransactionService(
				tt.mockCCTx(),
				&mockCreditCardRepository{},
				&mockCategoryRepository{},
				tt.mockInvoice(),
			)
			err := svc.AssignToInvoice(context.Background(), txID, invoiceID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AssignToInvoice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
