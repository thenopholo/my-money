package service

import (
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/thenopholo/my-money/internal/domain"
)

func TestParseBankCSV(t *testing.T) {
	tests := []struct {
		name        string
		csv         string
		wantLen     int
		wantFirst   domain.RawCSVTransaction
		wantErr     bool
		wantErrType error
	}{
		{
			name: "formato BR com separador ponto-e-virgula",
			csv: `Data;Descrição;Valor
01/02/2026;PIX RECEBIDO - EMPRESA XYZ;3500.00
01/02/2026;SUPERMERCADO CARREFOUR;-245.67
02/02/2026;NETFLIX;-55.90`,
			wantLen: 3,
			wantFirst: domain.RawCSVTransaction{
				Description:     "PIX RECEBIDO - EMPRESA XYZ",
				Amount:          decimal.NewFromFloat(3500.00),
				TransactionDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
				TransactionType: domain.TransactionTypeIncome,
			},
		},
		{
			name: "formato com virgula como separador",
			csv: `Data,Descrição,Valor
01/02/2026,PIX RECEBIDO,3500.00
02/02/2026,MERCADO,-123.45`,
			wantLen: 2,
			wantFirst: domain.RawCSVTransaction{
				Description:     "PIX RECEBIDO",
				Amount:          decimal.NewFromFloat(3500.00),
				TransactionDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
				TransactionType: domain.TransactionTypeIncome,
			},
		},
		{
			name: "formato de data YYYY-MM-DD",
			csv: `Date;Description;Amount
2026-02-01;SALARY;5000.00`,
			wantLen: 1,
			wantFirst: domain.RawCSVTransaction{
				Description:     "SALARY",
				Amount:          decimal.NewFromFloat(5000.00),
				TransactionDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
				TransactionType: domain.TransactionTypeIncome,
			},
		},
		{
			name: "formato BR de valor com milhar",
			csv: `Data;Descrição;Valor
01/02/2026;SALARIO;3.500,00
01/02/2026;ALUGUEL;-1.200,50`,
			wantLen: 2,
			wantFirst: domain.RawCSVTransaction{
				Description:     "SALARIO",
				Amount:          decimal.RequireFromString("3500.00"),
				TransactionDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
				TransactionType: domain.TransactionTypeIncome,
			},
		},
		{
			name:        "CSV vazio",
			csv:         "",
			wantErr:     true,
			wantErrType: domain.ErrEmptyCSV,
		},
		{
			name:        "apenas cabecalho",
			csv:         "Data;Descrição;Valor\n",
			wantErr:     true,
			wantErrType: domain.ErrEmptyCSV,
		},
		{
			name: "sem cabecalho",
			csv: `01/02/2026;PIX RECEBIDO;1500.00
02/02/2026;MERCADO;-200.00`,
			wantLen: 2,
			wantFirst: domain.RawCSVTransaction{
				Description:     "PIX RECEBIDO",
				Amount:          decimal.NewFromFloat(1500.00),
				TransactionDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
				TransactionType: domain.TransactionTypeIncome,
			},
		},
		{
			name:        "colunas insuficientes",
			csv:         "01/02/2026;PIX RECEBIDO\n",
			wantErr:     true,
			wantErrType: domain.ErrInvalidCSVFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactions, err := ParseBankCSV(strings.NewReader(tt.csv))

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(transactions) != tt.wantLen {
				t.Errorf("got %d transactions, want %d", len(transactions), tt.wantLen)
			}

			if tt.wantLen > 0 {
				got := transactions[0]
				if got.Description != tt.wantFirst.Description {
					t.Errorf("description = %q, want %q", got.Description, tt.wantFirst.Description)
				}
				if !got.Amount.Equal(tt.wantFirst.Amount) {
					t.Errorf("amount = %s, want %s", got.Amount, tt.wantFirst.Amount)
				}
				if !got.TransactionDate.Equal(tt.wantFirst.TransactionDate) {
					t.Errorf("date = %v, want %v", got.TransactionDate, tt.wantFirst.TransactionDate)
				}
				if got.TransactionType != tt.wantFirst.TransactionType {
					t.Errorf("type = %q, want %q", got.TransactionType, tt.wantFirst.TransactionType)
				}
			}
		})
	}
}

func TestParseCreditCardCSV(t *testing.T) {
	tests := []struct {
		name        string
		csv         string
		wantLen     int
		wantFirst   domain.RawCSVTransaction
		wantErr     bool
		wantErrType error
	}{
		{
			name: "fatura com parcelas na coluna",
			csv: `Data;Descrição;Valor;Parcela
15/01/2026;AMAZON - FONE BLUETOOTH;299.90;1/3
15/01/2026;SPOTIFY PREMIUM;34.90;
18/01/2026;RESTAURANTE OUTBACK;187.50;`,
			wantLen: 3,
			wantFirst: domain.RawCSVTransaction{
				Description:     "AMAZON - FONE BLUETOOTH",
				Amount:          decimal.NewFromFloat(299.90),
				TransactionDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
				TransactionType: domain.TransactionTypeExpense,
				Installments:    intPtr(3),
				CurrentInstall:  intPtr(1),
			},
		},
		{
			name: "fatura sem coluna de parcela",
			csv: `Data;Descrição;Valor
15/01/2026;SPOTIFY PREMIUM;34.90
18/01/2026;RESTAURANTE OUTBACK;187.50`,
			wantLen: 2,
			wantFirst: domain.RawCSVTransaction{
				Description:     "SPOTIFY PREMIUM",
				Amount:          decimal.NewFromFloat(34.90),
				TransactionDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
				TransactionType: domain.TransactionTypeExpense,
			},
		},
		{
			name: "parcela extraida da descricao",
			csv: `Data;Descrição;Valor
15/01/2026;AMAZON FONE 2/6;299.90`,
			wantLen: 1,
			wantFirst: domain.RawCSVTransaction{
				Description:     "AMAZON FONE 2/6",
				Amount:          decimal.NewFromFloat(299.90),
				TransactionDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
				TransactionType: domain.TransactionTypeExpense,
				Installments:    intPtr(6),
				CurrentInstall:  intPtr(2),
			},
		},
		{
			name:    "CSV vazio",
			csv:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactions, err := ParseCreditCardCSV(strings.NewReader(tt.csv))

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(transactions) != tt.wantLen {
				t.Errorf("got %d transactions, want %d", len(transactions), tt.wantLen)
			}

			if tt.wantLen > 0 {
				got := transactions[0]
				if got.Description != tt.wantFirst.Description {
					t.Errorf("description = %q, want %q", got.Description, tt.wantFirst.Description)
				}
				if !got.Amount.Equal(tt.wantFirst.Amount) {
					t.Errorf("amount = %s, want %s", got.Amount, tt.wantFirst.Amount)
				}
				if !got.TransactionDate.Equal(tt.wantFirst.TransactionDate) {
					t.Errorf("date = %v, want %v", got.TransactionDate, tt.wantFirst.TransactionDate)
				}
				if got.TransactionType != tt.wantFirst.TransactionType {
					t.Errorf("type = %q, want %q", got.TransactionType, tt.wantFirst.TransactionType)
				}
				if tt.wantFirst.Installments != nil {
					if got.Installments == nil || *got.Installments != *tt.wantFirst.Installments {
						t.Errorf("installments = %v, want %v", got.Installments, tt.wantFirst.Installments)
					}
				}
				if tt.wantFirst.CurrentInstall != nil {
					if got.CurrentInstall == nil || *got.CurrentInstall != *tt.wantFirst.CurrentInstall {
						t.Errorf("current_installment = %v, want %v", got.CurrentInstall, tt.wantFirst.CurrentInstall)
					}
				}
			}
		})
	}
}

func TestParseCSVAmount(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"US format", "1234.56", "1234.56", false},
		{"BR format", "1.234,56", "1234.56", false},
		{"negative US", "-245.67", "-245.67", false},
		{"negative BR", "-1.200,50", "-1200.50", false},
		{"with currency", "R$ 1.500,00", "1500.00", false},
		{"simple comma", "245,67", "245.67", false},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCSVAmount(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := decimal.RequireFromString(tt.want)
			if !got.Equal(want) {
				t.Errorf("parseCSVAmount(%q) = %s, want %s", tt.input, got, want)
			}
		})
	}
}

func TestParseCSVDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{"DD/MM/YYYY", "01/02/2026", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), false},
		{"YYYY-MM-DD", "2026-02-01", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), false},
		{"DD/MM/YY", "01/02/26", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), false},
		{"invalid", "abc", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCSVDate(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("parseCSVDate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDetectSeparator(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    rune
	}{
		{"semicolon", "Data;Descrição;Valor\n01/02/2026;PIX;100", ';'},
		{"comma", "Data,Descrição,Valor\n01/02/2026,PIX,100", ','},
		{"more semicolons", "a;b;c;d\n1;2;3;4", ';'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectSeparator(tt.content)
			if got != tt.want {
				t.Errorf("detectSeparator() = %c, want %c", got, tt.want)
			}
		})
	}
}

func TestParseInstallments(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantCurrent int
		wantTotal   int
		wantOk      bool
	}{
		{"valid 1/3", "1/3", 1, 3, true},
		{"valid 02/12", "02/12", 2, 12, true},
		{"empty", "", 0, 0, false},
		{"invalid current > total", "5/3", 0, 0, false},
		{"no match", "abc", 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current, total, ok := parseInstallments(tt.input)
			if ok != tt.wantOk {
				t.Errorf("ok = %v, want %v", ok, tt.wantOk)
			}
			if ok {
				if current != tt.wantCurrent {
					t.Errorf("current = %d, want %d", current, tt.wantCurrent)
				}
				if total != tt.wantTotal {
					t.Errorf("total = %d, want %d", total, tt.wantTotal)
				}
			}
		})
	}
}

func intPtr(v int) *int {
	return &v
}
