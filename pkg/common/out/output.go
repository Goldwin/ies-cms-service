package out

import "github.com/Goldwin/ies-pik-cms/pkg/common/commands"

type Output[T any] interface {
	OnError(err commands.AppErrorDetail)
	OnSuccess(result T)
}
