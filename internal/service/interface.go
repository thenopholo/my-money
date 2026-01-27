package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreditCardRepository interface {
	CreateCreditCard(ctx context.Context, cc *domain.CreditCard) error
	GetCreditCardByID(ctx context.Context, id uuid.UUID) (*domain.CreditCard, error)
	GetCreditCardUserByID(ctx context.Context, id uuid.UUID) (*domain.CreditCard, error)
	UpadateCreditCard(ctx context.Context, cc *domain.CreditCard) error
	DeleteCredirCard(ctx context.Context, id uuid.UUID) error
}
