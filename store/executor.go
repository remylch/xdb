package store

import (
	"fmt"
	"os"
	"path"
)

type BaseExecutor struct {
	dataDir           string
	dataBlockManager  DataBlockManager
	getCollectionHash func(string) string
}

func NewBaseExecutor(manager DataBlockManager, dataDir string, getCollectionHash func(string) string) QueryExecutor {
	return &BaseExecutor{
		dataBlockManager:  manager,
		dataDir:           dataDir,
		getCollectionHash: getCollectionHash,
	}
}

func (executor BaseExecutor) Execute(query Query) QueryResult {
	fmt.Printf("[Base - Executor] Executing query: %v\n", query)
	switch q := query.(type) {
	case WriteQuery:
		return QueryResult{}
	case ReadQuery:
		return executor.simpleReadQuery(q)
	default:
		return QueryResult{
			Data:  nil,
			Error: fmt.Errorf("unsupported query type: %T", query),
		}
	}
}

// For now we assume query is only the collection name
func (executor BaseExecutor) simpleReadQuery(query ReadQuery) QueryResult {
	//TODO: parse query
	collection := string(query)

	dataBlocks := make([]DataBlock, 0)

	collectionPath := path.Join(executor.dataDir, executor.getCollectionHash(collection))

	collectionDir, err := os.ReadDir(collectionPath)

	if err != nil {
		return QueryResult{
			Data:  nil,
			Error: fmt.Errorf("unable to read collection < %s > directory", query),
		}
	}

	for _, file := range collectionDir {
		infos, _ := file.Info()
		filePath := path.Join(collectionPath, infos.Name())

		newDataBlocks, err := executor.dataBlockManager.ReadDataBlock(filePath)

		if err != nil {
			return QueryResult{
				Data:  nil,
				Error: err,
			}
		}

		//TODO: dataBlocks should be only one part of the total data || for now dataBlocks is all the data of this collection
		dataBlocks = append(dataBlocks, newDataBlocks...)

	}

	merge := executor.dataBlockManager.mergeCorrelatedDataBlocks(dataBlocks)

	return QueryResult{
		Data:  merge,
		Error: nil,
	}
}

type AIExecutor struct {
}

func NewAIExecutor() QueryExecutor {
	return &AIExecutor{}
}

// TODO: Implement
func (executor AIExecutor) Execute(query Query) QueryResult {
	fmt.Printf("[AI - Executor] Executing query: %v\n", query)
	return QueryResult{}
}
