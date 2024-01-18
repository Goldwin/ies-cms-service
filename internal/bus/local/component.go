package local

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/internal/bus/common"
)

type localEventBusComponentImpl struct {
	subscribers map[string][]common.Consumer
}

func (l *localEventBusComponentImpl) Publish(ctx context.Context, event common.Event) error {

	if l.subscribers[event.Topic] != nil {
		for _, handler := range l.subscribers[event.Topic] {
			go handler(ctx, event)
		}
		return nil
	}
	return nil
}

func (l *localEventBusComponentImpl) Subscribe(topic string, handler common.Consumer) {
	if l.subscribers[topic] == nil {
		l.subscribers[topic] = make([]common.Consumer, 0)
	}
	l.subscribers[topic] = append(l.subscribers[topic], handler)
}

func NewEventBusComponent() *localEventBusComponentImpl {
	return &localEventBusComponentImpl{
		subscribers: make(map[string][]common.Consumer),
	}
}
