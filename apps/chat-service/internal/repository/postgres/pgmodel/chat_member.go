package pgmodel

import (
	"chat-service/internal/entity"
	"time"

	"github.com/google/uuid"
)

type ChatMember struct {
	ChatID            uuid.UUID  `db:"chat_id"`
	UserID            *uuid.UUID `db:"user_id"`
	JoinedAt          *time.Time `db:"joined_at"`
	LastReadMessageID *uuid.UUID `db:"last_read_message_id"`
	RequestorInChat   bool       `db:"req_in_chat"`
}

func ToEntityChatMember(m *ChatMember) *entity.ChatMember {
	return &entity.ChatMember{
		ChatID:            m.ChatID,
		UserID:            m.UserID,
		JoinedAt:          m.JoinedAt,
		LastReadMessageID: m.LastReadMessageID,
		RequestorInChat:   m.RequestorInChat,
	}
}
