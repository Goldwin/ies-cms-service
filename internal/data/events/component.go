package events

import (
	"log"

	"github.com/Goldwin/ies-pik-cms/internal/data"
	"github.com/Goldwin/ies-pik-cms/internal/data/events/mongo"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/Goldwin/ies-pik-cms/pkg/events"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
)

type churchEventDataLayerComponentImpl struct {
	commandWorker worker.UnitOfWork[repositories.CommandContext]
}

// CommandWorker implements events.ChurchDataLayerComponent.
func (c *churchEventDataLayerComponentImpl) CommandWorker() worker.UnitOfWork[repositories.CommandContext] {
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
		commandWorker: mongo.NewUnitOfWork(infra.Mongo(), config.CommandConfig.UseTransaction),
	}
}
