package repository

import (
	"context"
	"user-service/internal/entity"
	"user-service/internal/repository/postgres/pgmodel"

	"github.com/google/uuid"
)

func (s *PostgresStorage) SearchUsers(ctx context.Context, query string, cursor *string, limit int, excludeUserID uuid.UUID) ([]entity.User, error) {
	sqlQuery := `
	SELECT id, name, username FROM users
	WHERE (username ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
	AND id != $2
	AND ($3::text IS NULL OR username > $3)
	ORDER BY username ASC
	LIMIT $4
	`
	users := make([]pgmodel.User, 0, 10)
	err := s.DB.SelectContext(ctx, &users, sqlQuery, query, excludeUserID, cursor, limit)
	if err != nil {
		return nil, err
	}
	usersEntity := make([]entity.User, 0, 10)
	for _, u := range users {
		usersEntity = append(usersEntity, *u.ToEntity())
	}
	return usersEntity, nil

}
