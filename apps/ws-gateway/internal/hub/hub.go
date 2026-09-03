package hub

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	mu        sync.RWMutex
	conns     map[uuid.UUID]*websocket.Conn
	userChats map[uuid.UUID]map[uuid.UUID]bool // userID -> set of chatIDs
	chatUsers map[uuid.UUID]map[uuid.UUID]bool // chatID -> set of userIDs

	OnFirstChatUserJoined func(chatID uuid.UUID)
	OnLastChatUserLeft    func(chatID uuid.UUID)
}

func New(onSub, onUnsub func(chatID uuid.UUID)) *Hub {
	return &Hub{
		conns:                 make(map[uuid.UUID]*websocket.Conn),
		userChats:             make(map[uuid.UUID]map[uuid.UUID]bool),
		chatUsers:             make(map[uuid.UUID]map[uuid.UUID]bool),
		OnFirstChatUserJoined: onSub,
		OnLastChatUserLeft:    onUnsub,
	}
}

func (h *Hub) Register(userID uuid.UUID, conn *websocket.Conn, activeChatIDs []uuid.UUID) {
	h.mu.Lock()

	if old, ok := h.conns[userID]; ok {
		old.Close()
	}
	h.conns[userID] = conn

	if _, exists := h.userChats[userID]; !exists {
		h.userChats[userID] = make(map[uuid.UUID]bool)
	}

	chatsToSubscribe := make([]uuid.UUID, 0)

	for _, chatID := range activeChatIDs {
		h.userChats[userID][chatID] = true

		if _, exists := h.chatUsers[chatID]; !exists {
			h.chatUsers[chatID] = make(map[uuid.UUID]bool)
		}

		// Если это первый пользователь данного чата на ЭТОМ инстансе ws-gateway
		if len(h.chatUsers[chatID]) == 0 {
			chatsToSubscribe = append(chatsToSubscribe, chatID)
		}
		h.chatUsers[chatID][userID] = true
	}

	h.mu.Unlock()

	// Вызываем внешние коллбеки для оформления подписок в Redis вне блокировки мьютекса
	if h.OnFirstChatUserJoined != nil {
		for _, chatID := range chatsToSubscribe {
			h.OnFirstChatUserJoined(chatID)
		}
	}
}

func (h *Hub) Unregister(userID uuid.UUID) {
	h.mu.Lock()

	delete(h.conns, userID)

	chatsToUnsubscribe := make([]uuid.UUID, 0)

	if chats, ok := h.userChats[userID]; ok {
		for chatID := range chats {
			if users, ok := h.chatUsers[chatID]; ok {
				delete(users, userID)
				// Если в чате не осталось ни одного подключенного пользователя на этом инстансе
				if len(users) == 0 {
					delete(h.chatUsers, chatID)
					chatsToUnsubscribe = append(chatsToUnsubscribe, chatID)
				}
			}
		}
		delete(h.userChats, userID)
	}

	h.mu.Unlock()

	if h.OnLastChatUserLeft != nil {
		for _, chatID := range chatsToUnsubscribe {
			h.OnLastChatUserLeft(chatID)
		}
	}
}

// BroadcastToChat рассылает полученное из Redis событие всем онлайн-участникам чата на данном инстансе
func (h *Hub) BroadcastToChat(chatID uuid.UUID, payload any, excludeUserID uuid.UUID) {
	h.mu.RLock()
	userSet, exists := h.chatUsers[chatID]
	if !exists || len(userSet) == 0 {
		h.mu.RUnlock()
		return
	}

	// Собираем список соединений для отправки
	targets := make([]*websocket.Conn, 0, len(userSet))
	targetIDs := make([]uuid.UUID, 0, len(userSet))

	for uID := range userSet {
		if uID == excludeUserID {
			continue // Не отправляем инициатору события
		}
		if conn, ok := h.conns[uID]; ok {
			targets = append(targets, conn)
			targetIDs = append(targetIDs, uID)
		}
	}
	h.mu.RUnlock()

	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("failed to marshal ws payload", slog.String("error", err.Error()))
		return
	}

	for i, conn := range targets {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			slog.Warn("failed to write to socket, dropping connection",
				slog.String("user_id", targetIDs[i].String()),
				slog.String("error", err.Error()),
			)
			h.Unregister(targetIDs[i])
		}
	}
}
