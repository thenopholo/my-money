package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/repository/postgres"
)

type creditCardTransactionRepository struct {
	queries *postgres.Queries
}

func NewCreditCardTransactionRepository(q *postgres.Queries) *creditCardTransactionRepository {
	return &creditCardTransactionRepository{queries: q}
}

func (r *creditCardTransactionRepository) Create(ctx context.Context, t *domain.CreditCardTransaction) error {
	var invoiceID *uuid.UUID
	if t.InvoiceID != uuid.Nil {
		invoiceID = &t.InvoiceID
	}

	dbTx, err := r.queries.CreateCreditCardTransaction(ctx, postgres.CreateCreditCardTransactionParams{
		CardID:             t.CardID,
		CategoryID:         t.CategoryID,
		InvoiceID:          uuidPtrToPgUUID(invoiceID),
		Amount:             decimalToPgNumeric(t.Amount),
		Description:        t.Description,
		Installments:       int32(t.Installments),
		CurrentInstallment: int32(t.CurrentInstallments),
		InstallmentValue:   decimalToPgNumericPtr(t.InstallmentsValue),
		TransactionDate:    timeToPgDate(t.TransactionDate),
	})
	if err != nil {
		return err
	}

	t.ID = dbTx.ID
	t.CreatedAt = dbTx.CreatedAt.Time

	if invoice := pgUUIDToUUIDPtr(dbTx.InvoiceID); invoice != nil {
		t.InvoiceID = *invoice
	}

	return nil
}

func (r *creditCardTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CreditCardTransaction, error) {
	t, err := r.queries.GetCreditCardTransactionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCreditCardTransactionNotFound
		}

		return nil, err
	}

	return mapCreditCardTransaction(t), nil
}

func (r *creditCardTransactionRepository) GetByCardID(ctx context.Context, cardID uuid.UUID) ([]*domain.CreditCardTransaction, error) {
	dbTxs, err := r.queries.GetCreditCardTransactionsByCardID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	return mapCreditCardTransactions(dbTxs), nil
}

func (r *creditCardTransactionRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*domain.CreditCardTransaction, error) {
	dbTxs, err := r.queries.GetCreditCardTransactionsByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return mapCreditCardTransactions(dbTxs), nil
}

func (r *creditCardTransactionRepository) GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*domain.CreditCardTransaction, error) {
	dbTxs, err := r.queries.GetCreditCardTransactionsByInvoiceID(ctx, uuidPtrToPgUUID(&invoiceID))
	if err != nil {
		return nil, err
	}

	return mapCreditCardTransactions(dbTxs), nil
}

func (r *creditCardTransactionRepository) GetPendingByCardID(ctx context.Context, cardID uuid.UUID) ([]*domain.CreditCardTransaction, error) {
	dbTxs, err := r.queries.GetPendingCreditCardTransactions(ctx, cardID)
	if err != nil {
		return nil, err
	}

	return mapCreditCardTransactions(dbTxs), nil
}

func (r *creditCardTransactionRepository) AssignToInvoice(ctx context.Context, transactionID, invoiceID uuid.UUID) error {
	if _, err := r.queries.GetCreditCardTransactionByID(ctx, transactionID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrCreditCardTransactionNotFound
		}

		return err
	}

	return r.queries.AssignTransactionToInvoice(ctx, postgres.AssignTransactionToInvoiceParams{
		ID:        transactionID,
		InvoiceID: uuidPtrToPgUUID(&invoiceID),
	})
}

func (r *creditCardTransactionRepository) Update(ctx context.Context, t *domain.CreditCardTransaction) error {
	_, err := r.queries.UpdateCreditCardTransaction(ctx, postgres.UpdateCreditCardTransactionParams{
		ID:                 t.ID,
		CategoryID:         t.CategoryID,
		Amount:             decimalToPgNumeric(t.Amount),
		Description:        t.Description,
		Installments:       int32(t.Installments),
		CurrentInstallment: int32(t.CurrentInstallments),
		InstallmentValue:   decimalToPgNumericPtr(t.InstallmentsValue),
		TransactionDate:    timeToPgDate(t.TransactionDate),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrCreditCardTransactionNotFound
		}

		return err
	}

	return nil
}

func (r *creditCardTransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.queries.GetCreditCardTransactionByID(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrCreditCardTransactionNotFound
		}

		return err
	}

	return r.queries.DeleteCreditCardTransaction(ctx, id)
}

func (r *creditCardTransactionRepository) DeleteAllByCardID(ctx context.Context, cardID uuid.UUID) error {
	return r.queries.DeleteCreditCardTransactionsByCardID(ctx, cardID)
}

func mapCreditCardTransaction(t postgres.CreditCardTransaction) *domain.CreditCardTransaction {
	var invoiceID uuid.UUID
	if invoice := pgUUIDToUUIDPtr(t.InvoiceID); invoice != nil {
		invoiceID = *invoice
	}

	return &domain.CreditCardTransaction{
		ID:                  t.ID,
		CardID:              t.CardID,
		CategoryID:          t.CategoryID,
		InvoiceID:           invoiceID,
		Amount:              pgNumericToDecimal(t.Amount),
		Description:         t.Description,
		Installments:        int(t.Installments),
		CurrentInstallments: int(t.CurrentInstallment),
		InstallmentsValue:   pgNumericToDecimalPtr(t.InstallmentValue),
		TransactionDate:     pgDateToTime(t.TransactionDate),
		CreatedAt:           t.CreatedAt.Time,
	}
}

func mapCreditCardTransactions(dbTxs []postgres.CreditCardTransaction) []*domain.CreditCardTransaction {
	txs := make([]*domain.CreditCardTransaction, len(dbTxs))
	for i, tx := range dbTxs {
		txs[i] = mapCreditCardTransaction(tx)
	}

	return txs
}
