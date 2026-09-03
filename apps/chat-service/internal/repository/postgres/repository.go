package repository

import "github.com/jmoiron/sqlx"

type PostgresStorage struct {
	DB *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{DB: db}
}
