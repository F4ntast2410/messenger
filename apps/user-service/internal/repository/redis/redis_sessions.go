package cache

import (
	"context"
	"fmt"
	"time"
	customErrors "user-service/internal/errors"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (c RedisVerificationRepository) DeleteSessionToken(ctx context.Context, sessionID string) error {
	err := c.Client.Del(ctx, fmt.Sprintf("session:session_token:%s", sessionID)).Err()
	if err != nil {
		return err
	}
	return nil
}
func (c RedisVerificationRepository) GetUserID(ctx context.Context, sessionID string) (uuid.UUID, error) {
	result, err := c.Client.Get(ctx, fmt.Sprintf("session:session_token:%s", sessionID)).Result()
	if err != nil {
		if err == redis.Nil {
			return uuid.Nil, customErrors.ErrSessionInvalid
		}
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(result)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, err
}
func (c RedisVerificationRepository) SaveSessionToken(ctx context.Context, sessionID string, userID uuid.UUID, time_duration time.Duration) error {
	err := c.Client.Set(ctx, fmt.Sprintf("session:session_token:%s", sessionID), userID.String(), time_duration).Err()
	if err != nil {
		return err
	}
	return nil
}
