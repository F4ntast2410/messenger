package repository

import (
	"chat-service/internal/entity"
	"chat-service/internal/repository/postgres/pgmodel"
	"context"

	"github.com/google/uuid"
)

func (s *PostgresStorage) GetUsersChats(ctx context.Context, userID uuid.UUID) ([]entity.ChatSummary, error) {
	query := `SELECT last_msg.content, last_msg.created_at, last_msg.updated_at, last_msg.sender_id, cm.chat_id, c.title, c.type, msg_cnt.cnt
	FROM chat_members cm 
	LEFT JOIN LATERAL 
	(SELECT m.content, m.created_at, m.updated_at, m.sender_id
	FROM messages m 
	WHERE m.chat_id = cm.chat_id 
	ORDER BY m.created_at DESC 
	LIMIT 1) last_msg ON true 
	JOIN chats c ON c.id = cm.chat_id 
	JOIN LATERAL 
	(SELECT COUNT(m.id) AS cnt
	FROM messages m 
	WHERE m.chat_id = cm.chat_id AND (cm.last_read_message_id IS NULL OR m.id > cm.last_read_message_id)) msg_cnt ON true 
	WHERE cm.user_id = $1
	ORDER BY last_msg.created_at DESC NULLS LAST`
	chatSummary := make([]pgmodel.ChatSummary, 0, 16)
	err := s.DB.SelectContext(ctx, &chatSummary, query, userID)
	if err != nil {
		return nil, err
	}
	chatSummaryEntity := make([]entity.ChatSummary, 0, 16)
	for _, c := range chatSummary {
		chatSummaryEntity = append(chatSummaryEntity, *c.ToEntity())
	}
	return chatSummaryEntity, nil
}
