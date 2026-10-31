package repository

import (
	customErrors "chat-service/internal/errors"
	"chat-service/internal/repository/postgres/pgmodel"
	"context"
	"database/sql"
	"errors"

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

func (s *PostgresStorage) GetChatType(ctx context.Context, chatID uuid.UUID) (string, error) {
	query := `
		SELECT type FROM chats WHERE id = $1
	`
	var chat pgmodel.ChatType
	err := s.DB.GetContext(ctx, &chat, query, chatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", customErrors.ChatNotFound
		}
		return "", err
	}
	return chat.ChatType, nil
}
