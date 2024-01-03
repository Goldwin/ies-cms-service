package events

import (
	"log"

	"github.com/Goldwin/ies-pik-cms/internal/data"
	"github.com/Goldwin/ies-pik-cms/internal/data/events/mongo"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/Goldwin/ies-pik-cms/pkg/events"
	"github.com/Goldwin/ies-pik-cms/pkg/events/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/queries"
)

type churchEventDataLayerComponentImpl struct {
	commandWorker worker.UnitOfWork[commands.CommandContext]
	queryWorker   worker.QueryWorker[queries.QueryContext]
}

// QueryWorker implements events.ChurchDataLayerComponent.
func (c *churchEventDataLayerComponentImpl) QueryWorker() worker.QueryWorker[queries.QueryContext] {
	return c.queryWorker
}

// CommandWorker implements events.ChurchDataLayerComponent.
func (c *churchEventDataLayerComponentImpl) CommandWorker() worker.UnitOfWork[commands.CommandContext] {
	return c.commandWorker
}

func NewChurchEventDataLayerComponent(config data.DataLayerConfig, infra infra.InfraComponent) events.ChurchDataLayerComponent {
	if config.CommandConfig == nil {
		log.Fatalf("Command config is required for People Management Data Layer Component")
	}
	if config.CommandConfig.Mode != data.ModeMongo {
		log.Fatalf("Command mode %s is not supported for People Management Data Layer Component", config.CommandConfig.Mode)
	}
	return &churchEventDataLayerComponentImpl{
		commandWorker: mongo.NewUnitOfWork(infra.Mongo(), config.CommandConfig.DB, config.CommandConfig.UseTransaction),
		queryWorker:   mongo.NewQueryWorker(infra.Mongo(), config.QueryConfig.DB, config.QueryConfig.UseTransaction),
	}
}
