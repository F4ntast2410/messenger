package pubsub

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventNewMessage  EventType = "NEW_MESSAGE"
	EventUserTyping  EventType = "USER_TYPING"
	EventMessageRead EventType = "MESSAGE-READ"
)

type MessagePayload struct {
	ID        uuid.UUID `json:"id"`
	ChatID    uuid.UUID `json:"chat_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatEvent struct {
	Type      EventType `json:"type"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
	Payload   any       `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type TypeEvent struct {
	Type EventType `json:"type"`
}
