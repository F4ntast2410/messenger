package dto

import (
	"chat-service/internal/entity"
	"time"

	"github.com/google/uuid"
)

type SendMessageRequest struct {
	Content string `json:"content"`
}

type MessageResponse struct {
	ID        uuid.UUID  `json:"id"`
	ChatID    uuid.UUID  `json:"chat_id"`
	SenderID  uuid.UUID  `json:"sender_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	ReplyToID *uuid.UUID `json:"reply_to_id"`
	IsDeleted bool       `json:"is_deleted"`
}

func ToResponse(m *entity.Message) MessageResponse {
	return MessageResponse{
		ID:        m.ID,
		ChatID:    m.ChatID,
		SenderID:  m.SenderID,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		ReplyToID: m.ReplyToID,
		IsDeleted: m.IsDeleted,
	}
}

type MessageIDRequest struct {
	MessageID uuid.UUID `json:"message_id"`
}
