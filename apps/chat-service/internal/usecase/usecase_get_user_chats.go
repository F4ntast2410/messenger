package usecase

import (
	"chat-service/internal/entity"
	"context"

	"github.com/google/uuid"
)

func (s *ChatService) GetUsersChats(ctx context.Context, userID uuid.UUID) ([]entity.ChatSummary, error) {
	return s.repo.GetUsersChats(ctx, userID)
}
