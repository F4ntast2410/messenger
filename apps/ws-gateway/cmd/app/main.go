package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"messenger/pkg/pubsub"
	"ws-gateway/internal/config"
	"ws-gateway/internal/delivery/ws"
	"ws-gateway/internal/hub"
	gwpubsub "ws-gateway/internal/pubsub"

	"github.com/redis/go-redis/v9"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.MustLoad()

	// 1. Подключение к Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr(),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Error("failed to connect to redis", slog.String("error", err.Error()))
		panic(err)
	}
	defer rdb.Close()

	ps := pubsub.New(rdb)

	// 2. Создание менеджера подписок Redis
	manager := gwpubsub.NewDynamicPubSubManager(ps)

	// 3. Создание Hub с привязкой коллбэков подписки/отписки
	h := hub.New(
		manager.SubscribeChat,   // Вызывается при появлении первого юзера чата на этом инстансе
		manager.UnsubscribeChat, // Вызывается при уходе последнего юзера чата с этого инстанса
	)
	manager.SetHub(h)

	// 4. HTTP / WebSocket сервер
	wsHandler := ws.NewHandler(h)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/ws", wsHandler.HandleConnect)

	httpServer := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: mux,
	}

	fatalErr := make(chan error, 1)

	go func() {
		logger.Info("starting HTTP (ws) server", slog.String("port", cfg.ServerPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced shutdown", slog.String("err", err.Error()))
	}

	logger.Info("all services stopped")
}
