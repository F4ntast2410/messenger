package pubsub

import "time"

type EventType string

const (
	EventNewMessage  EventType = "NEW_MESSAGE"
	EventUserTyping  EventType = "USER_TYPING"
	EventMessageRead EventType = "MESSAGE-READ"
)

type ChatEvent struct {
	Type      EventType `json:"type"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
	Payload   any       `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}
