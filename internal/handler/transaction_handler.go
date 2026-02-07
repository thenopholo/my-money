package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/service"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID       string `json:"account_id"`
		CategoryID      string `json:"category_id"`
		Amount          string `json:"amount"`
		TransactionType string `json:"transaction_type"`
		Description     string `json:"description"`
		TransactionDate string `json:"transaction_date"`
	}

	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	accountID, err := parseUUID(req.AccountID)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid account_id")
		return
	}
	categoryID, err := parseUUID(req.CategoryID)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid category_id")
		return
	}
	amount, err := parseDecimal(req.Amount)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid amount")
		return
	}
	transactionDate, err := parseDate(req.TransactionDate)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid transaction_date")
		return
	}

	tx, err := h.service.Create(
		r.Context(),
		accountID,
		categoryID,
		amount,
		domain.TransactionType(req.TransactionType),
		req.Description,
		transactionDate,
	)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, tx)
}

func (h *TransactionHandler) CreateFromPlannedIncome(w http.ResponseWriter, r *http.Request) {
	plannedIncomeID, err := parseUUID(chi.URLParam(r, "plannedIncomeID"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid plannedIncomeID")
		return
	}

	var req struct {
		TransactionDate string `json:"transaction_date"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	transactionDate, err := parseDate(req.TransactionDate)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid transaction_date")
		return
	}

	tx, err := h.service.CreateFromPlannedIncome(r.Context(), plannedIncomeID, transactionDate)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, tx)
}

func (h *TransactionHandler) CreateFromPlannedExpense(w http.ResponseWriter, r *http.Request) {
	plannedExpenseID, err := parseUUID(chi.URLParam(r, "plannedExpenseID"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid plannedExpenseID")
		return
	}

	var req struct {
		TransactionDate string `json:"transaction_date"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	transactionDate, err := parseDate(req.TransactionDate)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid transaction_date")
		return
	}

	tx, err := h.service.CreateFromPlannedExpense(r.Context(), plannedExpenseID, transactionDate)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, tx)
}

func (h *TransactionHandler) PayInvoice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		InvoiceID   string `json:"invoice_id"`
		AccountID   string `json:"account_id"`
		CategoryID  string `json:"category_id"`
		PaymentDate string `json:"payment_date"`
		Description string `json:"description"`
	}

	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	invoiceID, err := parseUUID(req.InvoiceID)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid invoice_id")
		return
	}
	accountID, err := parseUUID(req.AccountID)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid account_id")
		return
	}
	categoryID, err := parseUUID(req.CategoryID)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid category_id")
		return
	}
	paymentDate, err := parseDate(req.PaymentDate)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid payment_date")
		return
	}

	tx, err := h.service.PayInvoice(r.Context(), invoiceID, accountID, categoryID, paymentDate, req.Description)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, tx)
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	tx, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, tx)
}

func (h *TransactionHandler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	accountID, err := parseUUID(chi.URLParam(r, "accountID"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid accountID")
		return
	}

	txs, err := h.service.GetByAccountID(r.Context(), accountID)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, txs)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Amount          string `json:"amount"`
		Description     string `json:"description"`
		TransactionDate string `json:"transaction_date"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	amount, err := parseDecimal(req.Amount)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid amount")
		return
	}
	transactionDate, err := parseDate(req.TransactionDate)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid transaction_date")
		return
	}

	tx, err := h.service.Update(r.Context(), id, amount, req.Description, transactionDate)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, tx)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
