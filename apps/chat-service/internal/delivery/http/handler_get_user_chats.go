package handler

import (
	"chat-service/internal/delivery/http/dto"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) HandleGetUserChats(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	chats, err := h.ChatService.GetUsersChats(r.Context(), userID)
	if err != nil {
		h.Logger.Error("failed to get user's chats", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	chatsDTOResponse := make([]dto.ChatListItemResponse, 0, 16)
	for _, c := range chats {
		chatsDTOResponse = append(chatsDTOResponse, *dto.ToDTOChatList(&c))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chatsDTOResponse)

}
