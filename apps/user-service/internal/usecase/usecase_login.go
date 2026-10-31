package usecase

import (
	"context"
	"user-service/internal/entity"
	customErrors "user-service/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

func (uc *UsecaseImpl) LoginUser(ctx context.Context, email string, password string) (*entity.LoginResult, error) {
	user, err := uc.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, customErrors.ErrInvalidEmailOrPassword
	}
	sessionID, err := uc.CreateSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &entity.LoginResult{
		SessionToken: sessionID,
		Required2FA:  false,
	}, nil
}
