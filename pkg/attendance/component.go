package attendance

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/commands"
	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	commonCmd "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/samber/lo"
)

type AttendanceDataLayerComponent interface {
	QueryWorker() worker.QueryWorker[QueryContext]
	CommandWorker() worker.UnitOfWork[commands.CommandContext]
}

type AttendanceComponent interface {
	CreateEventSchedule(ctx context.Context, schedule EventScheduleDTO, output out.Output[EventScheduleDTO])
	AddEventScheduleActivity(ctx context.Context, activity EventScheduleActivity)
	CreateEvent(ctx context.Context, event Event)
	CreateEventActivity(ctx context.Context, activity EventActivity)
}

type attendanceComponentImpl struct {
	dataLayer AttendanceDataLayerComponent
}

// AddEventScheduleActivity implements AttendanceComponent.
func (a *attendanceComponentImpl) AddEventScheduleActivity(ctx context.Context, activity EventScheduleActivity) {
	panic("unimplemented")
}

// CreateEvent implements AttendanceComponent.
func (a *attendanceComponentImpl) CreateEvent(ctx context.Context, event Event) {
	panic("unimplemented")
}

// CreateEventActivity implements AttendanceComponent.
func (a *attendanceComponentImpl) CreateEventActivity(ctx context.Context, activity EventActivity) {
	panic("unimplemented")
}

// CreateEventSchedule implements AttendanceComponent.
func (a *attendanceComponentImpl) CreateEventSchedule(ctx context.Context, schedule EventScheduleDTO, output out.Output[EventScheduleDTO]) {

	var result commonCmd.CommandExecutionResult[EventSchedule]
	err := a.dataLayer.CommandWorker().Execute(ctx, func(cc commands.CommandContext) error {
		result = commands.CreateEventScheduleCommand{
			Name:           schedule.Name,
			ScheduleType:   schedule.Type,
			TimezoneOffset: schedule.TimezoneOffset,
			Days:           schedule.Days,
			Date:           schedule.Date,
			StartDate:      schedule.StartDate,
			EndDate:        schedule.EndDate,
		}.Execute(cc)
		return nil
	})

	if err != nil {
		output.OnError(out.ConvertCommandErrorDetail(result.Error))
		return
	}

	output.OnSuccess(EventScheduleDTO{
		ID:             result.Result.ID,
		Name:           result.Result.Name,
		TimezoneOffset: result.Result.TimezoneOffset,
		Type:           string(result.Result.Type),
		Activities: lo.Map(result.Result.Activities,
			func(ea EventScheduleActivity, _ int) EventScheduleActivityDTO {
				return EventScheduleActivityDTO{ID: ea.ID, Name: ea.Name, Hour: ea.Hour, Minute: ea.Minute}
			}),
		Date:      result.Result.Date,
		Days:      result.Result.Days,
		StartDate: result.Result.StartDate,
		EndDate:   result.Result.EndDate,
	})
}

func NewAttendanceComponent(datalayer AttendanceDataLayerComponent) AttendanceComponent {
	return &attendanceComponentImpl{
		dataLayer: datalayer,
	}
}
