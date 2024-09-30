package auth

import (
	"context"
	"log"

	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/internal/bus/common"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/vmihailenco/msgpack/v5"
)

type registerOutputHandler struct {
	bus bus.EventBusComponent
}

// OnError implements out.Output.
func (*registerOutputHandler) OnError(err out.AppErrorDetail) {
	log.Default().Printf("Error Found when registering: %s", err.Error())
}

// OnSuccess implements out.Output.
func (o *registerOutputHandler) OnSuccess(result dto.AuthData) {
	body, err := msgpack.Marshal(result)
	if err != nil {
		log.Default().Printf("cancelling message publishing. Failed to marshall the body: %s", err.Error())
	}

	o.bus.Publish(context.Background(), common.Event{Topic: "auth.registered", Body: body})
}

func newRegisterOutputHandler(bus bus.EventBusComponent) out.Output[dto.AuthData] {
	return &registerOutputHandler{
		bus: bus,
	}
}
