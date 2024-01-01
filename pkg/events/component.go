package events

import (
	"github.com/Goldwin/ies-pik-cms/pkg/common/out"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
)

type ChurchEventComponent interface {
	CheckIn(input dto.CheckInEvent, output out.Output[dto.CheckInEvent])
	SaveEvent(input dto.ChurchEvent, output out.Output[dto.ChurchEvent])
	SaveEventSchedule(input dto.ChurchEventSchedule, output out.Output[dto.ChurchEventSchedule])
}
