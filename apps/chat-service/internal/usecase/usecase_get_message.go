package usecase

import (
	"chat-service/internal/entity"
	"context"
	"errors"

	"github.com/google/uuid"
)

func (s *ChatService) GetMessagesPaginated(ctx context.Context, chatID uuid.UUID, lastSeenID *uuid.UUID, limit int) ([]entity.Message, error) {
	if limit <= 0 {
		return nil, errors.New("limit is less than or equal to zero")
	}
	messages, err := s.repo.GetMessagesPaginated(ctx, chatID, lastSeenID, limit)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
