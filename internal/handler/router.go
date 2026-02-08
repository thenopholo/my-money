package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	userHandler *UserHandler,
	bankAccountHandler *BankAccountHandler,
	categoryHandler *CategoryHandler,
	creditCardHandler *CreditCardHandler,
	plannedIncomeHandler *PlannedIncomeHandler,
	plannedExpenseHandler *PlannedExpenseHandler,
	transactionHandler *TransactionHandler,
	creditCardTransactionHandler *CreditCardTransactionHandler,
	invoiceHandler *InvoiceHandler,
	importHandler *ImportHandler,
	authMiddleware func(http.Handler) http.Handler,
	corsAllowedOrigins []string,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	if userHandler != nil {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
		})
	}

	r.Route("/api", func(r chi.Router) {
		if authMiddleware != nil {
			r.Use(authMiddleware)
		}

		if userHandler != nil {
			r.Put("/me/password", userHandler.UpdatePassword)
		}

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

		if plannedIncomeHandler != nil {
			r.Route("/planned-incomes", func(r chi.Router) {
				r.Post("/", plannedIncomeHandler.Create)
				r.Get("/", plannedIncomeHandler.List)
				r.Get("/{id}", plannedIncomeHandler.GetByID)
				r.Put("/{id}", plannedIncomeHandler.Update)
				r.Delete("/{id}", plannedIncomeHandler.Delete)
			})
		}

		if plannedExpenseHandler != nil {
			r.Route("/planned-expenses", func(r chi.Router) {
				r.Post("/", plannedExpenseHandler.Create)
				r.Get("/", plannedExpenseHandler.List)
				r.Get("/{id}", plannedExpenseHandler.GetByID)
				r.Put("/{id}", plannedExpenseHandler.Update)
				r.Delete("/{id}", plannedExpenseHandler.Delete)
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

		if importHandler != nil {
			r.Route("/import", func(r chi.Router) {
				r.Post("/preview", importHandler.Preview)
				r.Post("/confirm", importHandler.Confirm)
			})
		}
	})

	return r
}
