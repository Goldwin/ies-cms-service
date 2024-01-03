package queries

type Query[T any, R any] interface {
	Execute(ctx T) (R, error)
}
