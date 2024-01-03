package worker

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/commands"
)

type AtomicOperation func(commands.CommandContext) error

type PeopleManagementUnitOfWork interface {
	Execute(ctx context.Context, op AtomicOperation) error
}
