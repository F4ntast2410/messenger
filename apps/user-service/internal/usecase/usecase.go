package usecase

import (
	"context"
	"time"
	"user-service/internal/entity"

	"github.com/google/uuid"
)

type AuthRepository interface {
	GetUserID(ctx context.Context, sessionID string) (uuid.UUID, error)
	SaveSessionToken(ctx context.Context, sessionID string, userID uuid.UUID, time_duration time.Duration) error
	DeleteSessionToken(ctx context.Context, sessionID string) error
}

type UsecaseImpl struct {
	AuthRepo AuthRepository
	UserRepo UserRepository
}

func NewUsecase(authRepo AuthRepository, userRepo UserRepository) *UsecaseImpl {
	return &UsecaseImpl{AuthRepo: authRepo, UserRepo: userRepo}
}

type UserRepository interface {
	RegisterUser(ctx context.Context, email string, username string, passwordHash string, name string) error
	UserExists(ctx context.Context, email string) (bool, error)
	GetUser(ctx context.Context, email string) (*entity.User, error)
}
