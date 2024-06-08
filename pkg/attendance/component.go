package attendance

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type AttendanceDataLayerComponent interface {
	QueryWorker() worker.QueryWorker[QueryContext]
	CommandWorker() worker.UnitOfWork[commands.CommandContext]
}

type AttendanceComponent interface {
	CreateEventSchedule(ctx context.Context, schedule EventSchedule)
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
func (a *attendanceComponentImpl) CreateEventSchedule(ctx context.Context, schedule EventSchedule) {
	panic("unimplemented")
}

func NewAttendanceComponent(datalayer AttendanceDataLayerComponent) AttendanceComponent {
	return &attendanceComponentImpl{
		dataLayer: datalayer,
	}
}
