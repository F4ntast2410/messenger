package usecase

import (
	"context"
	"user-service/internal/entity"

	"github.com/google/uuid"
)

func (uc *UsecaseImpl) SearchUsers(ctx context.Context, query string, cursor *string, limit int, excludeUserID uuid.UUID) ([]entity.User, error) {
	return uc.UserRepo.SearchUsers(ctx, query, cursor, limit, excludeUserID)
}
