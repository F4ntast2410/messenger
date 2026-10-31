package handler

import (
	"chat-service/internal/delivery/http/dto"
	customErrors "chat-service/internal/errors"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) HandleCreateChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	chatID, err := h.ChatService.CreateChat(r.Context(), req.ChatType, req.Title, senderID, req.UserIDs)
	resp := dto.CreateChatResponse{ChatID: chatID}
	if err != nil {
		if errors.Is(err, customErrors.DirectChatExist) {
			h.Logger.Info("direct chat already exists")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(resp)
			return
		}
		if errors.Is(err, customErrors.IncorrectChatType) {
			h.Logger.Info("incorrect chat type")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.Logger.Error("failed to create chat", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
