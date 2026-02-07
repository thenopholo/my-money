package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/service"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID       string  `json:"user_id"`
		Name         string  `json:"name"`
		CategoryType string  `json:"category_type"`
		Color        *string `json:"color"`
		Icon         *string `json:"icon"`
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

	category, err := h.service.Create(r.Context(), userID, domain.CategoryType(req.CategoryType), req.Name, req.Color, req.Icon)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, category)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUUID(r.URL.Query().Get("user_id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	categories, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	category, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Name         string  `json:"name"`
		CategoryType string  `json:"category_type"`
		Color        *string `json:"color"`
		Icon         *string `json:"icon"`
	}

	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	category, err := h.service.Update(r.Context(), id, domain.CategoryType(req.CategoryType), req.Name, req.Color, req.Icon)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
