package ws

import (
	"log/slog"
	"net/http"
	"strings"

	"ws-gateway/internal/hub"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// В проде тут нужна нормальная проверка Origin,
	// пока разрешаем всё, т.к. за nginx.
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	hub *hub.Hub
}

func NewHandler(h *hub.Hub) *Handler {
	return &Handler{hub: h}
}

func (h *Handler) HandleConnect(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем ID пользователя из заголовка, подставленного Nginx
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 2. Извлекаем список ID чатов из query-параметров (?chat_ids=uuid1,uuid2)
	var activeChatIDs []uuid.UUID
	if chatIDsParam := r.URL.Query().Get("chat_ids"); chatIDsParam != "" {
		rawIDs := strings.Split(chatIDsParam, ",")
		for _, rawID := range rawIDs {
			if id, err := uuid.Parse(strings.TrimSpace(rawID)); err == nil {
				activeChatIDs = append(activeChatIDs, id)
			}
		}
	}

	// 3. Апгрейдим HTTP-соединение до WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("ws upgrade failed", slog.String("error", err.Error()))
		return
	}

	// 4. Регистрируем пользователя и его чаты в Hub
	// Внутри Hub проверит, есть ли среди activeChatIDs новые чаты для этого инстанса,
	// и автоматически вызовет менеджер подписок Redis.
	h.hub.Register(userID, conn, activeChatIDs)
	slog.Info("user connected",
		slog.String("user_id", userID.String()),
		slog.Int("active_chats_count", len(activeChatIDs)),
	)

	// 5. Запускаем фоновую горутину для прослушивания закрытия сокета
	go h.readLoop(userID, conn)
}

func (h *Handler) readLoop(userID uuid.UUID, conn *websocket.Conn) {
	defer func() {
		h.hub.Unregister(userID)
		conn.Close()
		slog.Info("user disconnected", slog.String("user_id", userID.String()))
	}()

	for {
		// Читающий цикл нужен, чтобы вовремя узнать о закрытии соединения
		// (клиент сам не пишет нам ничего, кроме pong/close).
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
