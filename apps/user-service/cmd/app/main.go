package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/internal/config"
	httpHandler "user-service/internal/delivery/http"
	"user-service/internal/delivery/http/middleware"
	repository "user-service/internal/repository/postgres"
	cache "user-service/internal/repository/redis"
	"user-service/internal/usecase"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.MustLoad()

	var db *sqlx.DB
	var err error

	maxRetries := 10
	for i := 1; i <= maxRetries; i++ {
		db, err = sqlx.Connect("pgx", cfg.DatabaseURL)
		if err == nil {
			logger.Info("successfully connected to database")
			break
		}

		logger.Warn("database connection failed, retrying...",
			slog.Int("attempt", i),
			slog.Int("max_attempts", maxRetries),
			slog.String("error", err.Error()),
		)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Error("failed to connect to database after max retries", slog.String("error", err.Error()))
		panic(err)
	}
	defer db.Close()

	opt := &redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("redis connection failed", slog.String("error", err.Error()))
		panic(err)
	}
	cache := cache.NewRedisRepository(client)
	storage := repository.NewPostgresStorage(db)
	userUsecase := usecase.NewUsecase(cache, storage)
	userHandler := httpHandler.NewUserHandler(userUsecase, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/users/me", userHandler.GetUserHandler)
	mux.HandleFunc("POST /api/v1/users/login", userHandler.LoginUserHandler)
	mux.HandleFunc("POST /api/v1/users/logout", userHandler.LogoutUserHandler)
	mux.HandleFunc("POST /api/v1/users/register", userHandler.RegisterUserHandler)
	mux.HandleFunc("GET /api/v1/auth/verify", userHandler.CheckSessionIDHandler)
	mux.HandleFunc("GET /api/v1/users/search", userHandler.SearchUsersHandler)
	wrappedMux := middleware.Logger(mux)
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: wrappedMux,
	}
	// Канал для аварийных ошибок — буфер 2, чтобы обе горутины могли записать не блокируясь
	fatalErr := make(chan error, 1)
	// Запускаем HTTP-сервер
	go func() {
		logger.Info("starting HTTP server", slog.String("port", cfg.ServerPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", slog.String("err", err.Error()))
			fatalErr <- err // сообщаем об аварии
		}
	}()
	// Ждём сигнала остановки (Ctrl+C или docker stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ждём: либо сигнал от ОС, либо авария одного из сервисов
	select {
	case sig := <-quit:
		logger.Info("shutdown signal received", slog.String("signal", sig.String()))
	case err := <-fatalErr:
		logger.Error("fatal error, initiating shutdown", slog.String("err", err.Error()))
	}

	// Дальше — тот же код остановки
	logger.Info("stopping services...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced shutdown", slog.String("err", err.Error()))
	}

	logger.Info("all services stopped")
}
