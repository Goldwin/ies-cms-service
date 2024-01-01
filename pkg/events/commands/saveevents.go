package commands

import (
	"fmt"
	"time"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
)

const (
	CreateEventCommandsFailedToCreateEvent AppErrorCode = 30101
)

type SaveEventCommand struct {
	Input dto.ChurchEvent
}

func (cmd SaveEventCommand) Execute(ctx repositories.CommandContext) AppExecutionResult[dto.ChurchEvent] {
	event := entities.ChurchEvent{
		ID:        fmt.Sprintf("%d.%s", time.Now().Unix(), cmd.Input.Name),
		Name:      cmd.Input.Name,
		StartTime: cmd.Input.StartTime,
	}

	err := ctx.ChurchEventRepository().Save(event)

	if err != nil {
		return AppExecutionResult[dto.ChurchEvent]{
			Status: ExecutionStatusFailed,
			Error: AppErrorDetail{
				Code:    CreateEventCommandsFailedToCreateEvent,
				Message: fmt.Sprintf("Failed to Create Event: %s", err.Error()),
			},
		}
	}

	return AppExecutionResult[dto.ChurchEvent]{
		Status: ExecutionStatusSuccess,
		Result: dto.ChurchEvent{
			ID:        event.ID,
			Name:      event.Name,
			StartTime: event.StartTime,
		},
	}
}
