package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

func TestCreditCardService_Create(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	tests := []struct {
		name     string
		cardName string
		limit    decimal.Decimal
		closeDay int
		dueDay   int
		mockRepo func() *mockCreditCardRepository
		wantErr  bool
	}{
		{
			name:     "deve criar cartão com sucesso",
			cardName: "Nubank",
			limit:    decimal.NewFromInt(5000),
			closeDay: 15,
			dueDay:   25,
			mockRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					createFn: func(_ context.Context, _ *domain.CreditCard) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:     "deve retornar erro para nome vazio",
			cardName: "",
			limit:    decimal.NewFromInt(5000),
			closeDay: 15,
			dueDay:   25,
			mockRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{}
			},
			wantErr: true,
		},
		{
			name:     "deve retornar erro do repo",
			cardName: "Nubank",
			limit:    decimal.NewFromInt(5000),
			closeDay: 15,
			dueDay:   25,
			mockRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					createFn: func(_ context.Context, _ *domain.CreditCard) error {
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
			svc := NewCreditCardService(tt.mockRepo())
			card, err := svc.Create(context.Background(), userID, tt.cardName, tt.limit, tt.closeDay, tt.dueDay)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && card == nil {
				t.Error("Create() retornou nil sem erro")
			}
		})
	}
}

func TestCreditCardService_Update(t *testing.T) {
	t.Parallel()

	cardID := uuid.New()

	tests := []struct {
		name     string
		cardName string
		limit    decimal.Decimal
		closeDay int
		dueDay   int
		isActive bool
		mockRepo func() *mockCreditCardRepository
		wantErr  bool
	}{
		{
			name:     "deve atualizar com sucesso",
			cardName: "Inter",
			limit:    decimal.NewFromInt(10000),
			closeDay: 10,
			dueDay:   20,
			isActive: true,
			mockRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCard, error) {
						return &domain.CreditCard{
							ID:          cardID,
							Name:        "Nubank",
							CreditLimit: decimal.NewFromInt(5000),
							CloseDay:    15,
							DueDay:      25,
						}, nil
					},
					updateFn: func(_ context.Context, _ *domain.CreditCard) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:     "deve retornar erro quando não encontrado",
			cardName: "Inter",
			limit:    decimal.NewFromInt(10000),
			closeDay: 10,
			dueDay:   20,
			isActive: true,
			mockRepo: func() *mockCreditCardRepository {
				return &mockCreditCardRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.CreditCard, error) {
						return nil, domain.ErrCreditCardNotFound
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewCreditCardService(tt.mockRepo())
			_, err := svc.Update(context.Background(), cardID, tt.cardName, tt.limit, tt.closeDay, tt.dueDay, tt.isActive)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
