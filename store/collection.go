package store

import (
	"fmt"

	"github.com/google/uuid"
)

type JSON map[string]interface{}

type Collection struct {
	id        uuid.UUID // Since it's generated at the instantiation of the collection(on node startup) => 2 nodes will have different uuid for the same collection
	name      string
	indexes   []Index
	documents []Document
}

// Document is a part of the collection.  It's a certain array of bytes into the collection file
type Document struct {
	XDBId uuid.UUID
	Data  JSON
}

func newCollection(name string) *Collection {

	return &Collection{
		id:        uuid.New(),
		name:      name,
		indexes:   make([]Index, 0),
		documents: make([]Document, 0),
	}
}

func (c *Collection) addDocument(data JSON) uuid.UUID {
	docId := uuid.New()
	c.documents = append(c.documents, Document{
		XDBId: docId,
		Data:  data,
	})
	return docId
}

func (c *Collection) getDocument(docId uuid.UUID) (Document, error) {
	//TODO: find better way to find the document fast
	for _, doc := range c.documents {
		if doc.XDBId == docId {
			return doc, nil
		}
	}

	return Document{}, fmt.Errorf("document with id %s not found in collection %s", docId, c.name)
}

func (c *Collection) delete(docId uuid.UUID) error {
	idx := -1

	for i, doc := range c.documents {
		if doc.XDBId == docId {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fmt.Errorf("document with id %s not found in collection %s", docId, c.name)
	}

	c.documents[idx] = c.documents[len(c.documents)-1]
	c.documents = c.documents[:len(c.documents)-1]
	return nil
}
