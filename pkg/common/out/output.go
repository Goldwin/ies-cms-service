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

type outputAdapter[FROM, TO any] struct {
	transformer    OutputTransformer[TO, FROM]
	originalOutput Output[FROM]
}

func (o *outputAdapter[FROM, TO]) OnSuccess(result TO) {
	transformedResult := o.transformer(result)
	o.originalOutput.OnSuccess(transformedResult)
}

func (o *outputAdapter[FROM, TO]) OnError(err AppErrorDetail) {
	o.originalOutput.OnError(err)
}

type OutputTransformer[FROM, TO any] func(FROM) TO

func OutputAdapter[FROM, TO any](output Output[FROM], transformer OutputTransformer[TO, FROM]) Output[TO]{
	return &outputAdapter[FROM, TO] {
		transformer: transformer,
		originalOutput: output,
	}
}

type Waitable interface {
	Wait()
}
