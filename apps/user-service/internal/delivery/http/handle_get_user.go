package httpHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"user-service/internal/delivery/http/dto"

	"github.com/google/uuid"
)

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := h.UC.GetUserByID(r.Context(), userID)
	if err != nil {
		h.Logger.Error("failed to get user by id", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	meResponse := dto.ToDTOMeResponse(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(meResponse)
}
