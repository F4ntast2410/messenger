package dto

import (
	"chat-service/internal/entity"
	"time"

	"github.com/google/uuid"
)

type CreateChatRequest struct {
	Title    *string     `json:"title"`
	ChatType string      `json:"chat_type"`
	UserIDs  []uuid.UUID `json:"user_ids"`
}
type CreateChatResponse struct {
	ChatID uuid.UUID `json:"chat_id"`
}
type ChatListItemResponse struct {
	ChatID           uuid.UUID  `json:"chat_id"`
	ChatTitle        *string    `json:"chat_title"`
	ChatType         string     `json:"chat_type"`
	MessageText      *string    `json:"msg_text"`
	MessageSentAt    *time.Time `json:"msg_sent_at"`
	MessageUpdatedAt *time.Time `json:"msg_updated_at"`
	UnReadCount      int        `json:"unread_msg_cnt"`
}

func ToDTOChatList(c *entity.ChatSummary) *ChatListItemResponse {
	return &ChatListItemResponse{
		ChatID:           c.ChatID,
		ChatTitle:        c.ChatTitle,
		ChatType:         c.ChatType,
		MessageText:      c.MessageText,
		MessageSentAt:    c.MessageSentAt,
		MessageUpdatedAt: c.MessageUpdatedAt,
		UnReadCount:      c.UnReadCount,
	}
}
