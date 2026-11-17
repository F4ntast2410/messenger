package usecase

import (
	"chat-service/internal/entity"
	customErrors "chat-service/internal/errors"
	"context"

	"github.com/google/uuid"
)

func (s *ChatService) GetChatMembers(ctx context.Context, chatID uuid.UUID, requestorID uuid.UUID) ([]entity.ChatMember, error) {
	chatMembers, err := s.repo.GetChatMembers(ctx, chatID, requestorID)
	if err != nil {
		return nil, err
	}
	if len(chatMembers) == 0 {
		return nil, customErrors.ChatNotFound
	} else if !chatMembers[0].RequestorInChat {
		if chatMembers[0].UserID == nil {
			return nil, customErrors.NoUsersInChat
		}
		return nil, customErrors.NoPermission
	}
	return chatMembers, nil
}
