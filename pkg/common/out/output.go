package out

import (
	"github.com/Goldwin/ies-pik-cms/pkg/common/commands"
	"github.com/Goldwin/ies-pik-cms/pkg/common/queries"
)

type AppErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e AppErrorDetail) Error() string {
	return e.Message
}

func ConvertCommandErrorDetail(cmdError commands.CommandErrorDetail) AppErrorDetail {
	return AppErrorDetail{
		Code:    int(cmdError.Code),
		Message: cmdError.Message,
	}
}

func ConvertQueryErrorDetail(queryError queries.QueryErrorDetail) AppErrorDetail {
	return AppErrorDetail{
		Code:    int(queryError.Code),
		Message: queryError.Message,
	}
}

type Output[T any] interface {
	OnError(err AppErrorDetail)
	OnSuccess(result T)
}

type Waitable interface {
	Wait()
}