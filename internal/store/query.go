package store

import "errors"

type QueryExecutor interface {
	Execute(Query) QueryResult
}

type Query interface {
}

type ReadQuery string
type WriteQuery string

type QueryResult struct {
	Data  []byte
	Error error
}

type Operation string

type ComparisonOperator string
type LogicalOperator string

const (
	insert     Operation = "insert"
	insertMany Operation = "insertMany"
	delete     Operation = "delete"
)

const (
	OpEqual        ComparisonOperator = "eq"
	OpNotEqual     ComparisonOperator = "ne"
	OpGreaterThan  ComparisonOperator = "gt"
	OpLessThan     ComparisonOperator = "lt"
	OpGreaterEqual ComparisonOperator = "gte"
	OpLessEqual    ComparisonOperator = "lte"
)

const (
	OpAnd LogicalOperator = "and"
	OpOr  LogicalOperator = "or"
)

var EmptyQueryResult = QueryResult{
	Data:  nil,
	Error: errors.New("No data found"),
}
