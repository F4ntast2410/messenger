package pgmodel

import (
	"chat-service/internal/entity"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID  `db:"id"`
	ChatID    uuid.UUID  `db:"chat_id"`
	SenderID  uuid.UUID  `db:"sender_id"`
	Content   string     `db:"content"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" `
	ReplyToID *uuid.UUID `db:"reply_to_id"`
	IsDeleted bool       `db:"is_deleted"`
}

func (m *Message) ToEntity() entity.Message {
	return entity.Message{
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
