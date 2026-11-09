package repository

import (
	"chat-service/internal/entity"
	"context"
)

func (s *PostgresStorage) SaveMessage(ctx context.Context, msg *entity.Message) error {
	query := `
		INSERT INTO messages (id, chat_id, sender_id, content, created_at, reply_to_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.DB.ExecContext(ctx, query, msg.ID, msg.ChatID, msg.SenderID, msg.Content, msg.CreatedAt, msg.ReplyToID)
	return err
}
