package commands

// IDGenerator is an interface abstraction for generating unique identifiers
type IDGenerator interface {
	Generate(key string) int64
}

// UUIDGenerator is an interface abstraction for generating universal unique identifiers
type UUIDGenerator interface {
	GenerateUUID() string
}