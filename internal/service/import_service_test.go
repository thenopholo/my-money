package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/thenopholo/my-money/internal/domain"
)

func TestImportService_Confirm_BankAccount(t *testing.T) {
	userID := uuid.New()
	accountID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name               string
		req                domain.ImportConfirmRequest
		existingTxs        []*domain.Transaction
		wantCreated        int
		wantDuplicates     int
		wantCatCreated     int
		wantErr            bool
	}{
		{
			name: "cria transacao com categoria existente",
			req: domain.ImportConfirmRequest{
				ImportType: domain.ImportTypeBank,
				TargetID:   accountID,
				Transactions: []domain.ConfirmTransaction{
					{
						Description:     "PIX Recebido",
						Amount:          decimal.NewFromFloat(1000.00),
						TransactionDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
						TransactionType: domain.TransactionTypeIncome,
						CategoryID:      &categoryID,
					},
				},
			},
			existingTxs:    []*domain.Transaction{},
			wantCreated:    1,
			wantDuplicates: 0,
			wantCatCreated: 0,
		},
		{
			name: "detecta duplicata e pula",
			req: domain.ImportConfirmRequest{
				ImportType: domain.ImportTypeBank,
				TargetID:   accountID,
				Transactions: []domain.ConfirmTransaction{
					{
						Description:     "PIX Recebido",
						Amount:          decimal.NewFromFloat(1000.00),
						TransactionDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
						TransactionType: domain.TransactionTypeIncome,
						CategoryID:      &categoryID,
					},
				},
			},
			existingTxs: []*domain.Transaction{
				{
					ID:              uuid.New(),
					AccountID:       accountID,
					Description:     "PIX Recebido",
					Amount:          decimal.NewFromFloat(1000.00),
					TransactionDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
					TransactionType: domain.TransactionTypeIncome,
				},
			},
			wantCreated:    0,
			wantDuplicates: 1,
			wantCatCreated: 0,
		},
		{
			name: "cria nova categoria e transacao",
			req: domain.ImportConfirmRequest{
				ImportType: domain.ImportTypeBank,
				TargetID:   accountID,
				Transactions: []domain.ConfirmTransaction{
					{
						Description:     "Supermercado Carrefour",
						Amount:          decimal.NewFromFloat(245.67),
						TransactionDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
						TransactionType: domain.TransactionTypeExpense,
						NewCategoryName: strPtr("Alimentação"),
						NewCategoryType: catTypePtr(domain.CategoryTypeExpense),
					},
				},
			},
			existingTxs:    []*domain.Transaction{},
			wantCreated:    1,
			wantDuplicates: 0,
			wantCatCreated: 1,
		},
		{
			name:    "import type invalido",
			req:     domain.ImportConfirmRequest{ImportType: "invalid"},
			wantErr: true,
		},
		{
			name: "sem transacoes",
			req: domain.ImportConfirmRequest{
				ImportType:   domain.ImportTypeBank,
				TargetID:     accountID,
				Transactions: []domain.ConfirmTransaction{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txRepo := &mockTransactionRepository{
				getByAccountFn: func(_ context.Context, _ uuid.UUID) ([]*domain.Transaction, error) {
					return tt.existingTxs, nil
				},
				createFn: func(_ context.Context, _ *domain.Transaction) error {
					return nil
				},
			}

			ccTxRepo := &mockCreditCardTransactionRepository{
				getByCardFn: func(_ context.Context, _ uuid.UUID) ([]*domain.CreditCardTransaction, error) {
					return nil, nil
				},
				createFn: func(_ context.Context, _ *domain.CreditCardTransaction) error {
					return nil
				},
			}

			catRepo := &mockCategoryRepository{
				createFn: func(_ context.Context, _ *domain.Category) error {
					return nil
				},
				getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.Category, error) {
					return nil, nil
				},
			}

			account := &domain.BankAccount{
				ID:          accountID,
				UserID:      userID,
				AccountType: domain.AccountTypeChecking,
				Balance:     decimal.NewFromFloat(5000.00),
			}

			bankRepo := &mockBankAccountRepository{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.BankAccount, error) {
					return account, nil
				},
				updateFn: func(_ context.Context, _ *domain.BankAccount) error {
					return nil
				},
			}

			ccRepo := &mockCreditCardRepository{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCard, error) {
					return nil, domain.ErrCreditCardNotFound
				},
			}

			svc := NewImportService(
				NewLLMClient("http://localhost:8001"),
				txRepo,
				ccTxRepo,
				catRepo,
				bankRepo,
				ccRepo,
			)

			result, err := svc.Confirm(context.Background(), userID, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Created != tt.wantCreated {
				t.Errorf("created = %d, want %d", result.Created, tt.wantCreated)
			}
			if result.DuplicatesSkipped != tt.wantDuplicates {
				t.Errorf("duplicates = %d, want %d", result.DuplicatesSkipped, tt.wantDuplicates)
			}
			if result.CategoriesCreated != tt.wantCatCreated {
				t.Errorf("categories_created = %d, want %d", result.CategoriesCreated, tt.wantCatCreated)
			}
		})
	}
}

func TestIsWithinOneDay(t *testing.T) {
	base := time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		a, b time.Time
		want bool
	}{
		{"same time", base, base, true},
		{"12 hours apart", base, base.Add(12 * time.Hour), true},
		{"24 hours apart", base, base.Add(24 * time.Hour), true},
		{"25 hours apart", base, base.Add(25 * time.Hour), false},
		{"1 day before", base, base.Add(-24 * time.Hour), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isWithinOneDay(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("isWithinOneDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helpers para ponteiros
func strPtr(s string) *string {
	return &s
}

func catTypePtr(ct domain.CategoryType) *domain.CategoryType {
	return &ct
}
