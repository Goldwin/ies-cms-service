package commands

import (
	"fmt"
	"time"

	. "github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/events/dto"
	"github.com/Goldwin/ies-pik-cms/pkg/events/entities"
	"github.com/Goldwin/ies-pik-cms/pkg/events/repositories"
)

const (
	FetchSessionErrorFailedToCreateSession AppErrorCode = 30201
)

type CreateChurchEventSessionCommand struct {
	EventID string
}

func (cmd CreateChurchEventSessionCommand) Execute(ctx repositories.CommandContext) AppExecutionResult[dto.ChurchEventSession] {
	events, err := ctx.ChurchEventRepository().Get(cmd.EventID)
	if err != nil {
		return AppExecutionResult[dto.ChurchEventSession]{
			Status: ExecutionStatusFailed,
			Error:  AppErrorDetail{Code: FetchSessionErrorFailedToCreateSession, Message: err.Error()}}
	}

	dayAddition := 1
	if events.EventFrequency == "WEEKLY" {
		dayAddition = 7
	}
	session := entities.ChurchEventSession{
		ID:        fmt.Sprintf("%s.%d", events.ID, events.LatestSessionNo+1),
		SessionNo: events.LatestSessionNo + 1,
		StartTime: events.LatestShowAt.Add(time.Hour*24 + time.Duration(dayAddition)),
		ShowAt:    events.LatestShowAt.Add(time.Hour*24 + time.Duration(dayAddition)),
		HideAt:    events.LatestHideAt.Add(time.Hour*24 + time.Duration(dayAddition)),
	}

	err = ctx.ChurchEventSessionRepository().Save(session)
	if err != nil {
		return AppExecutionResult[dto.ChurchEventSession]{
			Status: ExecutionStatusFailed,
			Error:  AppErrorDetail{Code: FetchSessionErrorFailedToCreateSession, Message: err.Error()}}
	}

	events.LatestHideAt = session.HideAt
	events.LatestSessionNo++
	events.LatestShowAt = session.ShowAt
	events.LatestSessionStartTime = session.StartTime

	err = ctx.ChurchEventRepository().Save(*events)

	return AppExecutionResult[dto.ChurchEventSession]{
		Status: ExecutionStatusSuccess,
		Result: dto.ChurchEventSession{
			ID:        session.ID,
			Name:      session.Name,
			SessionNo: session.SessionNo,
			StartTime: session.StartTime,
			ShowAt:    session.ShowAt,
			HideAt:    session.HideAt,
		},
	}
}
