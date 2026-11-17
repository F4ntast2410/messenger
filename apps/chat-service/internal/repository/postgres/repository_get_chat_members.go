package repository

import (
	"chat-service/internal/entity"
	"chat-service/internal/repository/postgres/pgmodel"
	"context"

	"github.com/google/uuid"
)

func (s *PostgresStorage) GetChatMembers(ctx context.Context, chatID uuid.UUID, requestorID uuid.UUID) ([]entity.ChatMember, error) {
	query := `SELECT cm.user_id, cm.last_read_message_id, EXISTS(SELECT 1 FROM chat_members rq WHERE rq.chat_id = $1 AND rq.user_id = $2) AS req_in_chat FROM chats c LEFT JOIN chat_members cm ON c.id = cm.chat_id WHERE c.id = $1`
	var chatMembers []pgmodel.ChatMember
	if err := s.DB.SelectContext(ctx, &chatMembers, query, chatID, requestorID); err != nil {
		return nil, err
	}
	var entityChatMembers []entity.ChatMember
	for _, m := range chatMembers {
		entityChatMembers = append(entityChatMembers, *pgmodel.ToEntityChatMember(&m))
	}
	return entityChatMembers, nil

}
