package repository

import (
	"context"

	"github.com/google/uuid"
)

func (s *PostgresStorage) GetLastReadMessage(ctx context.Context, userID, chatID uuid.UUID) (*uuid.UUID, error) {
	query := `
			SELECT last_read_message_id
			FROM chat_members 
			WHERE chat_id = $1 AND user_id = $2`
	var lastReadMessage *uuid.UUID
	err := s.DB.GetContext(ctx, &lastReadMessage, query, chatID, userID)
	if err != nil {
		return nil, err
	}
	return lastReadMessage, nil
}
