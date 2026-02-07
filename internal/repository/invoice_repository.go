package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/repository/postgres"
)

type invoiceRepository struct {
	queries *postgres.Queries
}

func NewInvoiceRepository(q *postgres.Queries) *invoiceRepository {
	return &invoiceRepository{queries: q}
}

func (r *invoiceRepository) Create(ctx context.Context, in *domain.Invoice) error {
	dbInvoice, err := r.queries.CreateInvoice(ctx, postgres.CreateInvoiceParams{
		CardID:        in.CardID,
		ReferenceDate: timeToPgDate(in.ReferenceDate),
		DueDate:       timeToPgDate(in.DueDate),
		TotalAmount:   decimalToPgNumeric(in.TotalAmount),
		Status:        string(in.Status),
		PaidAt:        timePtrToPgTimestamptz(in.PaidAt),
	})
	if err != nil {
		return err
	}

	in.ID = dbInvoice.ID

	return nil
}

func (r *invoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	in, err := r.queries.GetInvoiceByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInvoiceNotFound
		}

		return nil, err
	}

	return &domain.Invoice{
		ID:            in.ID,
		CardID:        in.CardID,
		ReferenceDate: pgDateToTime(in.ReferenceDate),
		DueDate:       pgDateToTime(in.DueDate),
		TotalAmount:   pgNumericToDecimal(in.TotalAmount),
		Status:        domain.InvoiceStatus(in.Status),
		PaidAt:        pgTimestamptzToTimePtr(in.PaidAt),
	}, nil
}

func (r *invoiceRepository) GetByCardID(ctx context.Context, cardID uuid.UUID) ([]*domain.Invoice, error) {
	dbInvoices, err := r.queries.GetInvoiceByCardID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	invoices := make([]*domain.Invoice, len(dbInvoices))
	for i, in := range dbInvoices {
		invoices[i] = &domain.Invoice{
			ID:            in.ID,
			CardID:        in.CardID,
			ReferenceDate: pgDateToTime(in.ReferenceDate),
			DueDate:       pgDateToTime(in.DueDate),
			TotalAmount:   pgNumericToDecimal(in.TotalAmount),
			Status:        domain.InvoiceStatus(in.Status),
			PaidAt:        pgTimestamptzToTimePtr(in.PaidAt),
		}
	}

	return invoices, nil
}

func (r *invoiceRepository) Update(ctx context.Context, in *domain.Invoice) error {
	_, err := r.queries.UpdateInvoice(ctx, postgres.UpdateInvoiceParams{
		ID:          in.ID,
		TotalAmount: decimalToPgNumeric(in.TotalAmount),
		Status:      string(in.Status),
		PaidAt:      timePtrToPgTimestamptz(in.PaidAt),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrInvoiceNotFound
		}

		return err
	}

	return nil
}

func (r *invoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.queries.GetInvoiceByID(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrInvoiceNotFound
		}

		return err
	}

	return r.queries.DeleteInvoice(ctx, id)
}
