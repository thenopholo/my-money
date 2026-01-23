package domain

import "errors"

var (
	// Erros de User
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrEmptyName        = errors.New("name cannot be empty")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooWeak  = errors.New("password is too weak")
	ErrEmailInUse       = errors.New("email already in use")

	// Erros de BankAccount
	ErrEmptyBankName      = errors.New("bank name cannot be empty")
	ErrInvalidAccountType = errors.New("account type must be 'checking' or 'savings'")
	ErrNegativeBalance    = errors.New("initial balance cannot be negative")

	// Erros de Category
	ErrEmptyCategoryName   = errors.New("category name cannot be empty")
	ErrInvalidCategoryType = errors.New("category type must be 'income' or 'expense'")

	// Erros de CreditCard
	ErrEmptyCardName   = errors.New("card name cannot be empty")
	ErrInvalidCloseDay = errors.New("close day must be between 1 and 28")
	ErrInvalidLimit    = errors.New("credit limit must be positive")

	// Erros de Transaction
	ErrInvalidTransactionType = errors.New("transaction type must be 'income' or 'expense'")
	ErrTransactionInFuture    = errors.New("transaction date cannot be in the future")

	// Erros de Planned Income/Expense
	ErrInvalidDueDay      = errors.New("due day must be between 1 and 31")
	ErrInvalidAmount      = errors.New("amount must be positive")
	ErrInvalidFrequency   = errors.New("frequency must be 'once', 'monthly', or 'yearly'")
	ErrEndDateBeforeStart = errors.New("end date cannot be before start date")
	ErrEmptyDescription   = errors.New("description cannot be empty")

	// Erros de Invoice
	ErrInvoiceAlreadyPaid   = errors.New("invoice is already paid")
	ErrInvoiceAlreadyClosed = errors.New("invoice is already closed")
	ErrInvalidInvoice   = errors.New("invoice must have an amount bigger than zero")

	// Erros de CreditCardTransaction
	ErrInvalidInstallments = errors.New("installments must be between 1 and 48")
)