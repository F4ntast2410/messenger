package repository

import "context"

func (s *PostgresStorage) RegisterUser(ctx context.Context, email string, username string, passwordHash string, name string) error {
	query := `INSERT INTO users (email, password_hash, name, username) VALUES ($1, $2, $3, $4)`
	_, err := s.DB.ExecContext(ctx, query, email, passwordHash, name, username)
	return err
}
