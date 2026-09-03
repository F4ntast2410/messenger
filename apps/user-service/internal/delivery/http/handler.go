package httpHandler

import (
	"context"
	"log/slog"
	"user-service/internal/entity"
)

type Usecase interface {
	ValidateSession(ctx context.Context, sessionID string) (*entity.UserSession, error)
	LoginUser(ctx context.Context, email string, password string) (*entity.LoginResult, error)
	RegisterUser(ctx context.Context, email string, username string, name string, password string) error
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
