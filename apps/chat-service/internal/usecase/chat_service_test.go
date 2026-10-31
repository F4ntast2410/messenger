package usecase

import (
	"chat-service/internal/entity"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type FakeChatRepository struct {
	SaveMessageFunc           func(ctx context.Context, msg *entity.Message) error
	GetChatMemberIDsFunc      func(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error)
	GetMessagesPaginatedFunc  func(ctx context.Context, chatID uuid.UUID, lastSeenID *uuid.UUID, limit int) ([]entity.Message, error)
	CreateChatFunc            func(ctx context.Context, title *string, chatType string, userIDs []uuid.UUID) (uuid.UUID, error)
	JoinToChatFunc            func(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	GetLastReadMessageFunc    func(ctx context.Context, userID, chatID uuid.UUID) (*uuid.UUID, error)
	UpdateLastReadMessageFunc func(ctx context.Context, chatID, userID, messageID uuid.UUID) error
	GetUsersChatsFunc         func(ctx context.Context, userID uuid.UUID) ([]entity.ChatSummary, error)
	GetChatTypeFunc           func(ctx context.Context, chatID uuid.UUID) (string, error)
	IsDirectChatExistFunc     func(ctx context.Context, userIDs []uuid.UUID) (uuid.UUID, error)
}

func (r *FakeChatRepository) GetMessagesPaginated(ctx context.Context, chatID uuid.UUID, lastSeenID *uuid.UUID, limit int) ([]entity.Message, error) {
	return r.GetMessagesPaginatedFunc(ctx, chatID, lastSeenID, limit)
}
func (r *FakeChatRepository) SaveMessage(ctx context.Context, msg *entity.Message) error {
	return r.SaveMessageFunc(ctx, msg)
}
func (r *FakeChatRepository) GetChatMemberIDs(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	return r.GetChatMemberIDsFunc(ctx, chatID)
}
func (r *FakeChatRepository) CreateChat(ctx context.Context, title *string, chatType string, userIDs []uuid.UUID) (uuid.UUID, error) {
	return r.CreateChatFunc(ctx, title, chatType, userIDs)
}
func (r *FakeChatRepository) JoinToChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	return r.JoinToChatFunc(ctx, chatID, userID)
}
func (r *FakeChatRepository) GetLastReadMessage(ctx context.Context, userID, chatID uuid.UUID) (*uuid.UUID, error) {
	return r.GetLastReadMessageFunc(ctx, userID, chatID)
}
func (r *FakeChatRepository) UpdateLastReadMessage(ctx context.Context, chatID, userID, messageID uuid.UUID) error {
	return r.UpdateLastReadMessageFunc(ctx, chatID, userID, messageID)
}
func (r *FakeChatRepository) GetUsersChats(ctx context.Context, userID uuid.UUID) ([]entity.ChatSummary, error) {
	return r.GetUsersChatsFunc(ctx, userID)
}
func (r *FakeChatRepository) GetChatType(ctx context.Context, chatID uuid.UUID) (string, error) {
	return r.GetChatTypeFunc(ctx, chatID)
}
func (r *FakeChatRepository) IsDirectChatExist(ctx context.Context, userIDs []uuid.UUID) (uuid.UUID, error) {
	return r.IsDirectChatExistFunc(ctx, userIDs)
}

var _ ChatRepository = (*FakeChatRepository)(nil)

func TestChatService_GetMessagesPaginated_Success(t *testing.T) {
	messages := []entity.Message{{ID: uuid.New(), ChatID: uuid.New(), SenderID: uuid.New(), Content: "", CreatedAt: time.Now(), UpdatedAt: nil, ReplyToID: nil, IsDeleted: false}, {ID: uuid.New(), ChatID: uuid.New(), SenderID: uuid.New(), Content: "", CreatedAt: time.Now(), UpdatedAt: nil, ReplyToID: nil, IsDeleted: false}}
	fakeRepo := FakeChatRepository{GetMessagesPaginatedFunc: func(ctx context.Context, chatID uuid.UUID, lastSeenID *uuid.UUID, limit int) ([]entity.Message, error) {
		return messages, nil
	}}
	chatService := NewChatService(&fakeRepo, nil)
	tests := []struct {
		name       string
		chatID     uuid.UUID
		lastSeenID *uuid.UUID
		limit      int
		wantErr    bool
	}{
		{
			name:       "Успешное выполнение",
			chatID:     uuid.New(),
			lastSeenID: nil,
			limit:      10,
			wantErr:    false,
		},
		{
			name:       "Лимит меньше нуля",
			chatID:     uuid.New(),
			lastSeenID: nil,
			limit:      -1,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := chatService.GetMessagesPaginated(context.Background(), tt.chatID, tt.lastSeenID, tt.limit)
			if tt.wantErr {
				assert.Error(t, err, "GetMessagesPaginated() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				assert.Equal(t, messages, result, "GetMessagesPaginated() result = %v, expected = %v", result, messages)
				assert.NoError(t, err, "GetMessagesPaginated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
