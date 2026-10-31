package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"chat-service/internal/config"
	handler "chat-service/internal/delivery/http"
	"chat-service/internal/delivery/http/middleware"
	repository "chat-service/internal/repository/postgres"
	"chat-service/internal/usecase"
	"messenger/pkg/pubsub"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.MustLoad()

	// 1. Подключение к PostgreSQL
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

	// 2. Инициализация Redis Pub/Sub
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr(),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Error("failed to connect to redis", slog.String("error", err.Error()))
		panic(err)
	}
	defer rdb.Close()

	ps := pubsub.New(rdb)

	// 3. Сборка UseCase и HTTP Handlers
	storage := repository.NewPostgresStorage(db)
	chatService := usecase.NewChatService(storage, ps)
	h := handler.NewHandler(chatService, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/chats/{id}/messages", h.HandleSendMessage)
	mux.HandleFunc("GET /api/v1/chats/{id}/messages", h.HandleGetMessagesPaginated)
	mux.HandleFunc("POST /api/v1/chats", h.HandleCreateChat)
	mux.HandleFunc("POST /api/v1/chats/{id}", h.HandleJoinToChat)
	mux.HandleFunc("PUT /api/v1/chats/{id}/read", h.HandleUpdateLastReadMessageID)
	mux.HandleFunc("GET /api/v1/chats", h.HandleGetUserChats)

	wrappedMux := middleware.Logger(mux)
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: wrappedMux,
	}

	fatalErr := make(chan error, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("starting HTTP server", slog.String("port", cfg.ServerPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", slog.String("err", err.Error()))
			fatalErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logger.Info("shutdown signal received", slog.String("signal", sig.String()))
	case err := <-fatalErr:
		logger.Error("fatal error, initiating shutdown", slog.String("err", err.Error()))
	}

	logger.Info("stopping services...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced shutdown", slog.String("err", err.Error()))
	}

	wg.Wait()
	logger.Info("all services stopped")
}
