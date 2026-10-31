package entity

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID
	ChatType  string
	Title     *string
	CreatedAt time.Time
}

type ChatSummary struct {
	ChatID           uuid.UUID
	ChatTitle        *string
	ChatType         string
	MessageText      *string
	MessageSentAt    *time.Time
	MessageUpdatedAt *time.Time
	UnReadCount      int
}
