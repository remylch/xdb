package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	fixture = JSON{
		"id":   "1",
		"name": "John Doe",
		"age":  50,
		"nestedObject": JSON{
			"name": "object",
		},
		"nestedArray": []interface{}{
			JSON{
				"id": "item1",
			},
			JSON{
				"id": "item2",
			},
			JSON{
				"id": "item3",
			},
		},
	}
)

func TestNewCollection(t *testing.T) {
	testCollection := newCollection("test")

	if testCollection.name != "test" {
		t.Error("Expected collection name 'test', got", testCollection.name)
	}

	if len(testCollection.indexes) == 0 {
		t.Error("Expected at least one index with id in the collection")
	}

}

func TestCollection_AddDocument(t *testing.T) {
	testCollection := newCollection("test")
	_ = testCollection.addDocument(fixture)
	require.Len(t, testCollection.documents, 1, "doc array should have 3 items")
}

func TestCollection_GetDocument(t *testing.T) {
	testCollection := newCollection("test")
	docId := testCollection.addDocument(fixture)
	document, err := testCollection.getDocument(docId)
	if err != nil {
		t.Errorf("err retrieving document : %v", err)
	}
	require.Equal(t, document.XDBId, docId, "XDB doc id should match")
	require.Equal(t, fixture["id"], document.Data["id"], "doc id should match")
}

func TestCollection_DeleteDocument(t *testing.T) {
	testCollection := newCollection("test")
	docId := testCollection.addDocument(fixture)
	require.Len(t, testCollection.documents, 1, "collection should have 1 document")
	require.NoError(t, testCollection.delete(docId), "err deleting document")
	doc, err := testCollection.getDocument(docId)
	require.Error(t, err, "collection should not contain the deleted document")
	require.Equal(t, doc, Document{}, "document should be empty")
}
