package bus

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/internal/bus/common"
	"github.com/Goldwin/ies-pik-cms/internal/bus/local"
)

type EventBusComponent interface {
	Publish(ctx context.Context, event common.Event) error
	Subscribe(topic string, handler common.Consumer)
}

func NewLocalEventBusComponent() EventBusComponent {
	return local.NewEventBusComponent()
}
