package dto

import "github.com/google/uuid"

type CreateChatRequest struct {
	Title    string `json:"title"`
	ChatType string `json:"chat_type"`
}

type ChatListItemResponse struct {
	ID uuid.UUID `json:"chat_id"`
}
