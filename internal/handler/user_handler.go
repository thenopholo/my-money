package handler

import (
	"net/http"

	"github.com/thenopholo/my-money/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	JSON(w, http.StatusOK, map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"token": "jwt_token_aqui",
	})
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID          string `json:"user_id"`
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
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

	if err := h.userService.UpdatePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}
