package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/thenopholo/my-money/internal/domain"
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

func (r *creditCardRepository) Create(ctx context.Context, cc *domain.CreditCard) error {
	dbCC, err := r.queries.CreateCreditCard(ctx, postgres.CreateCreditCardParams{
		UserID:      cc.UserID,
		Name:        cc.Name,
		CreditLimit: decimalToPgNumeric(cc.CreditLimit),
		CloseDay:    int32(cc.CloseDay),
		DueDay:      int32(cc.DueDay),
		IsActive:    true,
	})

	if err != nil {
		return err
	}

	cc.ID = dbCC.ID
	cc.CreatedAt = dbCC.CreatedAt.Time

	return nil
}

func (r *creditCardRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CreditCard, error) {
	cc, err := r.queries.GetCreditCardByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCreditCardNotFound
		}
		return nil, err
	}

	return &domain.CreditCard{
		ID:          cc.ID,
		UserID:      cc.UserID,
		Name:        cc.Name,
		CreditLimit: pgNumericToDecimal(cc.CreditLimit),
		CloseDay:    int(cc.CloseDay),
		DueDay:      int(cc.DueDay),
		IsActive:    true,
		CreatedAt:   cc.CreatedAt.Time,
	}, nil
}

func (r *creditCardRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.CreditCard, error) {
	dbCards, err := r.queries.GetCreditCardByUserID(ctx, userID)
	if err != nil {
			return nil, err
	}

	cards := make([]*domain.CreditCard, len(dbCards))
	for i, cc := range dbCards {
			cards[i] = &domain.CreditCard{
					ID:          cc.ID,
					UserID:      cc.UserID,
					Name:        cc.Name,
					CreditLimit: pgNumericToDecimal(cc.CreditLimit),
					CloseDay:    int(cc.CloseDay),
					DueDay:      int(cc.DueDay),
					IsActive:    cc.IsActive,
					CreatedAt:   cc.CreatedAt.Time,
			}
	}

	return cards, nil
}

func (r *creditCardRepository) Update(ctx context.Context, cc *domain.CreditCard) error {
	dbCC, err := r.queries.UpdateCreditCard(ctx, postgres.UpdateCreditCardParams{
			ID:          cc.ID,
			Name:        cc.Name,
			CreditLimit: decimalToPgNumeric(cc.CreditLimit),
			CloseDay:    int32(cc.CloseDay),
			DueDay:      int32(cc.DueDay),
			IsActive:    cc.IsActive,
	})
	if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
					return domain.ErrCreditCardNotFound
			}
			return err
	}

	cc.UpdatedAt = dbCC.UpdatedAt.Time

	return nil
}

func (r *creditCardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteCreditCard(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrCreditCardNotFound
		}
		return err
	}
	return nil
}
