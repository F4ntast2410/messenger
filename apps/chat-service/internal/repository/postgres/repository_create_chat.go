package repository

import (
	"context"

	"github.com/google/uuid"
)

func (s *PostgresStorage) CreateChat(ctx context.Context, title *string, chatType string, userIDs []uuid.UUID) (uuid.UUID, error) {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()
	queryCreate := `INSERT INTO chats (title, type) VALUES ($1, $2) RETURNING id`
	var chatID uuid.UUID
	if err := tx.GetContext(ctx, &chatID, queryCreate, title, chatType); err != nil {
		return uuid.Nil, err
	}
	for _, u := range userIDs {
		queryJoin := `INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2)`
		if _, err := tx.ExecContext(ctx, queryJoin, chatID, u); err != nil {
			return uuid.Nil, err
		}
	}

	return chatID, tx.Commit()
}
