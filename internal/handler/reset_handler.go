package handler

import (
	"net/http"

	"github.com/thenopholo/my-money/internal/service"
)

type ResetHandler struct {
	service *service.ResetService
}

func NewResetHandler(s *service.ResetService) *ResetHandler {
	return &ResetHandler{service: s}
}

func (h *ResetHandler) ResetAllTransactions(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	if err := h.service.ResetAllTransactions(r.Context(), userID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
