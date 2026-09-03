package repository

import (
	"chat-service/internal/entity"
	"context"

	"github.com/google/uuid"
)

func (s *PostgresStorage) GetUsersChats(ctx context.Context, userID uuid.UUID) ([]entity.Chat, error) {

}
