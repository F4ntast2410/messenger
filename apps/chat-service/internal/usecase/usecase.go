package usecase

import (
	"chat-service/internal/entity"
	"context"
	"messenger/pkg/pubsub"

	"github.com/google/uuid"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, msg *entity.Message) error
	GetChatMemberIDs(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error)
	GetMessagesPaginated(ctx context.Context, chatID uuid.UUID, lastSeenID *uuid.UUID, limit int) ([]entity.Message, error)
	CreateChat(ctx context.Context, title *string, chatType string, userIDs []uuid.UUID) (uuid.UUID, error)
	JoinToChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	GetLastReadMessage(ctx context.Context, userID, chatID uuid.UUID) (*uuid.UUID, error)
	UpdateLastReadMessage(ctx context.Context, chatID, userID, messageID uuid.UUID) error
	GetUsersChats(ctx context.Context, userID uuid.UUID) ([]entity.ChatSummary, error)
	IsDirectChatExist(ctx context.Context, userIDs []uuid.UUID) (uuid.UUID, error)
	GetChatType(ctx context.Context, chatID uuid.UUID) (string, error)
}

type ChatService struct {
	repo ChatRepository
	ps   *pubsub.PubSub
}

func NewChatService(repo ChatRepository, ps *pubsub.PubSub) *ChatService {
	return &ChatService{
		repo: repo,
		ps:   ps,
	}
}
