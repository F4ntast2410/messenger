package usecase

import (
	"context"

	"github.com/google/uuid"
)

func (s *ChatService) JoinToChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	err := s.repo.JoinToChat(ctx, chatID, userID)
	return err
}
