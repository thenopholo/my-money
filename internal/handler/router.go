package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	userHandler *UserHandler,
	bankAccountHandler *BankAccountHandler,
	categoryHandler *CategoryHandler,
	creditCardHandler *CreditCardHandler,
	transactionHandler *TransactionHandler,
	creditCardTransactionHandler *CreditCardTransactionHandler,
	invoiceHandler *InvoiceHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	if userHandler != nil {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
			r.Put("/password", userHandler.UpdatePassword)
		})
	}

	r.Route("/api", func(r chi.Router) {
		if bankAccountHandler != nil {
			r.Route("/accounts", func(r chi.Router) {
				r.Post("/", bankAccountHandler.Create)
				r.Get("/", bankAccountHandler.List)
				r.Get("/{id}", bankAccountHandler.GetByID)
				r.Put("/{id}", bankAccountHandler.Update)
				r.Delete("/{id}", bankAccountHandler.Delete)
			})
		}

		if categoryHandler != nil {
			r.Route("/categories", func(r chi.Router) {
				r.Post("/", categoryHandler.Create)
				r.Get("/", categoryHandler.List)
				r.Get("/{id}", categoryHandler.GetByID)
				r.Put("/{id}", categoryHandler.Update)
				r.Delete("/{id}", categoryHandler.Delete)
			})
		}

		if creditCardHandler != nil {
			r.Route("/credit-cards", func(r chi.Router) {
				r.Post("/", creditCardHandler.Create)
				r.Get("/", creditCardHandler.List)
				r.Get("/{id}", creditCardHandler.GetByID)
				r.Put("/{id}", creditCardHandler.Update)
				r.Delete("/{id}", creditCardHandler.Delete)
			})
		}

		if transactionHandler != nil {
			r.Route("/transactions", func(r chi.Router) {
				r.Post("/", transactionHandler.Create)
				r.Get("/{id}", transactionHandler.GetByID)
				r.Put("/{id}", transactionHandler.Update)
				r.Delete("/{id}", transactionHandler.Delete)
				r.Post("/planned-income/{plannedIncomeID}", transactionHandler.CreateFromPlannedIncome)
				r.Post("/planned-expense/{plannedExpenseID}", transactionHandler.CreateFromPlannedExpense)
				r.Post("/pay-invoice", transactionHandler.PayInvoice)
			})
			r.Get("/accounts/{accountID}/transactions", transactionHandler.ListByAccount)
		}

		if creditCardTransactionHandler != nil {
			r.Post("/credit-cards/{cardID}/transactions", creditCardTransactionHandler.Create)
			r.Get("/credit-cards/{cardID}/transactions", creditCardTransactionHandler.ListByCard)
			r.Get("/credit-card-transactions/{id}", creditCardTransactionHandler.GetByID)
			r.Put("/credit-card-transactions/{id}", creditCardTransactionHandler.Update)
			r.Delete("/credit-card-transactions/{id}", creditCardTransactionHandler.Delete)
			r.Post("/credit-card-transactions/{id}/assign-invoice", creditCardTransactionHandler.AssignToInvoice)
		}

		if invoiceHandler != nil {
			r.Post("/credit-cards/{cardID}/invoices/close-month", invoiceHandler.CloseMonthInvoice)
			r.Get("/credit-cards/{cardID}/invoices", invoiceHandler.ListByCard)
			r.Get("/invoices/{id}", invoiceHandler.GetByID)
			r.Delete("/invoices/{id}", invoiceHandler.Delete)
		}
	})

	return r
}
