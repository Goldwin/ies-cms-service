package people

import (
	"context"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/common/worker"
	"github.com/Goldwin/ies-pik-cms/pkg/people/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/people/repositories"
)

type PeopleDataLayerComponent interface {
	CommandWorker() worker.UnitOfWork[repositories.CommandContext]
}

type PeopleManagementComponent interface {
	AddPerson(context.Context, dto.Person, out.Output[dto.Person])
	AddHousehold(context.Context, dto.HouseHoldInput, out.Output[dto.Household])
	UpdatePerson(context.Context, dto.Person, out.Output[dto.Person])
	UpdateHousehold(context.Context, dto.HouseHoldInput, out.Output[dto.Household])
}

func PeopleManagementComponents(worker worker.UnitOfWork[repositories.CommandContext]) PeopleManagementComponent {
	return &peopleManagementComponent{
		worker: worker,
	}
}

type peopleManagementComponent struct {
	worker worker.UnitOfWork[repositories.CommandContext]
}

// AddHousehold implements PeopleManagementComponent.
func (p *peopleManagementComponent) AddHousehold(ctx context.Context, input dto.HouseHoldInput, output out.Output[dto.Household]) {
	var result AppExecutionResult[dto.Household]
	_ = p.worker.Execute(ctx, func(ctx repositories.CommandContext) error {
		result = commands.AddHouseholdCommand{Input: input}.Execute(ctx)
		if result.Status == ExecutionStatusSuccess {
			go output.OnSuccess(result.Result)
		} else {
			go output.OnError(result.Error)
			return result.Error
		}
		return nil
	})
}

// AddPerson implements PeopleManagementComponent.
func (p *peopleManagementComponent) AddPerson(ctx context.Context, input dto.Person, output out.Output[dto.Person]) {
	var result AppExecutionResult[dto.Person]
	_ = p.worker.Execute(ctx, func(ctx repositories.CommandContext) error {
		result = commands.AddPersonCommand{Input: input}.Execute(ctx)
		if result.Status == ExecutionStatusSuccess {
			go output.OnSuccess(result.Result)
		} else {
			go output.OnError(result.Error)
			return result.Error
		}
		return nil
	})
}

// UpdateHousehold implements PeopleManagementComponent.
func (p *peopleManagementComponent) UpdateHousehold(ctx context.Context, input dto.HouseHoldInput, output out.Output[dto.Household]) {
	var result AppExecutionResult[dto.Household]
	_ = p.worker.Execute(ctx, func(ctx repositories.CommandContext) error {
		result = commands.UpdateHouseholdCommand{Input: input}.Execute(ctx)
		if result.Status == ExecutionStatusSuccess {
			go output.OnSuccess(result.Result)
		} else {
			go output.OnError(result.Error)
			return result.Error
		}
		return nil
	})
}

// UpdatePerson implements PeopleManagementComponent.
func (p *peopleManagementComponent) UpdatePerson(ctx context.Context, input dto.Person, output out.Output[dto.Person]) {
	var result AppExecutionResult[dto.Person]
	_ = p.worker.Execute(ctx, func(ctx repositories.CommandContext) error {
		result = commands.UpdatePersonCommand{Input: input}.Execute(ctx)
		if result.Status == ExecutionStatusSuccess {
			go output.OnSuccess(result.Result)
		} else {
			go output.OnError(result.Error)
			return result.Error
		}
		return nil
	})
}

func NewPeopleManagementComponent(data PeopleDataLayerComponent) PeopleManagementComponent {
	return &peopleManagementComponent{
		worker: data.CommandWorker(),
	}
}
