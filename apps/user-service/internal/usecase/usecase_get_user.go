package usecase

import (
	"context"
	"user-service/internal/entity"

	"github.com/google/uuid"
)

func (uc *UsecaseImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return uc.UserRepo.GetUserByID(ctx, userID)
}

func (uc *UsecaseImpl) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*entity.User, error) {
	return uc.UserRepo.GetUsersByIDs(ctx, userIDs)
}
