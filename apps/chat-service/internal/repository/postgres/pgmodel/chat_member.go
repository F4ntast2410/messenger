package pgmodel

import (
	"time"

	"github.com/google/uuid"
)

type ChatMember struct {
	ChatID            uuid.UUID  `db:"chat_id"`
	UserID            uuid.UUID  `db:"user_id"`
	JoinedAt          time.Time  `db:"joined_at"`
	LastReadMessageID *uuid.UUID `db:"last_read_message_id"`
}
