package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt *time.Time
	ReplyToID *uuid.UUID
	IsDeleted bool
}
