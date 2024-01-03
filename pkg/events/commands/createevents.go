package commands

import (
	"fmt"

	"github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
	"github.com/google/uuid"
)

const (
	CreateEventCommandsErrorAppError AppErrorCode = 30101
)

type CreateEventCommands struct {
	Input dto.ChurchEvent
}

func (c CreateEventCommands) Execute(ctx CommandContext) AppExecutionResult[dto.ChurchEvent] {
	c.Input.ID = uuid.NewString()
	event := entities.ChurchEvent{
		ID:                     c.Input.ID,
		Locations:              convertToLocationEntities(c.Input.Locations),
		Name:                   c.Input.Name,
		EventFrequency:         entities.Frequency(c.Input.EventFrequency),
		LatestSessionStartTime: c.Input.LatestSessionStartTime,
		LatestShowAt:           c.Input.ShowAt,
		LatestHideAt:           c.Input.HideAt,
	}

	session := entities.ChurchEventSession{
		ID:        fmt.Sprintf("%s.1", c.Input.ID),
		Name:      c.Input.Name,
		SessionNo: 1,
		StartTime: c.Input.LatestSessionStartTime,
		ShowAt:    c.Input.ShowAt,
		HideAt:    c.Input.HideAt,
	}

	err := ctx.ChurchEventSessionRepository().Save(session)
	if err != nil {
		return AppExecutionResult[dto.ChurchEvent]{
			Status: ExecutionStatusFailed,
			Error: commands.AppErrorDetail{
				Code:    CreateEventCommandsErrorAppError,
				Message: fmt.Sprintf("Failed to Create Event: %s", err.Error()),
			},
		}
	}

	err = ctx.ChurchEventRepository().Save(event)
	if err != nil {
		return AppExecutionResult[dto.ChurchEvent]{
			Status: ExecutionStatusFailed,
			Error: commands.AppErrorDetail{
				Code:    CreateEventCommandsErrorAppError,
				Message: fmt.Sprintf("Failed to Create Event: %s", err.Error()),
			},
		}
	}

	return AppExecutionResult[dto.ChurchEvent]{
		Status: ExecutionStatusSuccess,
		Result: c.Input,
	}
}

func convertToLocationEntities(locations []dto.Location) []entities.Location {
	var result []entities.Location
	for _, location := range locations {
		result = append(result, entities.Location{
			Name:      location.Name,
			AgeFilter: entities.AgeFilter(location.AgeFilter),
		})
	}
	return result
}
