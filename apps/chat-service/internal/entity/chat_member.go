package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChatMember struct {
	ChatID            uuid.UUID
	UserID            *uuid.UUID
	JoinedAt          *time.Time
	LastReadMessageID *uuid.UUID
	RequestorInChat   bool
}
