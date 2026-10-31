package usecase

import (
	customErrors "chat-service/internal/errors"
	"context"

	"github.com/google/uuid"
)

func (s *ChatService) JoinToChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	err := s.repo.JoinToChat(ctx, chatID, userID)
	return err
}
func (s *ChatService) CheckChatType(ctx context.Context, chatID uuid.UUID) error {
	chatType, err := s.repo.GetChatType(ctx, chatID)
	if err != nil {
		return err
	}
	if chatType != "group" {
		return customErrors.CannotJoinDirectChat
	}
	return nil
}
