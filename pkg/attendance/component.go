package attendance

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/attendance/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type AttendanceDataLayerComponent interface {
	QueryWorker() worker.QueryWorker[QueryContext]
	CommandWorker() worker.UnitOfWork[CommandContext]
}

type AttendanceComponent interface {
	CreateEventSchedule(ctx context.Context, schedule EventSchedule)
	AddEventScheduleActivity(ctx context.Context, activity EventScheduleActivity)
	CreateEvent(ctx context.Context, event Event)
	CreateEventActivity(ctx context.Context, activity EventActivity)
}
