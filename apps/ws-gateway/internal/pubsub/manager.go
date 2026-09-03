package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"messenger/pkg/pubsub"
	"ws-gateway/internal/hub"

	"github.com/google/uuid"
)

type DynamicPubSubManager struct {
	ps       *pubsub.PubSub
	hub      *hub.Hub
	mu       sync.Mutex
	cancelFn map[uuid.UUID]context.CancelFunc
}

func NewDynamicPubSubManager(ps *pubsub.PubSub) *DynamicPubSubManager {
	return &DynamicPubSubManager{
		ps:       ps,
		cancelFn: make(map[uuid.UUID]context.CancelFunc),
	}
}

func (m *DynamicPubSubManager) SetHub(h *hub.Hub) {
	m.hub = h
}

func (m *DynamicPubSubManager) SubscribeChat(chatID uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.cancelFn[chatID]; exists {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.cancelFn[chatID] = cancel

	channel := fmt.Sprintf("chat:%s", chatID.String())

	go func() {
		err := m.ps.Subscribe(ctx, channel, func(payload []byte) {
			var event pubsub.ChatEvent
			if err := json.Unmarshal(payload, &event); err != nil {
				slog.Error("failed to unmarshal chat event", slog.String("error", err.Error()))
				return
			}

			senderID, _ := uuid.Parse(event.SenderID)
			if m.hub != nil {
				m.hub.BroadcastToChat(chatID, event, senderID)
			}
		})

		if err != nil && ctx.Err() == nil {
			slog.Error("pubsub subscription error", slog.String("channel", channel), slog.String("error", err.Error()))
		}
	}()
}

func (m *DynamicPubSubManager) UnsubscribeChat(chatID uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cancel, exists := m.cancelFn[chatID]; exists {
		cancel()
		delete(m.cancelFn, chatID)
	}
}
