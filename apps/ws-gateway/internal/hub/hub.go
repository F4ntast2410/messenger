package hub

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const sendBufferSize = 32

type Client struct {
	conn   *websocket.Conn
	userID uuid.UUID
	send   chan []byte
}

type Hub struct {
	mu        sync.RWMutex
	clients   map[uuid.UUID]*Client
	userChats map[uuid.UUID]map[uuid.UUID]bool
	chatUsers map[uuid.UUID]map[uuid.UUID]bool

	OnFirstChatUserJoined func(chatID uuid.UUID)
	OnLastChatUserLeft    func(chatID uuid.UUID)
}

func New(onSub, onUnsub func(chatID uuid.UUID)) *Hub {
	return &Hub{
		clients:               make(map[uuid.UUID]*Client),
		userChats:             make(map[uuid.UUID]map[uuid.UUID]bool),
		chatUsers:             make(map[uuid.UUID]map[uuid.UUID]bool),
		OnFirstChatUserJoined: onSub,
		OnLastChatUserLeft:    onUnsub,
	}
}
func (c *Client) Conn() *websocket.Conn { return c.conn }
func (c *Client) Send() chan []byte     { return c.send }
func (c *Client) UserID() uuid.UUID     { return c.userID }
func (h *Hub) IsUserInChat(userID, chatID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.userChats[userID][chatID]
}

// Register создаёт нового Client и возвращает его — вызывающий код (HandleConnect)
// использует возвращённый *Client, чтобы запустить его readPump/writePump.
func (h *Hub) Register(userID uuid.UUID, conn *websocket.Conn, activeChatIDs []uuid.UUID) *Client {
	h.mu.Lock()

	if old, ok := h.clients[userID]; ok {
		old.conn.Close() // старое соединение больше не нужно — новое его заменяет
	}

	client := &Client{
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, sendBufferSize),
	}
	h.clients[userID] = client

	if _, exists := h.userChats[userID]; !exists {
		h.userChats[userID] = make(map[uuid.UUID]bool)
	}

	chatsToSubscribe := make([]uuid.UUID, 0)
	for _, chatID := range activeChatIDs {
		h.userChats[userID][chatID] = true
		if _, exists := h.chatUsers[chatID]; !exists {
			h.chatUsers[chatID] = make(map[uuid.UUID]bool)
		}
		if len(h.chatUsers[chatID]) == 0 {
			chatsToSubscribe = append(chatsToSubscribe, chatID)
		}
		h.chatUsers[chatID][userID] = true
	}

	h.mu.Unlock()

	if h.OnFirstChatUserJoined != nil {
		for _, chatID := range chatsToSubscribe {
			h.OnFirstChatUserJoined(chatID)
		}
	}

	return client
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()

	if c, ok := h.clients[client.userID]; ok && client == c {
		close(c.send)
		delete(h.clients, client.userID)
	} else {
		h.mu.Unlock()
		return
	}

	chatsToUnsubscribe := make([]uuid.UUID, 0)
	if chats, ok := h.userChats[client.userID]; ok {
		for chatID := range chats {
			if users, ok := h.chatUsers[chatID]; ok {
				delete(users, client.userID)
				if len(users) == 0 {
					delete(h.chatUsers, chatID)
					chatsToUnsubscribe = append(chatsToUnsubscribe, chatID)
				}
			}
		}
		delete(h.userChats, client.userID)
	}

	h.mu.Unlock()

	if h.OnLastChatUserLeft != nil {
		for _, chatID := range chatsToUnsubscribe {
			h.OnLastChatUserLeft(chatID)
		}
	}
}

func (h *Hub) BroadcastToChat(chatID uuid.UUID, payload any, excludeUserID uuid.UUID) {
	h.mu.RLock()
	userSet, exists := h.chatUsers[chatID]
	if !exists || len(userSet) == 0 {
		h.mu.RUnlock()
		return
	}

	targets := make([]*Client, 0, len(userSet))
	for uID := range userSet {
		if uID == excludeUserID {
			continue
		}
		if client, ok := h.clients[uID]; ok {
			targets = append(targets, client)
		}
	}
	h.mu.RUnlock()

	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("failed to marshal ws payload", slog.String("error", err.Error()))
		return
	}

	for _, client := range targets {
		select {
		case client.send <- data:
		default:
			slog.Warn("client send buffer full, dropping connection", slog.String("user_id", client.userID.String()))
			h.Unregister(client)
		}
	}
}
