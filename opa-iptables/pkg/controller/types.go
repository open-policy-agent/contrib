package controller

type operation string

const (
	insertOp operation = "insert"
	deleteOp operation = "delete"
	testOp   operation = "test"
)

type Payload struct {
	QueryPath string      `json:"query_path"`
	Input     interface{} `json:"input"`
	Op        operation   `json:"operation"`
}
