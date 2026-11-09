package repository

import (
	"context"
	"database/sql"
	"errors"
	"user-service/internal/entity"
	customErrors "user-service/internal/errors"
	"user-service/internal/repository/postgres/pgmodel"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (s *PostgresStorage) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT password_hash, id, email FROM users WHERE email = $1`
	var user pgmodel.User
	err := s.DB.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.ErrInvalidEmailOrPassword
		}
		return nil, err
	}
	return user.ToEntity(), nil
}

func (s *PostgresStorage) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	query := `SELECT id, email, username, name, created_at FROM users WHERE id = $1`
	var user pgmodel.User
	err := s.DB.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, err
	}
	return user.ToEntity(), nil
}

func (s *PostgresStorage) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*entity.User, error) {
	if len(userIDs) == 0 {
		return []*entity.User{}, nil
	}

	query := `
        SELECT id, email, username, name, created_at
        FROM users
        WHERE id IN (?)
    `

	query, args, err := sqlx.In(query, userIDs)
	if err != nil {
		return nil, err
	}

	query = s.DB.Rebind(query)

	var users []pgmodel.User
	err = s.DB.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.User, 0, len(users))
	for _, user := range users {
		result = append(result, user.ToEntity())
	}

	return result, nil
}
