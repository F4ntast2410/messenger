package repository

import (
	"context"

	"github.com/google/uuid"
)

func (s *PostgresStorage) UpdateLastReadMessage(ctx context.Context, chatID, userID, messageID uuid.UUID) error {
	query := `UPDATE chat_members cm
	SET last_read_message_id = m.id
	FROM messages m
	WHERE m.id = $1
	AND m.chat_id = $2
	AND cm.chat_id = $2
	AND cm.user_id = $3;`

	// Используем ExecContext, так как нам не нужно сканировать возвращаемые строки
	result, err := s.DB.ExecContext(
		ctx,
		query,
		messageID,
		chatID,
		userID,
	)
	if err != nil {
		return err
	}

	return checkAffected(result, "chat member not found")
}
