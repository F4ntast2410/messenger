package repository

import "context"

func (s *PostgresStorage) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := s.DB.GetContext(ctx, &exists, query, email)
	if err != nil {
		return false, err
	}

	return exists, nil
}
