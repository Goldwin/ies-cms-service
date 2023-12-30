package redis

import (
	"context"
	"log"

	"github.com/Goldwin/ies-pik-cms/internal/bus/common"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

type redisEventBusComponentImpl struct {
	redisClient redis.UniversalClient
}

func (r *redisEventBusComponentImpl) Publish(ctx context.Context, event common.Event) error {
	bytes, err := msgpack.Marshal(event)
	if err != nil {
		return err
	}
	return r.redisClient.Publish(ctx, event.Topic, bytes).Err()
}

func (r *redisEventBusComponentImpl) Subscribe(topic string, handler common.Consumer) {
	ctx := context.Background()
	pubsub := r.redisClient.Subscribe(ctx, topic)
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var event common.Event
			err := msgpack.Unmarshal([]byte(msg.Payload), &event)
			if err != nil {
				log.Printf("Consuming topic %v failed. Error: %s", topic, err.Error())
				continue
			}
			handler(context.Background(), event)
		}
	}()
}

func NewEventBusComponent(redis redis.UniversalClient) *redisEventBusComponentImpl {
	return &redisEventBusComponentImpl{
		redisClient: redis,
	}
}

func getKey(topic string) string {
	return topic
}
