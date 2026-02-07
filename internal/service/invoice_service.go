package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

type InvoiceService struct {
	invoiceRepo InvoiceRepository
	cardRepo    CreditCardRepository
	ccTxRepo    CreditCardTransactionRepository
}

func NewInvoiceService(
	invoiceRepo InvoiceRepository,
	cardRepo CreditCardRepository,
	ccTxRepo CreditCardTransactionRepository,
) *InvoiceService {
	return &InvoiceService{
		invoiceRepo: invoiceRepo,
		cardRepo:    cardRepo,
		ccTxRepo:    ccTxRepo,
	}
}

// CloseMonthInvoice fecha a fatura de um mes, somando compras pendentes do periodo e vinculando-as.
func (s *InvoiceService) CloseMonthInvoice(
	ctx context.Context,
	cardID uuid.UUID,
	referenceDate time.Time,
) (*domain.Invoice, error) {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return nil, err
	}
	if !card.IsActive {
		return nil, domain.ErrCreditCardInactive
	}

	pending, err := s.ccTxRepo.GetPendingByCardID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	monthTxs := filterTransactionsByMonth(pending, referenceDate)
	if len(monthTxs) == 0 {
		return nil, domain.ErrNoPendingTransactions
	}

	total := decimal.Zero
	for _, tx := range monthTxs {
		total = total.Add(tx.Amount)
	}

	refMonth := monthStart(referenceDate)
	dueDate := nextMonthWithDay(refMonth, card.DueDay)

	invoice, err := domain.NewInvoice(
		cardID,
		refMonth,
		dueDate,
		nil,
		total,
		domain.InvoiceStatusOpen,
	)
	if err != nil {
		return nil, err
	}

	if err := s.invoiceRepo.Create(ctx, invoice); err != nil {
		return nil, err
	}

	for _, tx := range monthTxs {
		if err := s.ccTxRepo.AssignToInvoice(ctx, tx.ID, invoice.ID); err != nil {
			return nil, err
		}
	}

	return invoice, nil
}

func (s *InvoiceService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	return s.invoiceRepo.GetByID(ctx, id)
}

func (s *InvoiceService) GetByCardID(ctx context.Context, cardID uuid.UUID) ([]*domain.Invoice, error) {
	return s.invoiceRepo.GetByCardID(ctx, cardID)
}

func (s *InvoiceService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.invoiceRepo.Delete(ctx, id)
}

func filterTransactionsByMonth(txs []*domain.CreditCardTransaction, date time.Time) []*domain.CreditCardTransaction {
	year, month := date.Year(), date.Month()
	filtered := make([]*domain.CreditCardTransaction, 0, len(txs))
	for _, tx := range txs {
		if tx.TransactionDate.Year() == year && tx.TransactionDate.Month() == month {
			filtered = append(filtered, tx)
		}
	}

	return filtered
}

func monthStart(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

func nextMonthWithDay(date time.Time, day int) time.Time {
	next := date.AddDate(0, 1, 0)
	lastDay := daysInMonth(next.Year(), next.Month())
	if day > lastDay {
		day = lastDay
	}

	return time.Date(next.Year(), next.Month(), day, 0, 0, 0, 0, date.Location())
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
