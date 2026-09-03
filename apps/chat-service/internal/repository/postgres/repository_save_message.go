package repository

import (
	"chat-service/internal/entity"
	"context"
)

func (s *PostgresStorage) SaveMessage(ctx context.Context, msg *entity.Message) error {
	query := `
		INSERT INTO messages (id, chat_id, sender_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := s.DB.ExecContext(ctx, query, msg.ID, msg.ChatID, msg.SenderID, msg.Content, msg.CreatedAt)
	return err
}
