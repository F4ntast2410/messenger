package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type PubSub struct {
	client *redis.Client
}

func New(client *redis.Client) *PubSub {
	return &PubSub{client: client}
}

// Publish сериализует событие в JSON и отправляет в канал Redis
func (p *PubSub) Publish(ctx context.Context, channel string, event ChatEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return p.client.Publish(ctx, channel, data).Err()
}

// Subscribe подписывается на канал и вызывает handler для каждого полученного сообщения
func (p *PubSub) Subscribe(ctx context.Context, channel string, handler func(payload []byte)) error {
	pubsub := p.client.Subscribe(ctx, channel)
	defer pubsub.Close()

	// Подтверждение установления подписки
	if _, err := pubsub.Receive(ctx); err != nil {
		return fmt.Errorf("subscribe failed for channel %s: %w", channel, err)
	}

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			handler([]byte(msg.Payload))
		}
	}
}
