package attendance

import (
	"github.com/Goldwin/ies-pik-cms/internal/data"
	"github.com/Goldwin/ies-pik-cms/internal/data/attendance/mongo"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/attendance/queries"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type attendanceDataLayerComponentImpl struct {
	commandWorker worker.UnitOfWork[commands.CommandContext]
	queryWorker   worker.QueryWorker[queries.QueryContext]
}

// CommandWorker implements attendance.AttendanceDataLayerComponent.
func (a *attendanceDataLayerComponentImpl) CommandWorker() worker.UnitOfWork[commands.CommandContext] {
	return a.commandWorker
}

// QueryWorker implements attendance.AttendanceDataLayerComponent.
func (a *attendanceDataLayerComponentImpl) QueryWorker() worker.QueryWorker[queries.QueryContext] {
	return a.queryWorker
}

func NewAttendanceDataLayerComponent(config data.DataLayerConfig, infraComponent infra.InfraComponent) attendance.AttendanceDataLayerComponent {
	return &attendanceDataLayerComponentImpl{
		commandWorker: mongo.NewUnitOfWork(infraComponent.Mongo(), config.CommandConfig.UseTransaction),
		queryWorker:   mongo.NewQueryWorker(infraComponent.Mongo(), config.QueryConfig.UseTransaction),
	}
}
