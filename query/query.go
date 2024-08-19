package query

type QueryExecutor interface {
	Execute(any) any
}

type QueryOperation string

const (
	insert     QueryOperation = "insert"
	insertMany QueryOperation = "insertMany"
	delete     QueryOperation = "delete"
)

type Query struct {
	Operation QueryOperation
	Value     interface{}
}

type LogicalQuery struct {
	Operator LogicalOperator
	Queries  []Query
}

type ComparisonOperator string

const (
	OpEqual        ComparisonOperator = "eq"
	OpNotEqual     ComparisonOperator = "ne"
	OpGreaterThan  ComparisonOperator = "gt"
	OpLessThan     ComparisonOperator = "lt"
	OpGreaterEqual ComparisonOperator = "gte"
	OpLessEqual    ComparisonOperator = "lte"
)

type LogicalOperator string

const (
	OpAnd LogicalOperator = "and"
	OpOr  LogicalOperator = "or"
)
