package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"messenger/pkg/pubsub"
	"net/http"
	"strings"
	"time"

	"ws-gateway/internal/hub"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	writeWait  = 10 * time.Second
)

type Handler struct {
	hub *hub.Hub
	ps  *pubsub.PubSub
}

func NewHandler(h *hub.Hub, ps *pubsub.PubSub) *Handler {
	return &Handler{hub: h, ps: ps}
}

func (h *Handler) HandleConnect(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var activeChatIDs []uuid.UUID
	if chatIDsParam := r.URL.Query().Get("chat_ids"); chatIDsParam != "" {
		rawIDs := strings.Split(chatIDsParam, ",")
		for _, rawID := range rawIDs {
			if id, err := uuid.Parse(strings.TrimSpace(rawID)); err == nil {
				activeChatIDs = append(activeChatIDs, id)
			}
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("ws upgrade failed", slog.String("error", err.Error()))
		return
	}

	client := h.hub.Register(userID, conn, activeChatIDs)
	slog.Info("user connected",
		slog.String("user_id", userID.String()),
		slog.Int("active_chats_count", len(activeChatIDs)),
	)

	go h.writePump(client)
	go h.readPump(client)
}

func (h *Handler) writePump(client *hub.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn().Close()
	}()

	for {
		select {
		case data, ok := <-client.Send():
			client.Conn().SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {

				client.Conn().WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := client.Conn().WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn().SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn().WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *Handler) readPump(client *hub.Client) {
	defer func() {
		h.hub.Unregister(client)
		slog.Info("user disconnected", slog.String("user_id", client.UserID().String()))
	}()

	conn := client.Conn()
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			slog.Warn("failed to read message with ws", slog.String("error", err.Error()))
			return
		}

		var eventType pubsub.TypeEvent
		if err := json.Unmarshal(data, &eventType); err != nil {
			slog.Warn("failed to unmarshal data to event type", slog.String("error", err.Error()))
			continue
		}

		switch eventType.Type {
		case pubsub.EventUserTyping:
			h.handleTyping(client, data)
		}
	}
}

func (h *Handler) handleTyping(client *hub.Client, data []byte) {
	var event pubsub.ChatEvent
	event.CreatedAt = time.Now()
	if err := json.Unmarshal(data, &event); err != nil {
		slog.Warn("failed to unmarshal data to userTypingEvent", slog.String("error", err.Error()))
		return
	}

	chatID, err := uuid.Parse(event.ChatID)
	if err != nil {
		slog.Warn("failed to parse chat_id to uuid from string", slog.String("error", err.Error()))
		return
	}

	if !h.hub.IsUserInChat(client.UserID(), chatID) {
		slog.Warn("user hasn't chat with id", slog.String("chat_id", event.ChatID))
		return
	}

	event.SenderID = client.UserID().String()
	channel := fmt.Sprintf("chat:%s", event.ChatID)
	if err := h.ps.Publish(context.Background(), channel, event); err != nil {
		slog.Error("failed to publish event", slog.String("error", err.Error()))
	}
}
