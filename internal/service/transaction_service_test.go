package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

func TestTransactionService_Create(t *testing.T) {
	t.Parallel()

	accountID := uuid.New()
	categoryID := uuid.New()
	yesterday := time.Now().AddDate(0, 0, -1)

	tests := []struct {
		name     string
		amount   decimal.Decimal
		txType   domain.TransactionType
		mockAcct func() *mockBankAccountRepository
		mockCat  func() *mockCategoryRepository
		mockTx   func() *mockTransactionRepository
		wantErr  bool
	}{
		{
			name:   "deve criar transação income com sucesso",
			amount: decimal.NewFromInt(500),
			txType: domain.TransactionTypeIncome,
			mockAcct: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.BankAccount, error) {
						return &domain.BankAccount{
							ID:          accountID,
							AccountType: domain.AccountTypeChecking,
							Balance:     decimal.NewFromInt(1000),
						}, nil
					},
					updateFn: func(_ context.Context, _ *domain.BankAccount) error { return nil },
				}
			},
			mockCat: func() *mockCategoryRepository {
				return &mockCategoryRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Category, error) {
						return &domain.Category{
							ID:           categoryID,
							CategoryType: domain.CategoryTypeIncome,
						}, nil
					},
				}
			},
			mockTx: func() *mockTransactionRepository {
				return &mockTransactionRepository{
					createFn: func(_ context.Context, _ *domain.Transaction) error { return nil },
				}
			},
			wantErr: false,
		},
		{
			name:   "deve retornar erro para category mismatch",
			amount: decimal.NewFromInt(500),
			txType: domain.TransactionTypeIncome,
			mockAcct: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.BankAccount, error) {
						return &domain.BankAccount{
							ID:          accountID,
							AccountType: domain.AccountTypeChecking,
							Balance:     decimal.NewFromInt(1000),
						}, nil
					},
				}
			},
			mockCat: func() *mockCategoryRepository {
				return &mockCategoryRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Category, error) {
						return &domain.Category{
							ID:           categoryID,
							CategoryType: domain.CategoryTypeExpense, // mismatch!
						}, nil
					},
				}
			},
			mockTx: func() *mockTransactionRepository {
				return &mockTransactionRepository{}
			},
			wantErr: true,
		},
		{
			name:   "deve retornar erro para conta não encontrada",
			amount: decimal.NewFromInt(500),
			txType: domain.TransactionTypeIncome,
			mockAcct: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.BankAccount, error) {
						return nil, domain.ErrBankAccountNotFound
					},
				}
			},
			mockCat: func() *mockCategoryRepository {
				return &mockCategoryRepository{}
			},
			mockTx: func() *mockTransactionRepository {
				return &mockTransactionRepository{}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewTransactionService(
				tt.mockTx(),
				tt.mockAcct(),
				tt.mockCat(),
				&mockPlannedIncomeRepository{},
				&mockPlannedExpenseRepository{},
				&mockInvoiceRepository{},
			)
			tx, err := svc.Create(context.Background(), accountID, categoryID, tt.amount, tt.txType, "Teste", yesterday)

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

func TestTransactionService_Delete(t *testing.T) {
	t.Parallel()

	txID := uuid.New()
	accountID := uuid.New()

	tests := []struct {
		name     string
		mockAcct func() *mockBankAccountRepository
		mockTx   func() *mockTransactionRepository
		wantErr  bool
	}{
		{
			name: "deve deletar transação income e reverter saldo",
			mockAcct: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.BankAccount, error) {
						return &domain.BankAccount{
							ID:          accountID,
							AccountType: domain.AccountTypeChecking,
							Balance:     decimal.NewFromInt(1500),
						}, nil
					},
					updateFn: func(_ context.Context, _ *domain.BankAccount) error { return nil },
				}
			},
			mockTx: func() *mockTransactionRepository {
				return &mockTransactionRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Transaction, error) {
						return &domain.Transaction{
							ID:              txID,
							AccountID:       accountID,
							Amount:          decimal.NewFromInt(500),
							TransactionType: domain.TransactionTypeIncome,
						}, nil
					},
					deleteFn: func(_ context.Context, _ uuid.UUID) error { return nil },
				}
			},
			wantErr: false,
		},
		{
			name: "deve retornar erro para transação não encontrada",
			mockAcct: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{}
			},
			mockTx: func() *mockTransactionRepository {
				return &mockTransactionRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Transaction, error) {
						return nil, domain.ErrTransactionNotFound
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewTransactionService(
				tt.mockTx(),
				tt.mockAcct(),
				&mockCategoryRepository{},
				&mockPlannedIncomeRepository{},
				&mockPlannedExpenseRepository{},
				&mockInvoiceRepository{},
			)
			err := svc.Delete(context.Background(), txID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionService_PayInvoice(t *testing.T) {
	t.Parallel()

	invoiceID := uuid.New()
	accountID := uuid.New()
	categoryID := uuid.New()
	payDate := time.Now()

	tests := []struct {
		name        string
		mockAcct    func() *mockBankAccountRepository
		mockCat     func() *mockCategoryRepository
		mockTx      func() *mockTransactionRepository
		mockInvoice func() *mockInvoiceRepository
		wantErr     bool
	}{
		{
			name: "deve pagar fatura com sucesso",
			mockAcct: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.BankAccount, error) {
						return &domain.BankAccount{
							ID:          accountID,
							AccountType: domain.AccountTypeChecking,
							Balance:     decimal.NewFromInt(5000),
						}, nil
					},
					updateFn: func(_ context.Context, _ *domain.BankAccount) error { return nil },
				}
			},
			mockCat: func() *mockCategoryRepository {
				return &mockCategoryRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Category, error) {
						return &domain.Category{
							CategoryType: domain.CategoryTypeExpense,
						}, nil
					},
				}
			},
			mockTx: func() *mockTransactionRepository {
				return &mockTransactionRepository{
					createFn: func(_ context.Context, _ *domain.Transaction) error { return nil },
				}
			},
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Invoice, error) {
						return &domain.Invoice{
							ID:          invoiceID,
							Status:      domain.InvoiceStatusOpen,
							TotalAmount: decimal.NewFromInt(1000),
						}, nil
					},
					updateFn: func(_ context.Context, _ *domain.Invoice) error { return nil },
				}
			},
			wantErr: false,
		},
		{
			name: "deve retornar erro para fatura já paga",
			mockAcct: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{}
			},
			mockCat: func() *mockCategoryRepository {
				return &mockCategoryRepository{}
			},
			mockTx: func() *mockTransactionRepository {
				return &mockTransactionRepository{}
			},
			mockInvoice: func() *mockInvoiceRepository {
				return &mockInvoiceRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Invoice, error) {
						return &domain.Invoice{
							ID:     invoiceID,
							Status: domain.InvoiceStatusPaid,
						}, nil
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewTransactionService(
				tt.mockTx(),
				tt.mockAcct(),
				tt.mockCat(),
				&mockPlannedIncomeRepository{},
				&mockPlannedExpenseRepository{},
				tt.mockInvoice(),
			)
			tx, err := svc.PayInvoice(context.Background(), invoiceID, accountID, categoryID, payDate, "Pagamento fatura")

			if (err != nil) != tt.wantErr {
				t.Errorf("PayInvoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tx == nil {
				t.Error("PayInvoice() retornou nil sem erro")
			}
		})
	}
}
