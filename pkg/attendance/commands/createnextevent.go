package commands

import (
	"fmt"
	"time"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

const (
	CreateEventCommandsErrorScheduleDoesntExists CommandErrorCode = 30301
	CreateEventCommandsErrorInvalidScheduleType  CommandErrorCode = 30302
	CreateEventCommandsErrorNoActivities         CommandErrorCode = 30303
)

type CreateNextEventCommand struct {
	ScheduleID string
}

func (c CreateNextEventCommand) Execute(ctx CommandContext) CommandExecutionResult[[]*entities.Event] {
	schedule, err := ctx.EventScheduleRepository().Get(c.ScheduleID)

	if err != nil {
		return CommandExecutionResult[[]*entities.Event]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	if schedule == nil {
		return CommandExecutionResult[[]*entities.Event]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CreateEventCommandsErrorScheduleDoesntExists,
				Message: "Schedule not found",
			},
		}
	}

	if schedule.Activities == nil || len(schedule.Activities) == 0 {
		return CommandExecutionResult[[]*entities.Event]{
			Status: ExecutionStatusFailed,
			Error: CommandErrorDetail{
				Code:    CreateEventCommandsErrorNoActivities,
				Message: "No activities found under this schedule. Please configure activity for this event first",
			},
		}
	}

	switch schedule.Type {
	case entities.EventScheduleTypeDaily:
		return createNextDailyEvent(schedule, ctx)
	case entities.EventScheduleTypeWeekly:
		return createNextWeeklyEvent(schedule, ctx)
	case entities.EventScheduleTypeOneTime:
		return createNextOneTimeEvent(schedule, ctx)
	}

	return CommandExecutionResult[[]*entities.Event]{
		Status: ExecutionStatusFailed,
		Error: CommandErrorDetail{
			Code:    CreateEventCommandsErrorInvalidScheduleType,
			Message: "Invalid schedule type. Please Configure the schedule first",
		},
	}
}

func generateEventId(scheduleId string, date time.Time) string {
	return fmt.Sprintf("%s.%04d%02d%02d", scheduleId, date.Year(), date.Month(), date.Day())
}

func createNextWeeklyEvent(weeklySchedule *entities.EventSchedule, ctx CommandContext) CommandExecutionResult[[]*entities.Event] {
	resultSet := make([]*entities.Event, 0)
	for d, i := time.Now(), 0; i <= 7; d, i = d.AddDate(0, 0, 1), i+1 {
		if lo.Contains(weeklySchedule.Days, d.Weekday()) {
			event := entities.Event{
				ID:         generateEventId(weeklySchedule.ID, d),
				Name:       weeklySchedule.Name,
				ScheduleID: weeklySchedule.ID,
				EventActivities: lo.Map(weeklySchedule.Activities, func(ea entities.EventScheduleActivity, _ int) *entities.EventActivity {
					return &entities.EventActivity{ID: uuid.NewString(), Name: ea.Name, Time: time.Date(d.Year(), d.Month(), d.Day(), ea.Hour-weeklySchedule.TimezoneOffset, ea.Minute, 0, 0, time.UTC)}
				}),
				StartDate: time.Date(d.Year(), d.Month(), d.Day(), weeklySchedule.StartTime.Hour-weeklySchedule.TimezoneOffset, weeklySchedule.StartTime.Minute, 0, 0, time.UTC),
				EndDate:   time.Date(d.Year(), d.Month(), d.Day(), weeklySchedule.EndTime.Hour-weeklySchedule.TimezoneOffset, weeklySchedule.EndTime.Minute, 0, 0, time.UTC),
			}
			result, err := ctx.EventRepository().Save(&event)

			if err != nil {
				return CommandExecutionResult[[]*entities.Event]{
					Status: ExecutionStatusFailed,
					Error:  CommandErrorDetailWorkerFailure(err),
				}
			}
			resultSet = append(resultSet, result)
		}
	}
	return CommandExecutionResult[[]*entities.Event]{
		Status: ExecutionStatusSuccess,
		Result: resultSet,
	}
}

func createNextDailyEvent(dailySchedule *entities.EventSchedule, ctx CommandContext) CommandExecutionResult[[]*entities.Event] {
	resultSet := make([]*entities.Event, 0)
	for d := dailySchedule.StartDate; d.Before(dailySchedule.EndDate) || d.Equal(dailySchedule.EndDate); d = d.AddDate(0, 0, 1) {
		event := entities.Event{
			ID:         generateEventId(dailySchedule.ID, d),
			ScheduleID: dailySchedule.ID,
			Name:       dailySchedule.Name,
			StartDate:  time.Date(d.Year(), d.Month(), d.Day(), dailySchedule.StartTime.Hour-dailySchedule.TimezoneOffset, dailySchedule.StartTime.Minute, 0, 0, time.UTC),
			EndDate:    time.Date(d.Year(), d.Month(), d.Day(), dailySchedule.EndTime.Hour-dailySchedule.TimezoneOffset, dailySchedule.EndTime.Minute, 0, 0, time.UTC),
			EventActivities: lo.Map(dailySchedule.Activities,
				func(ea entities.EventScheduleActivity, _ int) *entities.EventActivity {
					return &entities.EventActivity{
						ID:   uuid.NewString(),
						Name: ea.Name,
						Time: time.Date(d.Year(), d.Month(), d.Day(), ea.Hour-dailySchedule.TimezoneOffset, ea.Minute, 0, 0, time.UTC),
					}
				}),
		}

		result, err := ctx.EventRepository().Save(&event)

		if err != nil {
			return CommandExecutionResult[[]*entities.Event]{
				Status: ExecutionStatusFailed,
				Error:  CommandErrorDetailWorkerFailure(err),
			}
		}
		resultSet = append(resultSet, result)
	}
	return CommandExecutionResult[[]*entities.Event]{
		Status: ExecutionStatusSuccess,
		Result: resultSet,
	}
}

func createNextOneTimeEvent(oneTimeSchedule *entities.EventSchedule, ctx CommandContext) CommandExecutionResult[[]*entities.Event] {
	event := entities.Event{
		ID:         generateEventId(oneTimeSchedule.ID, oneTimeSchedule.Date),
		ScheduleID: oneTimeSchedule.ID,
		Name:       oneTimeSchedule.Name,
		StartDate:  time.Date(oneTimeSchedule.Date.Year(), oneTimeSchedule.Date.Month(), oneTimeSchedule.Date.Day(), oneTimeSchedule.StartTime.Hour-oneTimeSchedule.TimezoneOffset, oneTimeSchedule.StartTime.Minute, 0, 0, time.UTC),
		EndDate:    time.Date(oneTimeSchedule.Date.Year(), oneTimeSchedule.Date.Month(), oneTimeSchedule.Date.Day(), oneTimeSchedule.EndTime.Hour-oneTimeSchedule.TimezoneOffset, oneTimeSchedule.EndTime.Minute, 0, 0, time.UTC),
		EventActivities: lo.Map(oneTimeSchedule.Activities,
			func(ea entities.EventScheduleActivity, _ int) *entities.EventActivity {
				return &entities.EventActivity{
					ID:   uuid.NewString(),
					Name: ea.Name,
					Time: time.Date(oneTimeSchedule.Date.Day(), oneTimeSchedule.Date.Month(), oneTimeSchedule.Date.Day(), ea.Hour-oneTimeSchedule.TimezoneOffset, ea.Minute, 0, 0, time.UTC),
				}
			}),
	}

	result, err := ctx.EventRepository().Save(&event)

	if err != nil {
		return CommandExecutionResult[[]*entities.Event]{
			Status: ExecutionStatusFailed,
			Error:  CommandErrorDetailWorkerFailure(err),
		}
	}

	return CommandExecutionResult[[]*entities.Event]{
		Status: ExecutionStatusSuccess,
		Result: []*entities.Event{result},
	}

}
