package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

type CreditCardTransactionService struct {
	ccTxRepo     CreditCardTransactionRepository
	cardRepo     CreditCardRepository
	categoryRepo CategoryRepository
	invoiceRepo  InvoiceRepository
}

func NewCreditCardTransactionService(
	ccTxRepo CreditCardTransactionRepository,
	cardRepo CreditCardRepository,
	categoryRepo CategoryRepository,
	invoiceRepo InvoiceRepository,
) *CreditCardTransactionService {
	return &CreditCardTransactionService{
		ccTxRepo:     ccTxRepo,
		cardRepo:     cardRepo,
		categoryRepo: categoryRepo,
		invoiceRepo:  invoiceRepo,
	}
}

// Create registra uma compra no cartao; por padrao fica pendente de fatura (invoice nil).
func (s *CreditCardTransactionService) Create(
	ctx context.Context,
	cardID, categoryID uuid.UUID,
	amount decimal.Decimal,
	description string,
	installments int,
	transactionDate time.Time,
) (*domain.CreditCardTransaction, error) {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return nil, err
	}
	if !card.IsActive {
		return nil, domain.ErrCreditCardInactive
	}

	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	if category.CategoryType != domain.CategoryTypeExpense {
		return nil, domain.ErrCategoryTypeMismatch
	}

	if err := s.ensureCreditLimit(ctx, cardID, card.CreditLimit, amount); err != nil {
		return nil, err
	}

	currentInstallment := 1
	var installmentValue *decimal.Decimal
	if installments > 1 {
		value := amount.Div(decimal.NewFromInt(int64(installments)))
		installmentValue = &value
	}

	ccTx, err := domain.NewCreditCardTransaction(
		cardID,
		categoryID,
		uuid.Nil,
		amount,
		description,
		installments,
		currentInstallment,
		installmentValue,
	)
	if err != nil {
		return nil, err
	}
	ccTx.TransactionDate = transactionDate

	if err := s.ccTxRepo.Create(ctx, ccTx); err != nil {
		return nil, err
	}

	return ccTx, nil
}

func (s *CreditCardTransactionService) Update(
	ctx context.Context,
	id uuid.UUID,
	amount decimal.Decimal,
	description string,
	installments int,
	transactionDate time.Time,
) (*domain.CreditCardTransaction, error) {
	tx, err := s.ccTxRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, domain.ErrInvalidAmount
	}
	if description == "" {
		return nil, domain.ErrEmptyDescription
	}
	if installments < 1 || installments > 48 {
		return nil, domain.ErrInvalidInstallments
	}

	tx.Amount = amount
	tx.Description = description
	tx.Installments = installments
	tx.TransactionDate = transactionDate

	if installments > 1 {
		value := amount.Div(decimal.NewFromInt(int64(installments)))
		tx.InstallmentsValue = &value
	} else {
		tx.InstallmentsValue = nil
	}

	if err := s.ccTxRepo.Update(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *CreditCardTransactionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.ccTxRepo.Delete(ctx, id)
}

func (s *CreditCardTransactionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.CreditCardTransaction, error) {
	return s.ccTxRepo.GetByID(ctx, id)
}

func (s *CreditCardTransactionService) GetByCardID(ctx context.Context, cardID uuid.UUID) ([]*domain.CreditCardTransaction, error) {
	return s.ccTxRepo.GetByCardID(ctx, cardID)
}

func (s *CreditCardTransactionService) AssignToInvoice(
	ctx context.Context,
	transactionID, invoiceID uuid.UUID,
) error {
	tx, err := s.ccTxRepo.GetByID(ctx, transactionID)
	if err != nil {
		return err
	}

	invoice, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return err
	}
	if !invoice.Status.IsOpen() {
		return domain.ErrInvoiceNotOpen
	}
	if tx.CardID != invoice.CardID {
		return domain.ErrInvoiceCardMismatch
	}

	return s.ccTxRepo.AssignToInvoice(ctx, transactionID, invoiceID)
}

func (s *CreditCardTransactionService) ensureCreditLimit(
	ctx context.Context,
	cardID uuid.UUID,
	limit decimal.Decimal,
	newAmount decimal.Decimal,
) error {
	pending, err := s.ccTxRepo.GetPendingByCardID(ctx, cardID)
	if err != nil {
		return err
	}

	outstanding := decimal.Zero
	for _, tx := range pending {
		outstanding = outstanding.Add(tx.Amount)
	}

	invoices, err := s.invoiceRepo.GetByCardID(ctx, cardID)
	if err != nil {
		return err
	}

	for _, inv := range invoices {
		if inv.Status.IsOpen() {
			outstanding = outstanding.Add(inv.TotalAmount)
		}
	}

	if outstanding.Add(newAmount).GreaterThan(limit) {
		return domain.ErrCreditLimitExceeded
	}

	return nil
}
