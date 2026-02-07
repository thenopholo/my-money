package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

type CreditCardService struct {
	creditCardRepo CreditCardRepository
}

func NewCreditCardService(creditCardRepo CreditCardRepository) *CreditCardService {
	return &CreditCardService{creditCardRepo: creditCardRepo}
}

func (cs *CreditCardService) Create(ctx context.Context, userID uuid.UUID, name string, creditLimit decimal.Decimal, closeDay, dueDay int) (*domain.CreditCard, error) {
	creditCard, err := domain.NewCreditCard(userID, name, creditLimit, closeDay, dueDay)
	if err != nil {
		return nil, err
	}

	if err := cs.creditCardRepo.Create(ctx, creditCard); err != nil {
		return nil, err
	}

	return creditCard, nil
}

func (cs *CreditCardService) GetByID(ctx context.Context, id uuid.UUID) (*domain.CreditCard, error) {
	return cs.creditCardRepo.GetByID(ctx, id)
}

func (cs *CreditCardService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.CreditCard, error) {
	return cs.creditCardRepo.GetByUserID(ctx, userID)
}

func (cs *CreditCardService) Update(ctx context.Context, id uuid.UUID, name string, creditLimit decimal.Decimal, closeDay, dueDay int, isActive bool) (*domain.CreditCard, error) {
	creditCard, err := cs.creditCardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := creditCard.SetName(name); err != nil {
		return nil, err
	}
	if err := creditCard.SetCreditLimit(creditLimit); err != nil {
		return nil, err
	}
	if err := creditCard.SetCloseDay(closeDay); err != nil {
		return nil, err
	}
	if err := creditCard.SetDueDay(dueDay); err != nil {
		return nil, err
	}
	creditCard.SetIsActive(isActive)

	if err := cs.creditCardRepo.Update(ctx, creditCard); err != nil {
		return nil, err
	}

	return creditCard, nil
}

func (cs *CreditCardService) Delete(ctx context.Context, id uuid.UUID) error {
	return cs.creditCardRepo.Delete(ctx, id)
}
