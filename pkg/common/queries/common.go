package queries

type QueryErrorCode int

type QueryErrorDetail struct {
	Code    QueryErrorCode `json:"code"`
	Message string         `json:"message"`
}

func (e QueryErrorDetail) Error() string {
	return e.Message
}

type QueryResult[T any] struct {
	Data *T `json:"data"`
}

const (
	ResourceNotFound QueryErrorCode = 404
)

var (
	NoQueryError = QueryErrorDetail{Code: 0, Message: ""}
)
