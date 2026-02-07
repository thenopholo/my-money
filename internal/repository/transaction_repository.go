package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/repository/postgres"
)

type transactionRepository struct {
	queries *postgres.Queries
}

func NewTransactionRepository(q *postgres.Queries) *transactionRepository {
	return &transactionRepository{queries: q}
}

func (r *transactionRepository) Create(ctx context.Context, t *domain.Transaction) error {
	dbTx, err := r.queries.CreateTransaction(ctx, postgres.CreateTransactionParams{
		AccountID:        t.AccountID,
		CategoryID:       t.CategoryID,
		InvoiceID:        uuidPtrToPgUUID(t.InvoiceID),
		PlannedIncomeID:  uuidPtrToPgUUID(t.PlannedIncome),
		PlannedExpenseID: uuidPtrToPgUUID(t.PlannedExpense),
		Amount:           decimalToPgNumeric(t.Amount),
		TransactionType:  string(t.TransactionType),
		Description:      t.Description,
		TransactionDate:  timeToPgDate(t.TransactionDate),
	})
	if err != nil {
		return err
	}

	t.ID = dbTx.ID
	t.CreatedAt = dbTx.CreatedAt.Time

	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	t, err := r.queries.GetTransactionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTransactionNotFound
		}

		return nil, err
	}

	return &domain.Transaction{
		ID:              t.ID,
		AccountID:       t.AccountID,
		CategoryID:      t.CategoryID,
		InvoiceID:       pgUUIDToUUIDPtr(t.InvoiceID),
		PlannedIncome:   pgUUIDToUUIDPtr(t.PlannedIncomeID),
		PlannedExpense:  pgUUIDToUUIDPtr(t.PlannedExpenseID),
		Amount:          pgNumericToDecimal(t.Amount),
		TransactionType: domain.TransactionType(t.TransactionType),
		Description:     t.Description,
		TransactionDate: pgDateToTime(t.TransactionDate),
		CreatedAt:       t.CreatedAt.Time,
	}, nil
}

func (r *transactionRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*domain.Transaction, error) {
	dbTxs, err := r.queries.GetTransactionByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return mapTransactions(dbTxs), nil
}

func (r *transactionRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*domain.Transaction, error) {
	dbTxs, err := r.queries.GetTransactionByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return mapTransactions(dbTxs), nil
}

func (r *transactionRepository) GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*domain.Transaction, error) {
	dbTxs, err := r.queries.GetTransactionByInvoiceID(ctx, uuidPtrToPgUUID(&invoiceID))
	if err != nil {
		return nil, err
	}

	return mapTransactions(dbTxs), nil
}

func (r *transactionRepository) GetByPlannedIncomeID(ctx context.Context, plannedIncomeID uuid.UUID) ([]*domain.Transaction, error) {
	dbTxs, err := r.queries.GetTransactionByPlannedIncomeID(ctx, uuidPtrToPgUUID(&plannedIncomeID))
	if err != nil {
		return nil, err
	}

	return mapTransactions(dbTxs), nil
}

func (r *transactionRepository) GetByPlannedExpenseID(ctx context.Context, plannedExpenseID uuid.UUID) ([]*domain.Transaction, error) {
	dbTxs, err := r.queries.GetTransactionByPlannedExpenseID(ctx, uuidPtrToPgUUID(&plannedExpenseID))
	if err != nil {
		return nil, err
	}

	return mapTransactions(dbTxs), nil
}

func (r *transactionRepository) Update(ctx context.Context, t *domain.Transaction) error {
	_, err := r.queries.UpdateTransaction(ctx, postgres.UpdateTransactionParams{
		ID:              t.ID,
		Amount:          decimalToPgNumeric(t.Amount),
		TransactionType: string(t.TransactionType),
		Description:     t.Description,
		TransactionDate: timeToPgDate(t.TransactionDate),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrTransactionNotFound
		}

		return err
	}

	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.queries.GetTransactionByID(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrTransactionNotFound
		}

		return err
	}

	return r.queries.DeleteTransaction(ctx, id)
}

func mapTransactions(dbTxs []postgres.Transaction) []*domain.Transaction {
	txs := make([]*domain.Transaction, len(dbTxs))
	for i, t := range dbTxs {
		txs[i] = &domain.Transaction{
			ID:              t.ID,
			AccountID:       t.AccountID,
			CategoryID:      t.CategoryID,
			InvoiceID:       pgUUIDToUUIDPtr(t.InvoiceID),
			PlannedIncome:   pgUUIDToUUIDPtr(t.PlannedIncomeID),
			PlannedExpense:  pgUUIDToUUIDPtr(t.PlannedExpenseID),
			Amount:          pgNumericToDecimal(t.Amount),
			TransactionType: domain.TransactionType(t.TransactionType),
			Description:     t.Description,
			TransactionDate: pgDateToTime(t.TransactionDate),
			CreatedAt:       t.CreatedAt.Time,
		}
	}

	return txs
}
