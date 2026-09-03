package handler

import (
	"chat-service/internal/delivery/http/dto"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (h *Handler) HandleGetMessagesPaginated(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.PathValue("id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat_id", http.StatusBadRequest)
		return
	}
	queryParams := r.URL.Query()

	beforeIDStr := queryParams.Get("before_id")
	limitStr := queryParams.Get("limit")

	limit := 30
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	var beforeID *uuid.UUID
	if beforeIDStr != "" {
		parsedUUID, err := uuid.Parse(beforeIDStr)
		if err != nil {
			http.Error(w, "невалидный формат before_id", http.StatusBadRequest)
			return
		}
		beforeID = &parsedUUID
	}
	messages, err := h.ChatService.GetMessagesPaginated(r.Context(), chatID, beforeID, limit)
	if err != nil {
		http.Error(w, "failed to get messages", http.StatusInternalServerError)
		return
	}
	messagesResponse := make([]dto.MessageResponse, len(messages))
	for i := range messages {
		messagesResponse[i] = dto.ToResponse(&messages[i])
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messagesResponse)
}
