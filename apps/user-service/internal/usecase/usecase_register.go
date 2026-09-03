package usecase

import (
	"context"
	customErrors "user-service/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

func (uc *UsecaseImpl) RegisterUser(ctx context.Context, email string, username string, name string, password string) error {
	exists, err := uc.UserRepo.UserExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return customErrors.ErrUserExists
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = uc.UserRepo.RegisterUser(ctx, email, username, string(passwordHash), name)
	if err != nil {
		return err
	}
	return nil
}
