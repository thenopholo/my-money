package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/service"
)

type PlannedIncomeHandler struct {
	service *service.PlannedIncomeService
}

func NewPlannedIncomeHandler(s *service.PlannedIncomeService) *PlannedIncomeHandler {
	return &PlannedIncomeHandler{service: s}
}

func (h *PlannedIncomeHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	var req struct {
		AccountID   string  `json:"account_id"`
		CategoryID  string  `json:"category_id"`
		Amount      string  `json:"amount"`
		DueDay      int     `json:"due_day"`
		StartDate   *string `json:"start_date"`
		EndDate     *string `json:"end_date"`
		Description string  `json:"description"`
		Frequency   string  `json:"frequency"`
		IsActive    bool    `json:"is_active"`
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

	startDate, endDate, err := parseOptionalDates(req.StartDate, req.EndDate)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	planned, err := h.service.Create(
		r.Context(),
		userID,
		accountID,
		categoryID,
		amount,
		req.DueDay,
		startDate,
		endDate,
		req.Description,
		domain.Recurrence(req.Frequency),
		req.IsActive,
	)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, planned)
}

func (h *PlannedIncomeHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	items, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, items)
}

func (h *PlannedIncomeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	item, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	if item.UserID != userID {
		Error(w, http.StatusForbidden, "forbidden")
		return
	}

	JSON(w, http.StatusOK, item)
}

func (h *PlannedIncomeHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	existing, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	if existing.UserID != userID {
		Error(w, http.StatusForbidden, "forbidden")
		return
	}

	var req struct {
		Amount      string  `json:"amount"`
		DueDay      int     `json:"due_day"`
		StartDate   *string `json:"start_date"`
		EndDate     *string `json:"end_date"`
		Description string  `json:"description"`
		Frequency   string  `json:"frequency"`
		IsActive    bool    `json:"is_active"`
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

	startDate, endDate, err := parseOptionalDates(req.StartDate, req.EndDate)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.service.Update(r.Context(), id, amount, req.DueDay, startDate, endDate, req.Description, domain.Recurrence(req.Frequency), req.IsActive)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, item)
}

func (h *PlannedIncomeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	existing, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	if existing.UserID != userID {
		Error(w, http.StatusForbidden, "forbidden")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseOptionalDates(start, end *string) (*time.Time, *time.Time, error) {
	var startDate *time.Time
	var endDate *time.Time

	if start != nil && *start != "" {
		t, err := parseDate(*start)
		if err != nil {
			return nil, nil, err
		}
		startDate = &t
	}

	if end != nil && *end != "" {
		t, err := parseDate(*end)
		if err != nil {
			return nil, nil, err
		}
		endDate = &t
	}

	return startDate, endDate, nil
}
