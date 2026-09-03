package usecase

import (
	"context"
	"fmt"
	"messenger/pkg/pubsub"
	"time"

	"github.com/google/uuid"
)

func (s *ChatService) UpdateLastReadMessage(ctx context.Context, chatID, userID, messageID uuid.UUID) error {
	err := s.repo.UpdateLastReadMessage(ctx, chatID, userID, messageID)
	if err != nil {
		return err
	}

	// 2. Публикуем событие прочтения в Pub/Sub за O(1)
	event := pubsub.ChatEvent{
		Type:     pubsub.EventMessageRead,
		ChatID:   chatID.String(),
		SenderID: userID.String(), // В данном контексте — кто прочитал
		Payload: map[string]string{
			"message_id": messageID.String(),
		},
		CreatedAt: time.Now(),
	}

	channel := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.ps.Publish(ctx, channel, event); err != nil {
		// Ошибка публикации не должна отменять факт прочтения в БД
		return fmt.Errorf("read status updated, but failed to publish event: %w", err)
	}

	return nil
}
