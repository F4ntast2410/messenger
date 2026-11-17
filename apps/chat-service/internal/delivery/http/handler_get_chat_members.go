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

func (h *Handler) HandleGetChatMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	chatIDStr := r.PathValue("id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: "Invalid chat ID",
		})
		return
	}
	requestorIDStr := r.Header.Get("X-User-ID")
	requestorID, err := uuid.Parse(requestorIDStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	chatMembers, err := h.ChatService.GetChatMembers(r.Context(), chatID, requestorID)
	if err != nil {
		if errors.Is(err, customErrors.NoPermission) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Requestor does not have permission to view users list in this chat"})
			return
		} else if errors.Is(err, customErrors.ChatNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Chat not found"})
			return
		} else if errors.Is(err, customErrors.NoUsersInChat) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]dto.ChatMemberResponse{})
			return
		}
		h.Logger.Error("failed to get chat_members (get_chat_members)", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var dtoChatMembers []dto.ChatMemberResponse
	for _, m := range chatMembers {
		dtoChatMembers = append(dtoChatMembers, *dto.ToDTOChatMember(&m))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtoChatMembers)
}
