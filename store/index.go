package store

import (
	"fmt"
	"github.com/google/uuid"
	"regexp"
)

var (
	INDEX_RGX = "^index-([0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12})-([^0-9A-Fa-f])+"
)

type Index struct {
	Id   uuid.UUID
	Name string
}

func CreateIndex(name string) (*Index, error) {
	indexId := uuid.New()
	indexName := "index-" + indexId.String() + "-" + name

	match, _ := regexp.MatchString(INDEX_RGX, indexName)

	if !match {
		return nil, fmt.Errorf("invalid index name: %s", indexName)
	}

	return &Index{
		Id:   indexId,
		Name: name,
	}, nil
}

func (index Index) GetFullIndexName() string {
	return fmt.Sprintf("index-%s-%s", index.Id, index.Name)
}
