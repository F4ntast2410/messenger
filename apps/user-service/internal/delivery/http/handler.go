package httpHandler

import (
	"context"
	"log/slog"
	"user-service/internal/entity"

	"github.com/google/uuid"
)

type Usecase interface {
	ValidateSession(ctx context.Context, sessionID string) (*entity.UserSession, error)
	LoginUser(ctx context.Context, email string, password string) (*entity.LoginResult, error)
	RegisterUser(ctx context.Context, email string, username string, name string, password string) error
	SearchUsers(ctx context.Context, query string, cursor *string, limit int, excludeUserID uuid.UUID) ([]entity.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	DeleteSession(ctx context.Context, sessionToken string) error
}
type UserHandler struct {
	UC     Usecase
	Logger *slog.Logger
}

func NewUserHandler(uc Usecase, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		UC:     uc,
		Logger: logger,
	}
}
