package handler

import (
	"net/http"

	"github.com/thenopholo/my-money/internal/auth"
	handlermw "github.com/thenopholo/my-money/internal/handler/middleware"
	"github.com/thenopholo/my-money/internal/service"
)

type UserHandler struct {
	userService *service.UserService
	jwtManager  *auth.JWTManager
}

func NewUserHandler(userService *service.UserService, jwtManager *auth.JWTManager) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtManager:  jwtManager,
	}
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

	if h.jwtManager == nil {
		Error(w, http.StatusInternalServerError, "jwt manager not configured")
		return
	}

	token, err := h.jwtManager.Generate(user.ID, user.Email)
	if err != nil {
		Error(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	JSON(w, http.StatusOK, map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"token": token,
	})
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, ok := handlermw.UserIDFromContext(r.Context())
	if !ok {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.userService.UpdatePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}
