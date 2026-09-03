package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `env:"ENV" env-default:"development"`
	ServerPort  string `env:"SERVER_PORT" env-default:"8080"`
	DatabaseURL string `env:"DATABASE_URL" env-required:"true"`
	WSGateway   WSGatewayConfig
	Redis       RedisConfig
}

type WSGatewayConfig struct {
	Host string `env:"WS_GATEWAY_HOST" env-default:"localhost"`
	Port string `env:"WS_GATEWAY_GRPC_PORT" env-default:"50053"`
}

func (w WSGatewayConfig) Addr() string {
	return fmt.Sprintf("%s:%s", w.Host, w.Port)
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Port     string `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" env-default:"0"`
}

// Addr собирает host:port для redis.Options — так проще, чем городить URL
func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func MustLoad() *Config {
	var cfg Config

	_ = cleanenv.ReadConfig(".env", &cfg)

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		slog.Error("failed to load env config", slog.String("error", err.Error()))
		panic(err)
	}

	return &cfg
}
