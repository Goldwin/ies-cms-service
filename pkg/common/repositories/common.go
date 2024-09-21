package repositories

/*
Common repository is an interface abstraction for accessing entities from the persistence layer,
and updating the state of entities to the data layer

ID - Entity identifier type
T - Entity's model as an abstract data type
*/
type Repository[ID comparable, T any] interface {
	ReadOnlyRepository[ID, T]
	WriteOnlyRepository[ID, T]
}

type WriteOnlyRepository[ID comparable, T any] interface {
	/*
		Add entities model to the data layer
	*/
	Save(*T) (*T, error)

	/*
		Delete a single entity from the persistence layer
	*/
	Delete(*T) error
}

type ReadOnlyRepository[ID comparable, T any] interface {
	/*
		Get a single entity by it's identifier
	*/
	Get(ID) (*T, error)

	/*
		List multiple entites by a list of identifier
		The ordering of the slice will follow the exact sequence of the given identifier
	*/
	List([]ID) ([]*T, error)
}
