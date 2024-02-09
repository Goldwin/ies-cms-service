package commands

type ExecutionStatus string
type CommandErrorCode int

type CommandExecutionResult[T any] struct {
	Status ExecutionStatus
	Error  CommandErrorDetail
	Result T
}

type CommandErrorDetail struct {
	Code    CommandErrorCode `json:"code"`
	Message string           `json:"message"`
}

func (e CommandErrorDetail) Error() string {
	return e.Message
}

func CommandErrorDetailWorkerFailure(err error) CommandErrorDetail {
	return CommandErrorDetail{Code: CommandErrorCodeWorkerFailure, Message: err.Error()}
}

const (
	ExecutionStatusSuccess ExecutionStatus = "SUCCESS"
	ExecutionStatusFailed  ExecutionStatus = "FAILED"

	CommandErrorCodeNone          CommandErrorCode = 0
	CommandErrorCodeWorkerFailure CommandErrorCode = 1
)

var (
	CommandErrorDetailNone = CommandErrorDetail{Code: CommandErrorCodeNone, Message: ""}
)
