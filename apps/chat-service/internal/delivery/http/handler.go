package handler

import (
	"chat-service/internal/entity"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type ChatService interface {
	JoinToChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	CreateChat(ctx context.Context, chatType string, title *string, userID uuid.UUID, userIDs []uuid.UUID) (uuid.UUID, error)
	SendMessage(ctx context.Context, chatID, senderID uuid.UUID, replyID *uuid.UUID, content string) (*entity.Message, error)
	GetMessagesPaginated(ctx context.Context, chatID uuid.UUID, lastSeenID *uuid.UUID, limit int) ([]entity.Message, error)
	UpdateLastReadMessage(ctx context.Context, chatID, userID, messageID uuid.UUID) error
	GetUsersChats(ctx context.Context, userID uuid.UUID) ([]entity.ChatSummary, error)
	CheckChatType(ctx context.Context, chatID uuid.UUID) error
}

type Handler struct {
	ChatService ChatService
	Logger      *slog.Logger
}

func NewHandler(chatService ChatService, logger *slog.Logger) *Handler {
	return &Handler{ChatService: chatService, Logger: logger}
}
