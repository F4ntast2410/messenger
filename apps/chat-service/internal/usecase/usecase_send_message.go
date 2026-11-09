package usecase

import (
	"chat-service/internal/entity"
	"context"
	"fmt"
	"messenger/pkg/pubsub"
	"time"

	"github.com/google/uuid"
)

func ToPayload(m *entity.Message) pubsub.MessagePayload {
	return pubsub.MessagePayload{
		ID:        m.ID,
		ChatID:    m.ChatID,
		SenderID:  m.SenderID,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
	}
}
func (s *ChatService) SendMessage(ctx context.Context, chatID, senderID uuid.UUID, replyID *uuid.UUID, content string) (*entity.Message, error) {
	// 1. Генерируем UUIDv7 (или используем uuid.NewV7())
	msgID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	msg := &entity.Message{
		ID:        msgID,
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		CreatedAt: time.Now(),
		ReplyToID: replyID,
	}

	// 2. Атомарно сохраняем в БД
	if err := s.repo.SaveMessage(ctx, msg); err != nil {
		return nil, err
	}

	// 3. Публикуем событие в Redis Pub/Sub за O(1)
	// Вся логика фильтрации и разогрева WebSocket переходит на сторону ws-gateway
	event := pubsub.ChatEvent{
		Type:      pubsub.EventNewMessage,
		ChatID:    chatID.String(),
		SenderID:  senderID.String(),
		Payload:   ToPayload(msg),
		CreatedAt: msg.CreatedAt,
	}

	channel := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.ps.Publish(ctx, channel, event); err != nil {
		return msg, fmt.Errorf("message saved, but failed to publish event: %w", err)
	}

	return msg, nil
}
