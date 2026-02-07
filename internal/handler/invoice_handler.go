package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thenopholo/my-money/internal/service"
)

type InvoiceHandler struct {
	service *service.InvoiceService
}

func NewInvoiceHandler(s *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service: s}
}

func (h *InvoiceHandler) CloseMonthInvoice(w http.ResponseWriter, r *http.Request) {
	cardID, err := parseUUID(chi.URLParam(r, "cardID"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid cardID")
		return
	}

	var req struct {
		ReferenceDate string `json:"reference_date"`
	}
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	referenceDate, err := parseDateOrNow(req.ReferenceDate)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid reference_date")
		return
	}

	invoice, err := h.service.CloseMonthInvoice(r.Context(), cardID, referenceDate)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, invoice)
}

func (h *InvoiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	invoice, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, invoice)
}

func (h *InvoiceHandler) ListByCard(w http.ResponseWriter, r *http.Request) {
	cardID, err := parseUUID(chi.URLParam(r, "cardID"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid cardID")
		return
	}

	invoices, err := h.service.GetByCardID(r.Context(), cardID)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, invoices)
}

func (h *InvoiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
