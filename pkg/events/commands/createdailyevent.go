package commands

import (
	"time"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
)

type CreateDailyEventCommand struct {
	Day time.Weekday
}

func (cmd CreateDailyEventCommand) Execute(ctx repositories.CommandContext) AppExecutionResult[[]dto.ChurchEvent] {
	schedules := make([]entities.ChurchEventSchedule, 0)
	events := make([]dto.ChurchEvent, 0)

	for i := -12; i <= 14; i++ {
		schedule, _ := ctx.ChurchEventScheduleRepository().GetByTimezoneAndWeekDay(i, cmd.Day)
		if schedule != nil {
			schedules = append(schedules, *schedule)
		}
	}

	for _, schedule := range schedules {
		startTime := schedule.GetNextEventDate()
		result := SaveEventCommand{
			Input: dto.ChurchEvent{
				Name:      schedule.Name,
				StartTime: startTime,
			},
		}.Execute(ctx)
		if result.Status == ExecutionStatusSuccess {
			events = append(events, result.Result)
		}
	}
	return AppExecutionResult[[]dto.ChurchEvent]{
		Status: ExecutionStatusSuccess,
		Result: events,
	}
}
