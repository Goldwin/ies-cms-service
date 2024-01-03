package auth

import (
	"log"

	. "github.com/Goldwin/ies-pik-cms/internal/data"
	"github.com/Goldwin/ies-pik-cms/internal/data/auth/local"
	"github.com/Goldwin/ies-pik-cms/internal/data/auth/redis"
	"github.com/Goldwin/ies-pik-cms/internal/infra"
	"github.com/Goldwin/ies-pik-cms/pkg/auth"
	"github.com/Goldwin/ies-pik-cms/pkg/auth/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
)

type authDataLayerComponentImpl struct {
	commandWorker worker.UnitOfWork[commands.CommandContext]
}

// CommandWorker implements auth.AuthDataLayerComponent.
func (a *authDataLayerComponentImpl) CommandWorker() worker.UnitOfWork[commands.CommandContext] {
	return a.commandWorker
}

func NewAuthDataLayerComponent(config DataLayerConfig, infra infra.InfraComponent) auth.AuthDataLayerComponent {
	component := &authDataLayerComponentImpl{}
	if config.CommandConfig == nil {
		log.Fatalf("Command config is required")
	}
	switch config.CommandConfig.Mode {
	case ModeRedis:
		component.commandWorker = redis.NewUnitOfWork(infra.Redis())
	case ModeLocal:
		component.commandWorker = local.NewUnitOfWork()
	default:
		log.Fatalf("Command mode %s is not supported", config.CommandConfig.Mode)
	}
	return component
}
