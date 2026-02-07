package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/repository/postgres"
)

type bankAccountRepository struct {
	queries *postgres.Queries
}

func NewBankAccountRepository(q *postgres.Queries) *bankAccountRepository {
	return &bankAccountRepository{queries: q}
}

func (r *bankAccountRepository) Create(ctx context.Context, ba *domain.BankAccount) (*domain.BankAccount, error) {
	dbAccount, err := r.queries.CreateBankAccount(ctx, postgres.CreateBankAccountParams{
		UserID:      ba.UserID,
		Name:        ba.Name,
		BankName:    ba.BankName,
		AccountType: string(ba.AccountType),
		Balance:     decimalToPgNumeric(ba.Balance),
		IsActive:    ba.IsActive,
	})
	if err != nil {
		return nil, err
	}

	ba.ID = dbAccount.ID
	ba.CreatedAt = dbAccount.CreatedAt.Time
	ba.UpdatedAt = dbAccount.UpdatedAt.Time

	return ba, nil
}

func (r *bankAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.BankAccount, error) {
	ba, err := r.queries.GetBankAccountByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBankAccountNotFound
		}

		return nil, err
	}

	return &domain.BankAccount{
		ID:          ba.ID,
		UserID:      ba.UserID,
		Name:        ba.Name,
		BankName:    ba.BankName,
		AccountType: domain.AccountType(ba.AccountType),
		Balance:     pgNumericToDecimal(ba.Balance),
		IsActive:    ba.IsActive,
		CreatedAt:   ba.CreatedAt.Time,
		UpdatedAt:   ba.UpdatedAt.Time,
	}, nil
}

func (r *bankAccountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.BankAccount, error) {
	dbAccounts, err := r.queries.GetBankAccountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	accounts := make([]*domain.BankAccount, len(dbAccounts))
	for i, ba := range dbAccounts {
		accounts[i] = &domain.BankAccount{
			ID:          ba.ID,
			UserID:      ba.UserID,
			Name:        ba.Name,
			BankName:    ba.BankName,
			AccountType: domain.AccountType(ba.AccountType),
			Balance:     pgNumericToDecimal(ba.Balance),
			IsActive:    ba.IsActive,
			CreatedAt:   ba.CreatedAt.Time,
			UpdatedAt:   ba.UpdatedAt.Time,
		}
	}

	return accounts, nil
}

func (r *bankAccountRepository) Update(ctx context.Context, ba *domain.BankAccount) error {
	dbAccount, err := r.queries.UpdateBankAccount(ctx, postgres.UpdateBankAccountParams{
		ID:          ba.ID,
		Name:        ba.Name,
		BankName:    ba.BankName,
		AccountType: string(ba.AccountType),
		Balance:     decimalToPgNumeric(ba.Balance),
		IsActive:    ba.IsActive,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrBankAccountNotFound
		}

		return err
	}

	ba.UpdatedAt = dbAccount.UpdatedAt.Time

	return nil
}

func (r *bankAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.queries.GetBankAccountByID(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrBankAccountNotFound
		}

		return err
	}

	return r.queries.DeleteBankAccount(ctx, id)
}
