package cache

import "github.com/redis/go-redis/v9"

type RedisVerificationRepository struct {
	Client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisVerificationRepository {
	return &RedisVerificationRepository{Client: client}
}
