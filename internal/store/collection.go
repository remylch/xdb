package store

import (
	"github.com/google/uuid"
)

type JSON map[string]interface{}

type Collection struct {
	Id      uuid.UUID `json:"id"` // FIXME: Since it's generated at the instantiation of the collection(on node startup) => 2 nodes will have different uuid for the same collection
	Name    string    `json:"name"`
	Indexes []Index   `json:"indexes"`
}

func defaultIndex() Index {
	return Index("id")
}

func newCollection(name string) Collection {
	// Make a copy of the name string to ensure no external references
	nameCopy := string([]byte(name))
	return Collection{
		Id:      uuid.New(),
		Name:    nameCopy,
		Indexes: []Index{defaultIndex()},
	}
}
