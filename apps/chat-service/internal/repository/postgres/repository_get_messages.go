package repository

import (
	"chat-service/internal/entity"
	"chat-service/internal/repository/postgres/pgmodel"
	"context"

	"github.com/google/uuid"
)

func (s *PostgresStorage) GetMessagesPaginated(ctx context.Context, chatID uuid.UUID, lastSeenID *uuid.UUID, limit int) ([]entity.Message, error) {
	var query string
	var args []interface{}

	if lastSeenID == nil {
		// Первичная загрузка чата (самые свежие сообщения)
		query = `
			SELECT id, chat_id, sender_id, content, created_at 
			FROM messages 
			WHERE chat_id = $1 AND is_deleted = FALSE 
			ORDER BY id DESC LIMIT $2`
		args = []interface{}{chatID, limit}
	} else {
		// Подгрузка старых сообщений при скролле вверх (id < last_seen_id)
		query = `
			SELECT id, chat_id, sender_id, content, created_at 
			FROM messages 
			WHERE chat_id = $1 AND id < $2 AND is_deleted = FALSE 
			ORDER BY id DESC LIMIT $3`
		args = []interface{}{chatID, *lastSeenID, limit}
	}
	var messages []pgmodel.Message
	err := s.DB.SelectContext(ctx, &messages, query, args)
	if err != nil {
		return nil, err
	}
	entityMessages := make([]entity.Message, len(messages))
	for i := range messages {
		entityMessages[i] = messages[i].ToEntity()
	}
	return entityMessages, nil
}
