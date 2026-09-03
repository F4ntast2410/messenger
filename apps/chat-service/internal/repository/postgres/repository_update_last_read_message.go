package repository

import (
	"context"

	"github.com/google/uuid"
)

func (s *PostgresStorage) UpdateLastReadMessage(ctx context.Context, chatID, userID, messageID uuid.UUID) error {
	query := `UPDATE chat_members SET last_read_message_id = $1 WHERE chat_id = $2 AND user_id = $3`

	// Используем ExecContext, так как нам не нужно сканировать возвращаемые строки
	result, err := s.DB.ExecContext(ctx, query, chatID, userID, messageID)
	if err != nil {
		return err
	}
	return checkAffected(result, "chat member not found")
}
