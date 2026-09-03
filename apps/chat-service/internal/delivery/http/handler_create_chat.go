package handler

import (
	"chat-service/internal/delivery/http/dto"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) HandleCreateChat(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	senderID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var req dto.CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ChatService.CreateChat(r.Context(), req.ChatType, req.Title, senderID)
	if err != nil {
		h.Logger.Error("failed to create chat", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}
