package repository

import (
	"context"

	"github.com/google/uuid"
)

func (s *PostgresStorage) GetChatMemberIDs(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT user_id FROM chat_members WHERE chat_id = $1`
	var memberIDs []uuid.UUID
	err := s.DB.SelectContext(ctx, &memberIDs, query, chatID)
	if err != nil {
		return nil, err
	}
	return memberIDs, nil
}
