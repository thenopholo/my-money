package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/thenopholo/my-money/internal/domain"
)

type TransactionService struct {
	transactionRepo    TransactionRepository
	accountRepo        BankAccountRepository
	categoryRepo       CategoryRepository
	plannedIncomeRepo  PlannedIncomeRepository
	plannedExpenseRepo PlannedExpenseRepository
	invoiceRepo        InvoiceRepository
}

func NewTransactionService(
	transactionRepo TransactionRepository,
	accountRepo BankAccountRepository,
	categoryRepo CategoryRepository,
	plannedIncomeRepo PlannedIncomeRepository,
	plannedExpenseRepo PlannedExpenseRepository,
	invoiceRepo InvoiceRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo:    transactionRepo,
		accountRepo:        accountRepo,
		categoryRepo:       categoryRepo,
		plannedIncomeRepo:  plannedIncomeRepo,
		plannedExpenseRepo: plannedExpenseRepo,
		invoiceRepo:        invoiceRepo,
	}
}

// Create registra uma transacao manual (inclui PIX), atualiza saldo e persiste.
func (ts *TransactionService) Create(
	ctx context.Context,
	accountID, categoryID uuid.UUID,
	amount decimal.Decimal,
	txType domain.TransactionType,
	description string,
	transactionDate time.Time,
) (*domain.Transaction, error) {
	account, err := ts.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if err := ts.validateCategoryMatchesType(ctx, categoryID, txType); err != nil {
		return nil, err
	}

	tx, err := domain.NewTransaction(
		accountID,
		categoryID,
		nil,
		nil,
		nil,
		amount,
		txType,
		description,
		transactionDate,
	)
	if err != nil {
		return nil, err
	}

	if err := applyToAccount(account, txType, amount); err != nil {
		return nil, err
	}

	if err := ts.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	if err := ts.transactionRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// CreateFromPlannedIncome gera a transacao de uma receita prevista no dia de vencimento.
func (ts *TransactionService) CreateFromPlannedIncome(
	ctx context.Context,
	plannedIncomeID uuid.UUID,
	transactionDate time.Time,
) (*domain.Transaction, error) {
	planned, err := ts.plannedIncomeRepo.GetByID(ctx, plannedIncomeID)
	if err != nil {
		return nil, err
	}

	if !planned.IsActiveOn(transactionDate) {
		return nil, domain.ErrPlannedNotActive
	}
	if !planned.IsDueOn(transactionDate) {
		return nil, domain.ErrPlannedNotDueToday
	}

	account, err := ts.accountRepo.GetByID(ctx, planned.AccountID)
	if err != nil {
		return nil, err
	}

	if err := ts.validateCategoryMatchesType(ctx, planned.CategoryID, domain.TransactionTypeIncome); err != nil {
		return nil, err
	}

	tx, err := domain.NewTransaction(
		planned.AccountID,
		planned.CategoryID,
		nil,
		&planned.ID,
		nil,
		planned.Amount,
		domain.TransactionTypeIncome,
		planned.Description,
		transactionDate,
	)
	if err != nil {
		return nil, err
	}

	if err := account.ApplyIncome(planned.Amount); err != nil {
		return nil, err
	}

	if err := ts.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	if err := ts.transactionRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// CreateFromPlannedExpense gera a transacao de uma despesa prevista no dia de vencimento.
func (ts *TransactionService) CreateFromPlannedExpense(
	ctx context.Context,
	plannedExpenseID uuid.UUID,
	transactionDate time.Time,
) (*domain.Transaction, error) {
	planned, err := ts.plannedExpenseRepo.GetByID(ctx, plannedExpenseID)
	if err != nil {
		return nil, err
	}

	if !planned.IsActiveOn(transactionDate) {
		return nil, domain.ErrPlannedNotActive
	}
	if !planned.IsDueOn(transactionDate) {
		return nil, domain.ErrPlannedNotDueToday
	}

	account, err := ts.accountRepo.GetByID(ctx, planned.AccountID)
	if err != nil {
		return nil, err
	}

	if err := ts.validateCategoryMatchesType(ctx, planned.CategoryID, domain.TransactionTypeExpense); err != nil {
		return nil, err
	}

	tx, err := domain.NewTransaction(
		planned.AccountID,
		planned.CategoryID,
		nil,
		nil,
		&planned.ID,
		planned.Amount,
		domain.TransactionTypeExpense,
		planned.Description,
		transactionDate,
	)
	if err != nil {
		return nil, err
	}

	if err := account.ApplyExpense(planned.Amount); err != nil {
		return nil, err
	}

	if err := ts.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	if err := ts.transactionRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// PayInvoice paga fatura e cria uma transacao de despesa na conta.
func (ts *TransactionService) PayInvoice(
	ctx context.Context,
	invoiceID, accountID, categoryID uuid.UUID,
	paymentDate time.Time,
	description string,
) (*domain.Transaction, error) {
	invoice, err := ts.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	if err := invoice.MarkAsPaid(paymentDate); err != nil {
		return nil, err
	}

	account, err := ts.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if err := ts.validateCategoryMatchesType(ctx, categoryID, domain.TransactionTypeExpense); err != nil {
		return nil, err
	}

	tx, err := domain.NewTransaction(
		accountID,
		categoryID,
		&invoiceID,
		nil,
		nil,
		invoice.TotalAmount,
		domain.TransactionTypeExpense,
		description,
		paymentDate,
	)
	if err != nil {
		return nil, err
	}

	if err := account.ApplyExpense(invoice.TotalAmount); err != nil {
		return nil, err
	}

	if err := ts.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	if err := ts.invoiceRepo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	if err := ts.transactionRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (ts *TransactionService) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := ts.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	account, err := ts.accountRepo.GetByID(ctx, tx.AccountID)
	if err != nil {
		return err
	}

	// Reverte o impacto da transacao no saldo.
	switch tx.TransactionType {
	case domain.TransactionTypeIncome:
		if err := account.ApplyExpense(tx.Amount); err != nil {
			return err
		}
	case domain.TransactionTypeExpense:
		if err := account.ApplyIncome(tx.Amount); err != nil {
			return err
		}
	default:
		return domain.ErrInvalidTransactionType
	}

	if err := ts.accountRepo.Update(ctx, account); err != nil {
		return err
	}

	return ts.transactionRepo.Delete(ctx, id)
}

// Update altera dados da transacao e ajusta saldo da conta corretamente:
// reverte o impacto antigo e aplica o novo.
func (ts *TransactionService) Update(
	ctx context.Context,
	id uuid.UUID,
	categoryID uuid.UUID,
	amount decimal.Decimal,
	description string,
	transactionDate time.Time,
) (*domain.Transaction, error) {
	tx, err := ts.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := ts.validateCategoryMatchesType(ctx, categoryID, tx.TransactionType); err != nil {
		return nil, err
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, domain.ErrInvalidAmount
	}
	if transactionDate.After(time.Now().AddDate(0, 0, 1)) {
		return nil, domain.ErrTransactionInFuture
	}
	if description == "" {
		description = "Transacao"
	}

	account, err := ts.accountRepo.GetByID(ctx, tx.AccountID)
	if err != nil {
		return nil, err
	}

	// Reverte o impacto antigo no saldo.
	switch tx.TransactionType {
	case domain.TransactionTypeIncome:
		if err := account.ApplyExpense(tx.Amount); err != nil {
			return nil, err
		}
	case domain.TransactionTypeExpense:
		if err := account.ApplyIncome(tx.Amount); err != nil {
			return nil, err
		}
	default:
		return nil, domain.ErrInvalidTransactionType
	}

	// Aplica o novo impacto no saldo.
	if err := applyToAccount(account, tx.TransactionType, amount); err != nil {
		return nil, err
	}
	tx.CategoryID = categoryID
	tx.Amount = amount
	tx.Description = description
	tx.TransactionDate = transactionDate

	if err := ts.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}
	if err := ts.transactionRepo.Update(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (ts *TransactionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	return ts.transactionRepo.GetByID(ctx, id)
}

func (ts *TransactionService) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*domain.Transaction, error) {
	return ts.transactionRepo.GetByAccountID(ctx, accountID)
}

func (ts *TransactionService) validateCategoryMatchesType(
	ctx context.Context,
	categoryID uuid.UUID,
	txType domain.TransactionType,
) error {
	category, err := ts.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return err
	}

	if string(category.CategoryType) != string(txType) {
		return domain.ErrCategoryTypeMismatch
	}

	return nil
}

func applyToAccount(account *domain.BankAccount, txType domain.TransactionType, amount decimal.Decimal) error {
	switch txType {
	case domain.TransactionTypeIncome:
		return account.ApplyIncome(amount)
	case domain.TransactionTypeExpense:
		return account.ApplyExpense(amount)
	default:
		return domain.ErrInvalidTransactionType
	}
}
