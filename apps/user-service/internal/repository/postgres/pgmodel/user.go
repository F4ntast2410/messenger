package pgmodel

import (
	"time"
	"user-service/internal/entity"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

func (u *User) ToEntity() *entity.User {
	return &entity.User{
		ID:           u.ID,
		Name:         u.Name,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
	}
}
