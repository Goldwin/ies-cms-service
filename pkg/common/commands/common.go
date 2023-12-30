package commands

type ExecutionStatus string
type AppErrorCode int

type AppExecutionResult[T any] struct {
	Status ExecutionStatus
	Error  AppErrorDetail
	Result T
}

type AppErrorDetail struct {
	Code    AppErrorCode `json:"code"`
	Message string       `json:"message"`
}

func (e AppErrorDetail) Error() string {
	return e.Message
}

func AppErrorDetailWorkerFailure(err error) AppErrorDetail {
	return AppErrorDetail{Code: AppErrorCodeWorkerFailure, Message: err.Error()}
}

const (
	ExecutionStatusSuccess ExecutionStatus = "SUCCESS"
	ExecutionStatusFailed  ExecutionStatus = "FAILED"

	AppErrorCodeNone          AppErrorCode = 0
	AppErrorCodeWorkerFailure AppErrorCode = 1
)

var (
	AppErrorDetailNone = AppErrorDetail{Code: AppErrorCodeNone, Message: ""}
)
