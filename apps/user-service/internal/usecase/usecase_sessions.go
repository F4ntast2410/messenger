package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
	"user-service/internal/entity"
	customErrors "user-service/internal/errors"

	"github.com/google/uuid"
)

func (uc *UsecaseImpl) ValidateSession(ctx context.Context, sessionID string) (*entity.UserSession, error) {
	userID, err := uc.AuthRepo.GetUserID(ctx, hashToken(sessionID))
	if err != nil {
		if errors.Is(err, customErrors.ErrSessionInvalid) {
			return &entity.UserSession{IsValid: false}, nil
		}
		return nil, err
	}
	return &entity.UserSession{
		UserID:  userID,
		IsValid: true,
	}, nil
}
func generateSessionToken() (rawToken string, tokenHash string, err error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", err
	}
	rawToken = hex.EncodeToString(raw)
	tokenHash = hashToken(rawToken)
	return rawToken, tokenHash, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
func (uc *UsecaseImpl) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	rawToken, tokenHash, err := generateSessionToken()
	if err != nil {
		return "", err
	}
	if err := uc.AuthRepo.SaveSessionToken(ctx, tokenHash, userID, 24*time.Hour); err != nil {
		return "", err
	}
	return rawToken, err
}

func (uc *UsecaseImpl) DeleteSession(ctx context.Context, sessionToken string) error {
	return uc.AuthRepo.DeleteSessionToken(ctx, sessionToken)
}
