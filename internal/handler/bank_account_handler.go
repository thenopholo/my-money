package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/service"
)

type BankAccountHandler struct {
	service *service.BankAccountService
}

func NewBankAccountHandler(s *service.BankAccountService) *BankAccountHandler {
	return &BankAccountHandler{service: s}
}

func (h *BankAccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID      string `json:"user_id"`
		AccountType string `json:"account_type"`
		Name        string `json:"name"`
		BankName    string `json:"bank_name"`
		Balance     string `json:"balance"`
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

	balance, err := parseDecimal(req.Balance)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid balance")
		return
	}

	account, err := h.service.Create(
		r.Context(),
		userID,
		domain.AccountType(req.AccountType),
		req.Name,
		req.BankName,
		balance,
	)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, account)
}

func (h *BankAccountHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUUID(r.URL.Query().Get("user_id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	accounts, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, accounts)
}

func (h *BankAccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	account, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, account)
}

func (h *BankAccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Name     string `json:"name"`
		BankName string `json:"bank_name"`
		Balance  string `json:"balance"`
		IsActive bool   `json:"is_active"`
	}

	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	balance, err := parseDecimal(req.Balance)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid balance")
		return
	}

	account, err := h.service.Update(r.Context(), id, req.Name, req.BankName, balance, req.IsActive)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, account)
}

func (h *BankAccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
