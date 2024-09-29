package queries

import "log"

type QueryErrorCode int

type QueryErrorDetail struct {
	Code    QueryErrorCode `json:"code"`
	Message string         `json:"message"`
}

func (e *QueryErrorDetail) Error() string {
	return e.Message
}

func (e *QueryErrorDetail) NoError() bool {
	return e == nil || e.Code == NoQueryError.Code
}

type QueryResult[T any] struct {
	Data *T `json:"data"`
}

// Common error code. coincides with http status codes
const (
	ResourceNotFound QueryErrorCode = 404
	InternalError    QueryErrorCode = 500
)

var (
	NoQueryError = QueryErrorDetail{Code: 0, Message: ""}
)

func ResourceNotFoundError(msg string) QueryErrorDetail {
	return QueryErrorDetail{Code: ResourceNotFound, Message: msg}
}

func InternalServerError(err error) QueryErrorDetail {
	log.Default().Printf("Failed to execute query: %s\n", err.Error())
	return QueryErrorDetail{Code: InternalError, Message: "Failed to execute query because of an internal error. please contact the system administrator."}
}
