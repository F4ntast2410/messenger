package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func (s *PostgresStorage) IsDirectChatExist(ctx context.Context, userIDs []uuid.UUID) (uuid.UUID, error) {
	query := `SELECT c.id FROM chats c
		JOIN chat_members cm1 ON cm1.chat_id = c.id AND cm1.user_id = $1
		JOIN chat_members cm2 ON cm2.chat_id = c.id AND cm2.user_id = $2
		WHERE c.type = 'direct'`
	var chatID uuid.UUID
	err := s.DB.GetContext(ctx, &chatID, query, userIDs[0], userIDs[1])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}

	return chatID, nil
}
