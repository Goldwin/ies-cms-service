package people

import (
	"log"

	"github.com/Goldwin/ies-pik-cms/internal/data"
	"github.com/Goldwin/ies-pik-cms/internal/data/people/mongo"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/Goldwin/ies-pik-cms/pkg/people"
	"github.com/Goldwin/ies-pik-cms/pkg/people/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/queries"
)

type peopleDataLayerComponentImpl struct {
	commandWorker worker.UnitOfWork[commands.CommandContext]
	queryWorker   worker.QueryWorker[queries.QueryContext]
}

// QueryWorker implements people.PeopleDataLayerComponent.
func (p *peopleDataLayerComponentImpl) QueryWorker() worker.QueryWorker[queries.QueryContext] {
	return p.queryWorker
}

// CommandWorker implements people.PeopleDataLayerComponent.
func (p *peopleDataLayerComponentImpl) CommandWorker() worker.UnitOfWork[commands.CommandContext] {
	return p.commandWorker
}

func NewPeopleDataLayerComponent(config data.DataLayerConfig, infra infra.InfraComponent) people.PeopleDataLayerComponent {
	if config.CommandConfig == nil {
		log.Fatalf("Command config is required for People Management Data Layer Component")
	}
	if config.CommandConfig.Mode != data.ModeMongo {
		log.Fatalf("Command mode %s is not supported for People Management Data Layer Component", config.CommandConfig.Mode)
	}
	return &peopleDataLayerComponentImpl{
		commandWorker: mongo.NewUnitOfWork(infra.Mongo(), config.CommandConfig.UseTransaction),
		queryWorker:   mongo.NewQueryWorker(infra.Mongo(), config.QueryConfig.UseTransaction),
	}
}
