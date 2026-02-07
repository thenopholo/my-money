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
	Create(ctx context.Context, cc *domain.CreditCard) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.CreditCard, error)
	GetByUserID(ctx context.Context, id uuid.UUID) ([]*domain.CreditCard, error)
	Update(ctx context.Context, cc *domain.CreditCard) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetByUserID(ctx context.Context, id uuid.UUID) ([]*domain.Category, error)
	Update(ctx context.Context, c *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BankAccountRepository interface {
	Create(ctx context.Context, ba *domain.BankAccount) (*domain.BankAccount, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BankAccount, error)
	GetByUserID(ctx context.Context, id uuid.UUID) ([]*domain.BankAccount, error)
	Update(ctx context.Context, ba *domain.BankAccount) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type InvoiceRepository interface {
	Create(ctx context.Context, in *domain.Invoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error)
	GetByCardID(ctx context.Context, cardID uuid.UUID) ([]*domain.Invoice, error)
	Update(ctx context.Context, in *domain.Invoice) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PlannedIncomeRepository interface {
	Create(ctx context.Context, pi *domain.PlannedIncome) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PlannedIncome, error)
	GetByUserID(ctx context.Context, id uuid.UUID) ([]*domain.PlannedIncome, error)
	Update(ctx context.Context, pi *domain.PlannedIncome) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PlannedExpenseRepository interface {
	Create(ctx context.Context, pe *domain.PlannedExpense) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PlannedExpense, error)
	GetByUserID(ctx context.Context, id uuid.UUID) ([]*domain.PlannedExpense, error)
	Update(ctx context.Context, pe *domain.PlannedExpense) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TransactionRepository interface {
	Create(ctx context.Context, t *domain.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*domain.Transaction, error)
	GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*domain.Transaction, error)
	GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*domain.Transaction, error)
	GetByPlannedIncomeID(ctx context.Context, plannedIncomeID uuid.UUID) ([]*domain.Transaction, error)
	GetByPlannedExpenseID(ctx context.Context, plannedExpenseID uuid.UUID) ([]*domain.Transaction, error)
	Update(ctx context.Context, t *domain.Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreditCardTransactionRepository interface {
	Create(ctx context.Context, t *domain.CreditCardTransaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.CreditCardTransaction, error)
	GetByCardID(ctx context.Context, cardID uuid.UUID) ([]*domain.CreditCardTransaction, error)
	GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]*domain.CreditCardTransaction, error)
	GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*domain.CreditCardTransaction, error)
	GetPendingByCardID(ctx context.Context, cardID uuid.UUID) ([]*domain.CreditCardTransaction, error)
	AssignToInvoice(ctx context.Context, transactionID, invoiceID uuid.UUID) error
	Update(ctx context.Context, t *domain.CreditCardTransaction) error
	Delete(ctx context.Context, id uuid.UUID) error
}