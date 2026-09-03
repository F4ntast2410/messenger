package repository

import (
	"context"
	"user-service/internal/entity"
	"user-service/internal/repository/postgres/pgmodel"
)

func (s *PostgresStorage) GetUser(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT password_hash, id, email FROM users WHERE email = $1`
	var user pgmodel.User
	err := s.DB.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return user.ToEntity(), nil
}
