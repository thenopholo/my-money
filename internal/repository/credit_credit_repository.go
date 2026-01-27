package repository

import (
	// "context"

	// "github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/repository/postgres"
)

type creditCardRepository struct {
	queries *postgres.Queries
}

func NewCreditCardRepository(q *postgres.Queries) *creditCardRepository {
  return &creditCardRepository{
  	queries: q,
  }
}

// func (r *creditCardRepository) CreateCreditCard(ctx context.Context, cc *domain.CreditCard) error {

// }