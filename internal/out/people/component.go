package people

import (
	"github.com/Goldwin/ies-pik-cms/internal/bus"
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/people/dto"
)

type PeopleManagementOutputComponent interface {
	AddPersonOutput() out.Output[dto.Person]
}

type peopleManagementOutputComponent struct {
	addPersonOutput out.Output[dto.Person]
}

// AddPersonOutput implements PeopleManagementOutputComponent.
func (p *peopleManagementOutputComponent) AddPersonOutput() out.Output[dto.Person] {
	return p.addPersonOutput
}

func NewPeopleManagementOutputComponent(bus bus.EventBusComponent) PeopleManagementOutputComponent {
	return &peopleManagementOutputComponent{
		addPersonOutput: AddPersonOutput(bus),
	}
}
