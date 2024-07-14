package common

// Component is a generic interface abstraction of deployable business logic modules.
type Component interface {
	Start()
	Stop()
}
