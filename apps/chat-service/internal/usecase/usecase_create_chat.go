package usecase

import (
	"context"

	"github.com/google/uuid"
)

func (s *ChatService) CreateChat(ctx context.Context, chatType string, title string, userID uuid.UUID) error {
	return s.repo.CreateChat(ctx, title, chatType, userID)
}
