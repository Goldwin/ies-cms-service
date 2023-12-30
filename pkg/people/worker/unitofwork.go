package worker

import (
	"context"

	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
)

type AtomicOperation func(repositories.CommandContext) error

type PeopleManagementUnitOfWork interface {
	Execute(ctx context.Context, op AtomicOperation) error
}
