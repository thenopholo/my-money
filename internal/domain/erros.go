package domain

import "errors"

var (
	// Erros de User
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrEmptyName          = errors.New("name cannot be empty")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrPasswordTooWeak    = errors.New("password is too weak")
	ErrEmailInUse         = errors.New("email already in use")
	ErrUserNotFound       = errors.New("user not found in DB")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSamePassword       = errors.New("new password must be different from current password")

	// Erros de BankAccount
	ErrEmptyBankName       = errors.New("bank name cannot be empty")
	ErrInvalidAccountType  = errors.New("account type must be 'checking' or 'savings'")
	ErrNegativeBalance     = errors.New("initial balance cannot be negative")
	ErrBankAccountNotFound = errors.New("bank account not found")

	// Erros de Category
	ErrEmptyCategoryName   = errors.New("category name cannot be empty")
	ErrInvalidCategoryType = errors.New("category type must be 'income' or 'expense'")
	ErrCategoryNotFound    = errors.New("category not found")

	// Erros de CreditCard
	ErrEmptyCardName      = errors.New("card name cannot be empty")
	ErrInvalidCloseDay    = errors.New("close day must be between 1 and 28")
	ErrInvalidLimit       = errors.New("credit limit must be positive")
	ErrCreditCardNotFound = errors.New("credit card not found")

	// Erros de Invoice
	ErrInvoiceNotFound = errors.New("invoice not found")

	// Erros de Transaction
	ErrInvalidTransactionType = errors.New("transaction type must be 'income' or 'expense'")
	ErrTransactionInFuture    = errors.New("transaction date cannot be in the future")
	ErrTransactionNotFound    = errors.New("transaction not found")

	// Erros de Planned Income/Expense
	ErrInvalidDueDay          = errors.New("due day must be between 1 and 31")
	ErrInvalidAmount          = errors.New("amount must be positive")
	ErrInvalidFrequency       = errors.New("frequency must be 'once', 'monthly', or 'yearly'")
	ErrEndDateBeforeStart     = errors.New("end date cannot be before start date")
	ErrEmptyDescription       = errors.New("description cannot be empty")
	ErrPlannedIncomeNotFound  = errors.New("planned income not found")
	ErrPlannedExpenseNotFound = errors.New("planned expense not found")

	// Erros de Invoice (regra de negocio)
	ErrInvoiceAlreadyPaid   = errors.New("invoice is already paid")
	ErrInvoiceAlreadyClosed = errors.New("invoice is already closed")
	ErrInvalidInvoice       = errors.New("invoice must have an amount bigger than zero")

	// Erros de CreditCardTransaction
	ErrInvalidInstallments           = errors.New("installments must be between 1 and 48")
	ErrCreditCardTransactionNotFound = errors.New("credit card transaction not found")
)
