package attendance

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type AttendanceDataLayerComponent interface {
	QueryWorker() worker.QueryWorker[QueryContext]
	CommandWorker() worker.UnitOfWork[commands.CommandContext]
}

type AttendanceComponent interface {
	CreateEventSchedule(ctx context.Context, schedule EventScheduleDTO, output out.Output[EventScheduleDTO])
	AddEventScheduleActivity(ctx context.Context, activity EventScheduleActivityDTO, output out.Output[EventScheduleDTO])
	CreateEvent(ctx context.Context, scheduleID string, output out.Output[EventDTO])
	UpdateEventSchedule(ctx context.Context, schedule EventScheduleDTO, output out.Output[EventScheduleDTO])
	RemoveEventScheduleActivity(ctx context.Context, activity EventScheduleActivityDTO, output out.Output[EventScheduleDTO])
}

type attendanceComponentImpl struct {
	dataLayer AttendanceDataLayerComponent
}

// UpdateEventSchedule implements AttendanceComponent.
func (a *attendanceComponentImpl) UpdateEventSchedule(ctx context.Context, schedule EventScheduleDTO, output out.Output[EventScheduleDTO]) {
	var result CommandExecutionResult[entities.EventSchedule]
	a.dataLayer.CommandWorker().Execute(ctx, func(cc commands.CommandContext) error {
		result = commands.UpdateEventScheduleCommand{
			ID:             schedule.ID,
			Name:           schedule.Name,
			ScheduleType:   string(schedule.Type),
			TimezoneOffset: schedule.TimezoneOffset,
			Days:           schedule.Days,
			Date:           schedule.Date,
			StartDate:      schedule.StartDate,
			EndDate:        schedule.EndDate,
		}.Execute(cc)
		return nil
	})

	if result.Status == ExecutionStatusFailed {
		output.OnError(out.ConvertCommandErrorDetail(result.Error))
		return
	}

	output.OnSuccess(fromEntities(&result.Result))
}

// RemoveEventScheduleActivity implements AttendanceComponent.
func (a *attendanceComponentImpl) RemoveEventScheduleActivity(ctx context.Context, activity EventScheduleActivityDTO, output out.Output[EventScheduleDTO]) {
	var result CommandExecutionResult[entities.EventSchedule]
	a.dataLayer.CommandWorker().Execute(ctx, func(cc commands.CommandContext) error {
		result = commands.RemoveScheduleActivityCommand{
			ScheduleID: activity.ScheduleID,
			ActivityID: activity.ID,
		}.Execute(cc)
		return nil
	})

	if result.Status == ExecutionStatusFailed {
		output.OnError(out.ConvertCommandErrorDetail(result.Error))
		return
	}

	output.OnSuccess(fromEntities(&result.Result))
}

// AddEventScheduleActivity implements AttendanceComponent.
func (a *attendanceComponentImpl) AddEventScheduleActivity(ctx context.Context, activity EventScheduleActivityDTO, output out.Output[EventScheduleDTO]) {
	var result CommandExecutionResult[entities.EventSchedule]
	a.dataLayer.CommandWorker().Execute(ctx, func(cc commands.CommandContext) error {
		result = commands.AddEventScheduleActivityCommand{
			ScheduleID: activity.ScheduleID,
			Name:       activity.Name,
			Hour:       activity.Hour,
			Minute:     activity.Minute,
		}.Execute(cc)
		return nil
	})

	if result.Status == ExecutionStatusFailed {
		output.OnError(out.ConvertCommandErrorDetail(result.Error))
		return
	}

	output.OnSuccess(fromEntities(&result.Result))
}

// CreateEvent implements AttendanceComponent.
func (a *attendanceComponentImpl) CreateEvent(ctx context.Context, scheduleID string, output out.Output[EventDTO]) {
	panic("unimplemented")
}

// CreateEventSchedule implements AttendanceComponent.
func (a *attendanceComponentImpl) CreateEventSchedule(ctx context.Context, schedule EventScheduleDTO, output out.Output[EventScheduleDTO]) {

	var result CommandExecutionResult[entities.EventSchedule]
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

	output.OnSuccess(fromEntities(&result.Result))
}

func NewAttendanceComponent(datalayer AttendanceDataLayerComponent) AttendanceComponent {
	return &attendanceComponentImpl{
		dataLayer: datalayer,
	}
}
