package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

func TestBankAccountService_Create(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	tests := []struct {
		name        string
		accountType domain.AccountType
		accName     string
		bankName    string
		balance     decimal.Decimal
		mockRepo    func() *mockBankAccountRepository
		wantErr     bool
	}{
		{
			name:        "deve criar conta com sucesso",
			accountType: domain.AccountTypeChecking,
			accName:     "Principal",
			bankName:    "Nubank",
			balance:     decimal.NewFromInt(1000),
			mockRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					createFn: func(_ context.Context, _ *domain.BankAccount) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:        "deve retornar erro de domínio para tipo inválido",
			accountType: domain.AccountType("invalid"),
			accName:     "Principal",
			bankName:    "Nubank",
			balance:     decimal.Zero,
			mockRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{}
			},
			wantErr: true,
		},
		{
			name:        "deve retornar erro do repo",
			accountType: domain.AccountTypeChecking,
			accName:     "Principal",
			bankName:    "Nubank",
			balance:     decimal.Zero,
			mockRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					createFn: func(_ context.Context, _ *domain.BankAccount) error {
						return errors.New("db error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewBankAccountService(tt.mockRepo())
			account, err := svc.Create(context.Background(), userID, tt.accountType, tt.accName, tt.bankName, tt.balance)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && account == nil {
				t.Error("Create() retornou nil sem erro")
			}
		})
	}
}

func TestBankAccountService_GetByID(t *testing.T) {
	t.Parallel()

	accountID := uuid.New()

	tests := []struct {
		name     string
		mockRepo func() *mockBankAccountRepository
		wantErr  error
	}{
		{
			name: "deve retornar conta com sucesso",
			mockRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.BankAccount, error) {
						return &domain.BankAccount{ID: accountID}, nil
					},
				}
			},
			wantErr: nil,
		},
		{
			name: "deve retornar erro quando não encontrada",
			mockRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.BankAccount, error) {
						return nil, domain.ErrBankAccountNotFound
					},
				}
			},
			wantErr: domain.ErrBankAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewBankAccountService(tt.mockRepo())
			account, err := svc.GetByID(context.Background(), accountID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && account == nil {
				t.Error("GetByID() retornou nil sem erro")
			}
		})
	}
}

func TestBankAccountService_Update(t *testing.T) {
	t.Parallel()

	accountID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name     string
		bankName string
		mockRepo func() *mockBankAccountRepository
		wantErr  bool
	}{
		{
			name:     "deve atualizar com sucesso",
			bankName: "Inter",
			mockRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.BankAccount, error) {
						return &domain.BankAccount{
							ID:          accountID,
							UserID:      userID,
							AccountType: domain.AccountTypeChecking,
							BankName:    "Nubank",
							Balance:     decimal.NewFromInt(1000),
						}, nil
					},
					updateFn: func(_ context.Context, _ *domain.BankAccount) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:     "deve retornar erro para bankName vazio",
			bankName: "",
			mockRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.BankAccount, error) {
						return &domain.BankAccount{
							ID:          accountID,
							AccountType: domain.AccountTypeChecking,
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
			svc := NewBankAccountService(tt.mockRepo())
			_, err := svc.Update(context.Background(), accountID, "Conta", tt.bankName, decimal.NewFromInt(500), true)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBankAccountService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mockRepo func() *mockBankAccountRepository
		wantErr  bool
	}{
		{
			name: "deve deletar com sucesso",
			mockRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					deleteFn: func(_ context.Context, _ uuid.UUID) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "deve retornar erro do repo",
			mockRepo: func() *mockBankAccountRepository {
				return &mockBankAccountRepository{
					deleteFn: func(_ context.Context, _ uuid.UUID) error {
						return errors.New("db error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewBankAccountService(tt.mockRepo())
			err := svc.Delete(context.Background(), uuid.New())
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
