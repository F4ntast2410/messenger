package pgmodel

import (
	"chat-service/internal/entity"
	"time"

	"github.com/google/uuid"
)

type ChatSummary struct {
	ChatID           uuid.UUID  `db:"chat_id"`
	ChatTitle        *string    `db:"title"`
	ChatType         string     `db:"type"`
	MessageText      *string    `db:"content"`
	MessageSentAt    *time.Time `db:"created_at"`
	MessageUpdatedAt *time.Time `db:"updated_at"`
	UnReadCount      int        `db:"cnt"`
}
type ChatType struct {
	ChatType string `db:"type"`
}

func (c ChatSummary) ToEntity() *entity.ChatSummary {
	return &entity.ChatSummary{
		ChatID:           c.ChatID,
		ChatTitle:        c.ChatTitle,
		ChatType:         c.ChatType,
		MessageText:      c.MessageText,
		MessageSentAt:    c.MessageSentAt,
		MessageUpdatedAt: c.MessageUpdatedAt,
		UnReadCount:      c.UnReadCount,
	}
}
