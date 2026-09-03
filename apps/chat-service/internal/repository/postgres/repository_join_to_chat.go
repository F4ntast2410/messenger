package repository

import (
	"context"

	"github.com/google/uuid"
)

func (s *PostgresStorage) JoinToChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	query := `
		INSERT INTO chat_members (chat_id, user_id)
		VALUES ($1, $2)
	`
	_, err := s.DB.ExecContext(ctx, query, chatID, userID)
	return err
}
