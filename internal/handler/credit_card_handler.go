package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thenopholo/my-money/internal/service"
)

type CreditCardHandler struct {
	service *service.CreditCardService
}

func NewCreditCardHandler(s *service.CreditCardService) *CreditCardHandler {
	return &CreditCardHandler{service: s}
}

func (h *CreditCardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID      string `json:"user_id"`
		Name        string `json:"name"`
		CreditLimit string `json:"credit_limit"`
		CloseDay    int    `json:"close_day"`
		DueDay      int    `json:"due_day"`
	}

	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := parseUUID(req.UserID)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid user_id")
		return
	}
	creditLimit, err := parseDecimal(req.CreditLimit)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid credit_limit")
		return
	}

	card, err := h.service.Create(r.Context(), userID, req.Name, creditLimit, req.CloseDay, req.DueDay)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, card)
}

func (h *CreditCardHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUUID(r.URL.Query().Get("user_id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	cards, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, cards)
}

func (h *CreditCardHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	card, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, card)
}

func (h *CreditCardHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Name        string `json:"name"`
		CreditLimit string `json:"credit_limit"`
		CloseDay    int    `json:"close_day"`
		DueDay      int    `json:"due_day"`
		IsActive    bool   `json:"is_active"`
	}

	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	creditLimit, err := parseDecimal(req.CreditLimit)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid credit_limit")
		return
	}

	card, err := h.service.Update(r.Context(), id, req.Name, creditLimit, req.CloseDay, req.DueDay, req.IsActive)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, card)
}

func (h *CreditCardHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
