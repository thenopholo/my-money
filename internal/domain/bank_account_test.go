package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestAccountType_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		at   AccountType
		want bool
	}{
		{"checking é válido", AccountTypeChecking, true},
		{"savings é válido", AccountTypeSavings, true},
		{"string vazia é inválido", AccountType(""), false},
		{"tipo desconhecido é inválido", AccountType("investment"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.at.IsValid(); got != tt.want {
				t.Errorf("AccountType(%q).IsValid() = %v, want %v", tt.at, got, tt.want)
			}
		})
	}
}

func TestNewBankAccount_Validations(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	tests := []struct {
		name        string
		userID      uuid.UUID
		accountType AccountType
		accountName string
		bankName    string
		balance     decimal.Decimal
		wantErr     error
	}{
		{
			name:        "deve criar conta checking com saldo positivo",
			userID:      userID,
			accountType: AccountTypeChecking,
			accountName: "Conta Principal",
			bankName:    "Nubank",
			balance:     decimal.NewFromInt(1000),
			wantErr:     nil,
		},
		{
			name:        "deve criar conta checking com saldo negativo (permite overdraft)",
			userID:      userID,
			accountType: AccountTypeChecking,
			accountName: "Conta Principal",
			bankName:    "Nubank",
			balance:     decimal.NewFromInt(-500),
			wantErr:     nil,
		},
		{
			name:        "deve criar conta savings com saldo zero",
			userID:      userID,
			accountType: AccountTypeSavings,
			accountName: "Poupança",
			bankName:    "Inter",
			balance:     decimal.Zero,
			wantErr:     nil,
		},
		{
			name:        "deve retornar erro para tipo de conta inválido",
			userID:      userID,
			accountType: AccountType("investment"),
			accountName: "Conta",
			bankName:    "Banco",
			balance:     decimal.Zero,
			wantErr:     ErrInvalidAccountType,
		},
		{
			name:        "deve retornar erro para nome do banco vazio",
			userID:      userID,
			accountType: AccountTypeChecking,
			accountName: "Conta",
			bankName:    "",
			balance:     decimal.Zero,
			wantErr:     ErrEmptyBankName,
		},
		{
			name:        "deve retornar erro para saldo negativo em poupança",
			userID:      userID,
			accountType: AccountTypeSavings,
			accountName: "Poupança",
			bankName:    "Inter",
			balance:     decimal.NewFromInt(-100),
			wantErr:     ErrNegativeBalance,
		},
		{
			name:        "deve usar bankName como name quando name é vazio",
			userID:      userID,
			accountType: AccountTypeChecking,
			accountName: "",
			bankName:    "Nubank",
			balance:     decimal.Zero,
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ba, err := NewBankAccount(tt.userID, tt.accountType, tt.accountName, tt.bankName, tt.balance)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewBankAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if ba == nil {
					t.Fatal("NewBankAccount() retornou nil sem erro")
				}
				if ba.UserID != tt.userID {
					t.Errorf("UserID = %v, want %v", ba.UserID, tt.userID)
				}
				if ba.BankName != tt.bankName {
					t.Errorf("BankName = %q, want %q", ba.BankName, tt.bankName)
				}
				if !ba.Balance.Equal(tt.balance) {
					t.Errorf("Balance = %v, want %v", ba.Balance, tt.balance)
				}
				if !ba.IsActive {
					t.Error("IsActive deve ser true para conta nova")
				}
				if tt.accountName == "" && ba.Name != tt.bankName {
					t.Errorf("Name = %q, want bankName %q quando nome é vazio", ba.Name, tt.bankName)
				}
			}
		})
	}
}

func TestBankAccount_ApplyIncome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		initial     decimal.Decimal
		amount      decimal.Decimal
		wantBalance decimal.Decimal
		wantErr     error
	}{
		{
			name:        "deve somar receita ao saldo",
			initial:     decimal.NewFromInt(1000),
			amount:      decimal.NewFromInt(500),
			wantBalance: decimal.NewFromInt(1500),
			wantErr:     nil,
		},
		{
			name:    "deve retornar erro para valor zero",
			initial: decimal.NewFromInt(1000),
			amount:  decimal.Zero,
			wantErr: ErrInvalidAmount,
		},
		{
			name:    "deve retornar erro para valor negativo",
			initial: decimal.NewFromInt(1000),
			amount:  decimal.NewFromInt(-100),
			wantErr: ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ba := &BankAccount{
				AccountType: AccountTypeChecking,
				Balance:     tt.initial,
			}

			err := ba.ApplyIncome(tt.amount)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApplyIncome() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && !ba.Balance.Equal(tt.wantBalance) {
				t.Errorf("Balance = %v, want %v", ba.Balance, tt.wantBalance)
			}
		})
	}
}

func TestBankAccount_ApplyExpense(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accountType AccountType
		initial     decimal.Decimal
		amount      decimal.Decimal
		wantBalance decimal.Decimal
		wantErr     error
	}{
		{
			name:        "deve subtrair despesa do saldo em conta corrente",
			accountType: AccountTypeChecking,
			initial:     decimal.NewFromInt(1000),
			amount:      decimal.NewFromInt(500),
			wantBalance: decimal.NewFromInt(500),
			wantErr:     nil,
		},
		{
			name:        "deve permitir saldo negativo em conta corrente",
			accountType: AccountTypeChecking,
			initial:     decimal.NewFromInt(100),
			amount:      decimal.NewFromInt(500),
			wantBalance: decimal.NewFromInt(-400),
			wantErr:     nil,
		},
		{
			name:        "deve retornar erro para saldo negativo em poupança",
			accountType: AccountTypeSavings,
			initial:     decimal.NewFromInt(100),
			amount:      decimal.NewFromInt(500),
			wantErr:     ErrNegativeBalance,
		},
		{
			name:        "deve retornar erro para valor zero",
			accountType: AccountTypeChecking,
			initial:     decimal.NewFromInt(1000),
			amount:      decimal.Zero,
			wantErr:     ErrInvalidAmount,
		},
		{
			name:        "deve retornar erro para valor negativo",
			accountType: AccountTypeChecking,
			initial:     decimal.NewFromInt(1000),
			amount:      decimal.NewFromInt(-100),
			wantErr:     ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ba := &BankAccount{
				AccountType: tt.accountType,
				Balance:     tt.initial,
			}

			err := ba.ApplyExpense(tt.amount)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApplyExpense() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && !ba.Balance.Equal(tt.wantBalance) {
				t.Errorf("Balance = %v, want %v", ba.Balance, tt.wantBalance)
			}
		})
	}
}

func TestBankAccount_AllowsOverdraft(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accountType AccountType
		want        bool
	}{
		{"checking permite overdraft", AccountTypeChecking, true},
		{"savings não permite overdraft", AccountTypeSavings, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ba := &BankAccount{AccountType: tt.accountType}
			if got := ba.AllowsOverdraft(); got != tt.want {
				t.Errorf("AllowsOverdraft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankAccount_SetBankName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"deve aceitar nome válido", "Nubank", nil},
		{"deve retornar erro para nome vazio", "", ErrEmptyBankName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ba := &BankAccount{}
			err := ba.SetBankName(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SetBankName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && ba.BankName != tt.input {
				t.Errorf("BankName = %q, want %q", ba.BankName, tt.input)
			}
		})
	}
}

func TestBankAccount_SetBalance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accountType AccountType
		balance     decimal.Decimal
		wantErr     error
	}{
		{"checking aceita saldo negativo", AccountTypeChecking, decimal.NewFromInt(-100), nil},
		{"savings rejeita saldo negativo", AccountTypeSavings, decimal.NewFromInt(-100), ErrNegativeBalance},
		{"savings aceita saldo zero", AccountTypeSavings, decimal.Zero, nil},
		{"savings aceita saldo positivo", AccountTypeSavings, decimal.NewFromInt(100), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ba := &BankAccount{AccountType: tt.accountType}
			err := ba.SetBalance(tt.balance)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SetBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
