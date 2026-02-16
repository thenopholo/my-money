package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thenopholo/my-money/internal/service"
)

type CreditCardTransactionHandler struct {
	service     *service.CreditCardTransactionService
	cardService *service.CreditCardService
}

func NewCreditCardTransactionHandler(s *service.CreditCardTransactionService, cardService *service.CreditCardService) *CreditCardTransactionHandler {
	return &CreditCardTransactionHandler{service: s, cardService: cardService}
}

func (h *CreditCardTransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	cardID, err := parseUUID(chi.URLParam(r, "cardID"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid cardID")
		return
	}

	// Verifica ownership do cartão
	card, err := h.cardService.GetByID(r.Context(), cardID)
	if err != nil {
		Error(w, http.StatusNotFound, "credit card not found")
		return
	}
	if card.UserID != userID {
		Error(w, http.StatusForbidden, "forbidden")
		return
	}

	var req struct {
		CategoryID      string `json:"category_id"`
		Amount          string `json:"amount"`
		Description     string `json:"description"`
		Installments    int    `json:"installments"`
		TransactionDate string `json:"transaction_date"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
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
	transactionDate, err := parseDateOrNow(req.TransactionDate)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid transaction_date")
		return
	}

	tx, err := h.service.Create(r.Context(), cardID, categoryID, amount, req.Description, req.Installments, transactionDate)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, tx)
}

func (h *CreditCardTransactionHandler) ListByCard(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	cardID, err := parseUUID(chi.URLParam(r, "cardID"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid cardID")
		return
	}

	// Verifica ownership do cartão
	card, err := h.cardService.GetByID(r.Context(), cardID)
	if err != nil {
		Error(w, http.StatusNotFound, "credit card not found")
		return
	}
	if card.UserID != userID {
		Error(w, http.StatusForbidden, "forbidden")
		return
	}

	txs, err := h.service.GetByCardID(r.Context(), cardID)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, txs)
}

func (h *CreditCardTransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

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

	// Verifica ownership via cartão
	card, err := h.cardService.GetByID(r.Context(), tx.CardID)
	if err != nil {
		Error(w, http.StatusNotFound, "credit card not found")
		return
	}
	if card.UserID != userID {
		Error(w, http.StatusForbidden, "forbidden")
		return
	}

	JSON(w, http.StatusOK, tx)
}

func (h *CreditCardTransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Verifica ownership via transação → cartão
	existing, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	card, err := h.cardService.GetByID(r.Context(), existing.CardID)
	if err != nil {
		Error(w, http.StatusNotFound, "credit card not found")
		return
	}
	if card.UserID != userID {
		Error(w, http.StatusForbidden, "forbidden")
		return
	}

	var req struct {
		CategoryID      string `json:"category_id"`
		Amount          string `json:"amount"`
		Description     string `json:"description"`
		Installments    int    `json:"installments"`
		TransactionDate string `json:"transaction_date"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
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
	transactionDate, err := parseDateOrNow(req.TransactionDate)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid transaction_date")
		return
	}

	tx, err := h.service.Update(r.Context(), id, categoryID, amount, req.Description, req.Installments, transactionDate)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, tx)
}

func (h *CreditCardTransactionHandler) AssignToInvoice(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Verifica ownership via transação → cartão
	existing, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	card, err := h.cardService.GetByID(r.Context(), existing.CardID)
	if err != nil {
		Error(w, http.StatusNotFound, "credit card not found")
		return
	}
	if card.UserID != userID {
		Error(w, http.StatusForbidden, "forbidden")
		return
	}

	var req struct {
		InvoiceID string `json:"invoice_id"`
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

	if err := h.service.AssignToInvoice(r.Context(), id, invoiceID); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{"message": "transaction assigned to invoice"})
}

func (h *CreditCardTransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Verifica ownership via transação → cartão
	existing, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	card, err := h.cardService.GetByID(r.Context(), existing.CardID)
	if err != nil {
		Error(w, http.StatusNotFound, "credit card not found")
		return
	}
	if card.UserID != userID {
		Error(w, http.StatusForbidden, "forbidden")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
