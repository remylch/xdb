package store

import (
	"github.com/google/uuid"
)

type JSON map[string]interface{}

type Collection struct {
	id      uuid.UUID // FIXME: Since it's generated at the instantiation of the collection(on node startup) => 2 nodes will have different uuid for the same collection
	name    string
	indexes []Index
}

func defaultIndex() Index {
	return Index("id")
}

func newCollection(name string) *Collection {
	return &Collection{
		id:      uuid.New(),
		name:    name,
		indexes: []Index{defaultIndex()},
	}
}
