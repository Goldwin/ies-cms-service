package bus

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/internal/bus/common"
	"github.com/Goldwin/ies-pik-cms/internal/bus/local"
	"github.com/Goldwin/ies-pik-cms/internal/bus/redis"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
)

type EventBusComponent interface {
	Publish(ctx context.Context, event common.Event) error
	Subscribe(topic string, handler common.Consumer)
}

func Redis(infraComponent infra.InfraComponent) EventBusComponent {
	return redis.NewEventBusComponent(infraComponent.Redis())
}

func Local() EventBusComponent {
	return local.NewEventBusComponent()
}
