package usecase

import (
	customErrors "chat-service/internal/errors"
	"context"

	"github.com/google/uuid"
)

func (s *ChatService) CreateChat(ctx context.Context, chatType string, title *string, userID uuid.UUID, userIDs []uuid.UUID) (uuid.UUID, error) {
	userIDs = append(userIDs, userID)
	switch {
	case chatType == "direct" && len(userIDs) == 2:
		chatID, err := s.repo.IsDirectChatExist(ctx, userIDs)
		if err != nil {
			return uuid.Nil, err
		}
		if chatID != uuid.Nil {
			return chatID, customErrors.DirectChatExist
		}
		return s.repo.CreateChat(ctx, nil, chatType, userIDs)
	case chatType == "group":
		return s.repo.CreateChat(ctx, title, chatType, userIDs)
	default:
		return uuid.Nil, customErrors.IncorrectChatType
	}
}
