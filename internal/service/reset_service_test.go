package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

func TestResetService_ResetAllTransactions(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	accountID1 := uuid.New()
	accountID2 := uuid.New()
	cardID1 := uuid.New()

	tests := []struct {
		name       string
		txRepo     func() *mockTransactionRepository
		ccTxRepo   func() *mockCreditCardTransactionRepository
		acctRepo   func() *mockBankAccountRepository
		ccRepo     func() *mockCreditCardRepository
		wantErr    bool
		errMessage string
	}{
		{
			name: "should reset all transactions successfully",
			txRepo: func() *mockTransactionRepository {
				return &mockTransactionRepository{
					deleteAllByAccountFn: func(_ context.Context, _ uuid.UUID) error {
						return nil
					},
				}
			},
			ccTxRepo: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{
					deleteAllByCardFn: func(_ context.Context, _ uuid.UUID) error {
						return nil
					},
				}
			},
			acctRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.BankAccount, error) {
						return []*domain.BankAccount{
							{ID: accountID1, UserID: userID, AccountType: domain.AccountTypeChecking, Balance: decimal.NewFromInt(1000)},
							{ID: accountID2, UserID: userID, AccountType: domain.AccountTypeChecking, Balance: decimal.NewFromInt(500)},
						}, nil
					},
					updateFn: func(_ context.Context, ba *domain.BankAccount) error {
						if !ba.Balance.IsZero() {
							t.Errorf("expected balance to be zero, got %s", ba.Balance.String())
						}
						return nil
					},
				}
			},
			ccRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.CreditCard, error) {
						return []*domain.CreditCard{
							{ID: cardID1, UserID: userID},
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "should succeed when user has no accounts and no cards",
			txRepo: func() *mockTransactionRepository {
				return &mockTransactionRepository{}
			},
			ccTxRepo: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{}
			},
			acctRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.BankAccount, error) {
						return []*domain.BankAccount{}, nil
					},
				}
			},
			ccRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.CreditCard, error) {
						return []*domain.CreditCard{}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "should return error when fetching accounts fails",
			txRepo: func() *mockTransactionRepository {
				return &mockTransactionRepository{}
			},
			ccTxRepo: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{}
			},
			acctRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.BankAccount, error) {
						return nil, errors.New("db error")
					},
				}
			},
			ccRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{}
			},
			wantErr: true,
		},
		{
			name: "should return error when deleting transactions fails",
			txRepo: func() *mockTransactionRepository {
				return &mockTransactionRepository{
					deleteAllByAccountFn: func(_ context.Context, _ uuid.UUID) error {
						return errors.New("delete error")
					},
				}
			},
			ccTxRepo: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{}
			},
			acctRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.BankAccount, error) {
						return []*domain.BankAccount{
							{ID: accountID1, UserID: userID, AccountType: domain.AccountTypeChecking},
						}, nil
					},
				}
			},
			ccRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{}
			},
			wantErr: true,
		},
		{
			name: "should return error when updating account balance fails",
			txRepo: func() *mockTransactionRepository {
				return &mockTransactionRepository{
					deleteAllByAccountFn: func(_ context.Context, _ uuid.UUID) error {
						return nil
					},
				}
			},
			ccTxRepo: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{}
			},
			acctRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.BankAccount, error) {
						return []*domain.BankAccount{
							{ID: accountID1, UserID: userID, AccountType: domain.AccountTypeChecking},
						}, nil
					},
					updateFn: func(_ context.Context, _ *domain.BankAccount) error {
						return errors.New("update error")
					},
				}
			},
			ccRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{}
			},
			wantErr: true,
		},
		{
			name: "should return error when fetching credit cards fails",
			txRepo: func() *mockTransactionRepository {
				return &mockTransactionRepository{}
			},
			ccTxRepo: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{}
			},
			acctRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.BankAccount, error) {
						return []*domain.BankAccount{}, nil
					},
				}
			},
			ccRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.CreditCard, error) {
						return nil, errors.New("db error")
					},
				}
			},
			wantErr: true,
		},
		{
			name: "should return error when deleting credit card transactions fails",
			txRepo: func() *mockTransactionRepository {
				return &mockTransactionRepository{}
			},
			ccTxRepo: func() *mockCreditCardTransactionRepository {
				return &mockCreditCardTransactionRepository{
					deleteAllByCardFn: func(_ context.Context, _ uuid.UUID) error {
						return errors.New("delete error")
					},
				}
			},
			acctRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.BankAccount, error) {
						return []*domain.BankAccount{}, nil
					},
				}
			},
			ccRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByUserFn: func(_ context.Context, _ uuid.UUID) ([]*domain.CreditCard, error) {
						return []*domain.CreditCard{
							{ID: cardID1, UserID: userID},
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

			svc := NewResetService(tt.txRepo(), tt.ccTxRepo(), tt.acctRepo(), tt.ccRepo())
			err := svc.ResetAllTransactions(context.Background(), userID)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
