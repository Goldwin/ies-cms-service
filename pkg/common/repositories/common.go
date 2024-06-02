package repositories

/*
Common repository is an interface abstraction for accessing entities from the persistence layer, 
and updating the state of entities to the data layer

ID - Entity identifier type
T - Entity's model as an abstract data type
*/
type Repository[ID comparable, T any] interface {
	/*
	Add entities model to the data layer if it doesn't exists yet,
	Or updating the state of the entity if it's already exists.
	*/
	Save(*T) (*T, error)

	/*
	Get a single entity by it's identifier
	*/
	Get(ID) (*T, error)

	/*
	Delete a single entity from the persistence layer
	*/
	Delete(*T) error

	/*
	List multiple entites by a list of identifier
	The ordering of the slice will follow the exact sequence of the given identifier
	*/
	List([]ID) ([]*T, error)
}