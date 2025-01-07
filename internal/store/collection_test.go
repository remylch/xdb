package store

import (
	"testing"
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

	if testCollection.Name != "test" {
		t.Error("Expected collection name 'test', got", testCollection.Name)
	}

	if len(testCollection.Indexes) == 0 {
		t.Error("Expected at least one index with id in the collection")
	}

}
