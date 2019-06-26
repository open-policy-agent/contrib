package controller

import (
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/iptables"
)

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

type Result struct {
	Rules []iptables.Rule `json:"result"`
}
