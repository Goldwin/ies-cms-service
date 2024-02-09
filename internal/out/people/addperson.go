package people

import (
	"context"
	"log"

	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/internal/bus/common"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/vmihailenco/msgpack/v5"
)

type addPersonOutput struct {
	eventBus bus.EventBusComponent
}

// OnError implements out.Output.
func (*addPersonOutput) OnError(err out.AppErrorDetail) {
	log.Default().Printf("Error Found when adding a person: %s", err.Error())
}

// OnSuccess implements out.Output.
func (a *addPersonOutput) OnSuccess(result dto.Person) {
	body, err := msgpack.Marshal(result)
	if err != nil {
		log.Default().Printf("cancelling message publishing. Failed to marshall the body: %s", err.Error())
	}
	a.eventBus.Publish(context.Background(), common.Event{Topic: "people.added", Body: body})
}

func AddPersonOutput(bus bus.EventBusComponent) out.Output[dto.Person] {
	return &addPersonOutput{
		eventBus: bus,
	}
}
