package indexing

type EntityStatus int

const (
	EntityStatusPending EntityStatus = iota
	EntityStatusEmbedding
	EntityStatusEmbedded
)

type Entity struct {
	name   string
	status EntityStatus
}

func NewEntity(name string) *Entity {
	return &Entity{
		name:   name,
		status: EntityStatusPending,
	}
}

// UnmarshalEntityFromDatabase unmarshals Entity from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalEntityFromDatabase as constructor - It may put domain into the invalid state!
func UnmarshalEntityFromDatabase(name string, status EntityStatus) *Entity {
	return &Entity{
		name:   name,
		status: status,
	}
}

func (e *Entity) Name() string {
	return e.name
}

func (e *Entity) Status() EntityStatus {
	return e.status
}
