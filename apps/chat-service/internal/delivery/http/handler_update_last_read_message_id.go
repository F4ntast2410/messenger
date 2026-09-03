package handler

import (
	"chat-service/internal/delivery/http/dto"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) HandleUpdateLastReadMessageID(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	senderID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	chatIDStr := r.PathValue("id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat_id", http.StatusBadRequest)
		return
	}
	var req dto.MessageIDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Debug("failed to decode json", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ChatService.UpdateLastReadMessage(r.Context(), chatID, senderID, req.MessageID)
	if err != nil {
		h.Logger.Error("failed to update last read message id", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
