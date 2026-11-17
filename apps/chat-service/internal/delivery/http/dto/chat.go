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
	ChatID            uuid.UUID  `json:"chat_id"`
	ChatTitle         *string    `json:"chat_title"`
	ChatType          string     `json:"chat_type"`
	MessageText       *string    `json:"msg_text"`
	MessageSentAt     *time.Time `json:"msg_sent_at"`
	MessageUpdatedAt  *time.Time `json:"msg_updated_at"`
	MessageSenderID   *uuid.UUID `json:"msg_sender_id"`
	MessageID         *uuid.UUID `json:"msg_id"`
	LastReadMessageID *uuid.UUID `json:"last_read_message_id"`
	UnReadCount       int        `json:"unread_msg_cnt"`
}
type ChatMemberResponse struct {
	UserID            uuid.UUID  `json:"user_id"`
	LastReadMessageID *uuid.UUID `json:"last_read_message_id"`
}
type ErrorResponse struct {
	Error string
}

func ToDTOChatMember(m *entity.ChatMember) *ChatMemberResponse {
	return &ChatMemberResponse{
		UserID:            *m.UserID,
		LastReadMessageID: m.LastReadMessageID,
	}
}

func ToDTOChatList(c *entity.ChatSummary) *ChatListItemResponse {
	return &ChatListItemResponse{
		ChatID:            c.ChatID,
		ChatTitle:         c.ChatTitle,
		ChatType:          c.ChatType,
		MessageText:       c.MessageText,
		MessageSentAt:     c.MessageSentAt,
		MessageUpdatedAt:  c.MessageUpdatedAt,
		MessageSenderID:   c.MessageSenderID,
		MessageID:         c.MessageID,
		LastReadMessageID: c.LastReadMessageID,
		UnReadCount:       c.UnReadCount,
	}
}
