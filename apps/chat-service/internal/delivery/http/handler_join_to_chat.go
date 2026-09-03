package handler

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) HandleJoinToChat(w http.ResponseWriter, r *http.Request) {
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
	err = h.ChatService.JoinToChat(r.Context(), chatID, senderID)
	if err != nil {
		h.Logger.Error("failed to join to chat", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
