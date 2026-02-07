package handler

import (
	"net/http"

	"github.com/google/uuid"
	handlermw "github.com/thenopholo/my-money/internal/handler/middleware"
)

func requireUserIDFromContext(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	userID, ok := handlermw.UserIDFromContext(r.Context())
	if !ok {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return uuid.Nil, false
	}

	return userID, true
}
