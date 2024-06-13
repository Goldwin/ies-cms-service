package attendance

import (
	"github.com/Goldwin/ies-pik-cms/internal/data"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type attendanceDataLayerComponentImpl struct {
}

// CommandWorker implements attendance.AttendanceDataLayerComponent.
func (a *attendanceDataLayerComponentImpl) CommandWorker() worker.UnitOfWork[commands.CommandContext] {
	panic("unimplemented")
}

// QueryWorker implements attendance.AttendanceDataLayerComponent.
func (a *attendanceDataLayerComponentImpl) QueryWorker() worker.QueryWorker[attendance.QueryContext] {
	panic("unimplemented")
}

func NewAttendanceDataLayerComponent(config data.DataLayerConfig, infraComponent infra.InfraComponent) attendance.AttendanceDataLayerComponent {
	return &attendanceDataLayerComponentImpl{}
}
