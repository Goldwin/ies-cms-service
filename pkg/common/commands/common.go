package commands

import "log"

type ExecutionStatus string
type CommandErrorCode int

type CommandExecutionResult[T any] struct {
	Status ExecutionStatus
	Error  CommandErrorDetail
	Result T
}

type Command[CTX any, R any] interface {
	Execute(ctx CTX) CommandExecutionResult[R]
}
type CommandErrorDetail struct {
	Code    CommandErrorCode `json:"code"`
	Message string           `json:"message"`
	Details []string         `json:"detail"`
}

func (e CommandErrorDetail) Error() string {
	return e.Message
}

func CommandErrorDetailWorkerFailure(err error) CommandErrorDetail {
	log.Default().Printf("Failed to execute command: %s\n", err.Error())
	return CommandErrorDetail{Code: CommandErrorCodeWorkerFailure, Message: "Failed to execute command because of an internal error. please contact the system administrator."}
}

const (
	ExecutionStatusSuccess ExecutionStatus = "SUCCESS"
	ExecutionStatusFailed  ExecutionStatus = "FAILED"

	CommandErrorCodeNone          CommandErrorCode = 0
	CommandErrorCodeWorkerFailure CommandErrorCode = 500
)

var (
	CommandErrorDetailNone = CommandErrorDetail{Code: CommandErrorCodeNone, Message: ""}
)
