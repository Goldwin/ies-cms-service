package attendance

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/attendance/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/common/utils"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type AttendanceDataLayerComponent interface {
	QueryWorker() worker.QueryWorker[queries.QueryContext]
	CommandWorker() worker.UnitOfWork[commands.CommandContext]
}

type AttendanceCommandComponent interface {
	CreateEventSchedule(ctx context.Context, schedule dto.EventScheduleDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable
	AddEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable
	CreateEvent(ctx context.Context, scheduleID string, output out.Output[dto.EventDTO]) out.Waitable
	UpdateEventSchedule(ctx context.Context, schedule dto.EventScheduleDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable
	RemoveEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable
	UpdateEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable
	CheckIn(ctx context.Context, command commands.CheckInCommand, output out.Output[entities.Attendance]) out.Waitable
}

type AttendanceQueryComponent interface {
	GetEventSchedule(ctx context.Context, query queries.GetEventScheduleQuery, output out.Output[queries.GetEventScheduleResult])
	ListEventSchedules(ctx context.Context, query queries.ListEventScheduleQuery, output out.Output[queries.ListEventScheduleResult])
	ListEventsBySchedule(ctx context.Context, query queries.ListEventByScheduleQuery, output out.Output[queries.ListEventByScheduleResult])
	GetEvent(ctx context.Context, query queries.GetEventQuery, output out.Output[queries.GetEventResult])
	ListEventAttendance(ctx context.Context, query queries.ListEventAttendanceQuery, output out.Output[queries.ListEventAttendanceResult])
}

type AttendanceComponent interface {
	AttendanceCommandComponent
	AttendanceQueryComponent
}

type attendanceComponentImpl struct {
	dataLayer AttendanceDataLayerComponent
}

// GetEvent implements AttendanceComponent.
func (a *attendanceComponentImpl) GetEvent(ctx context.Context, query queries.GetEventQuery, output out.Output[queries.GetEventResult]) {

	result, err := a.dataLayer.QueryWorker().Query(ctx).GetEvent().Execute(query)
	if err.NoError() {
		output.OnSuccess(result)
		return
	}
	output.OnError(out.ConvertQueryErrorDetail(err))

}

// GetEventSchedule implements AttendanceComponent.
func (a *attendanceComponentImpl) GetEventSchedule(ctx context.Context, query queries.GetEventScheduleQuery, output out.Output[queries.GetEventScheduleResult]) {

	result, err := a.dataLayer.QueryWorker().Query(ctx).GetEventSchedule().Execute(query)
	if err.NoError() {
		output.OnSuccess(result)
		return
	}
	output.OnError(out.ConvertQueryErrorDetail(err))

}

// ListEventAttendance implements AttendanceComponent.
func (a *attendanceComponentImpl) ListEventAttendance(ctx context.Context, query queries.ListEventAttendanceQuery, output out.Output[queries.ListEventAttendanceResult]) {
	result, err := a.dataLayer.QueryWorker().Query(ctx).ListEventAttendance().Execute(query)

	if err.NoError() {
		output.OnSuccess(result)
		return
	}
	output.OnError(out.ConvertQueryErrorDetail(err))

}

// ListEventSchedules implements AttendanceComponent.
func (a *attendanceComponentImpl) ListEventSchedules(ctx context.Context, query queries.ListEventScheduleQuery, output out.Output[queries.ListEventScheduleResult]) {

	result, err := a.dataLayer.QueryWorker().Query(ctx).ListEventSchedules().Execute(queries.ListEventScheduleQuery{Limit: query.Limit, LastID: query.LastID})
	if err.NoError() {
		output.OnSuccess(result)
		return
	}

	output.OnError(out.ConvertQueryErrorDetail(err))

}

// ListEventsBySchedule implements AttendanceComponent.
func (a *attendanceComponentImpl) ListEventsBySchedule(ctx context.Context, query queries.ListEventByScheduleQuery, output out.Output[queries.ListEventByScheduleResult]) {
	result, err := a.dataLayer.QueryWorker().Query(ctx).ListEventsBySchedule().Execute(query)
	if err.NoError() {
		output.OnSuccess(result)
		return
	}
	output.OnError(out.ConvertQueryErrorDetail(err))

}

// UpdateEventSchedule implements AttendanceComponent.
func (a *attendanceComponentImpl) UpdateEventSchedule(ctx context.Context, schedule dto.EventScheduleDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.UpdateEventScheduleCommand{
		ID:             schedule.ID,
		Name:           schedule.Name,
		ScheduleType:   string(schedule.Type),
		TimezoneOffset: schedule.TimezoneOffset,
		Days:           schedule.Days,
		Date:           schedule.Date,
		StartDate:      schedule.StartDate,
		EndDate:        schedule.EndDate,
	}).WithOutput(
		out.OutputAdapter(output, func(e entities.EventSchedule) dto.EventScheduleDTO {
			return dto.FromEntities(&e)
		}),
	).Execute(ctx)
}

// RemoveEventScheduleActivity implements AttendanceComponent.
func (a *attendanceComponentImpl) RemoveEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.RemoveScheduleActivityCommand{
		ScheduleID: activity.ScheduleID,
		ActivityID: activity.ID,
	}).WithOutput(
		out.OutputAdapter(output, func(e entities.EventSchedule) dto.EventScheduleDTO {
			return dto.FromEntities(&e)
		}),
	).Execute(ctx)
}

// AddEventScheduleActivity implements AttendanceComponent.
func (a *attendanceComponentImpl) AddEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.AddEventScheduleActivityCommand{
		ScheduleID: activity.ScheduleID,
		Name:       activity.Name,
		Hour:       activity.Hour,
		Minute:     activity.Minute,
	}).WithOutput(
		out.OutputAdapter(output, func(e entities.EventSchedule) dto.EventScheduleDTO {
			return dto.FromEntities(&e)
		}),
	).Execute(ctx)
}

// CreateEvent implements AttendanceComponent.
func (a *attendanceComponentImpl) CreateEvent(ctx context.Context, scheduleID string, output out.Output[dto.EventDTO]) out.Waitable {
	panic("unimplemented")
}

// CreateEventSchedule implements AttendanceComponent.
func (a *attendanceComponentImpl) CreateEventSchedule(ctx context.Context, schedule dto.EventScheduleDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.CreateEventScheduleCommand{
		Name:           schedule.Name,
		ScheduleType:   schedule.Type,
		TimezoneOffset: schedule.TimezoneOffset,
		Days:           schedule.Days,
		Date:           schedule.Date,
		StartDate:      schedule.StartDate,
		EndDate:        schedule.EndDate,
	}).WithOutput(
		out.OutputAdapter(output, func(e entities.EventSchedule) dto.EventScheduleDTO {
			return dto.FromEntities(&e)
		}),
	).Execute(ctx)
}

// UpdateEventScheduleActivity implements AttendanceComponent.
func (a *attendanceComponentImpl) UpdateEventScheduleActivity(ctx context.Context, activity dto.EventScheduleActivityDTO, output out.Output[dto.EventScheduleDTO]) out.Waitable {
	return utils.SingleCommandExecution(a.dataLayer.CommandWorker(), commands.UpdateEventScheduleActivityCommand{
		ScheduleID: activity.ScheduleID,
		ActivityID: activity.ID,
		Name:       activity.Name,
		Hour:       activity.Hour,
		Minute:     activity.Minute,
	}).WithOutput(
		out.OutputAdapter(output, func(e entities.EventSchedule) dto.EventScheduleDTO {
			return dto.FromEntities(&e)
		})).Execute(ctx)
}

// CheckIn implements AttendanceComponent.
func (a *attendanceComponentImpl) CheckIn(ctx context.Context, command commands.CheckInCommand, output out.Output[entities.Attendance]) out.Waitable {
	panic("unimplemented")
}

func NewAttendanceComponent(datalayer AttendanceDataLayerComponent) AttendanceComponent {
	return &attendanceComponentImpl{
		dataLayer: datalayer,
	}
}
