package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/open-policy-agent/contrib/data_filter_mongodb/internal/mongo"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"go.uber.org/zap"
)

var (
	unknowns     = []string{"data.employees"}
	defaultQuery = "data.details.authz.allow == true"
)

type opaClient struct {
	logger      *zap.Logger
	mongoClient mongo.DBClient
	policyFile  string
}

type requestBody struct {
	Input input `json:"input"`
}

type input struct {
	Method string   `json:"method"`
	Path   []string `json:"path"`
	User   string   `json:"user"`
}

// handle http request
func (c *opaClient) queryMongoOPA(r *http.Request) (mongo.Result, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.logger.Error("Failed to read the http request", zap.Error(err))
		return mongo.Result{}, err
	}

	var inputBody requestBody
	err = json.Unmarshal(body, &inputBody)
	if err != nil {
		c.logger.Error("Failed to unmarshal the http request", zap.Error(err))
		return mongo.Result{}, err
	}

	input := map[string]interface{}{
		"method": r.Method,
		"path":   inputBody.Input.Path,
		"user":   inputBody.Input.User,
	}

	c.logger.Info("received request", zap.Any("request", input))
	return c.Compile(r.Context(), input)
}

// Compile compiles OPA query and partially evaluates it.
func (c *opaClient) Compile(ctx context.Context, input map[string]interface{}) (mongo.Result, error) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		c.logger.Error("Failed to marshal the request in compile phase.", zap.Error(err))
		return mongo.Result{}, fmt.Errorf("JSON Encoding error %v", err)
	}

	inputTerm, err := ast.ParseTerm(string(inputBytes))
	if err != nil {
		c.logger.Error("Failed to parse term.", zap.Error(err))
		return mongo.Result{}, err
	}

	// load policy
	policy, err := ioutil.ReadFile(c.policyFile)
	if err != nil {
		c.logger.Error("Failed to read policy file", zap.Error(err))
		return mongo.Result{}, fmt.Errorf("failed to read policy: %s", err)
	}

	r := rego.New(
		rego.Query(defaultQuery),
		rego.Module(c.policyFile, string(policy)),
		rego.ParsedInput(inputTerm.Value),
		rego.Unknowns(unknowns),
	)

	pq, err := r.Partial(ctx)
	if err != nil {
		c.logger.Error("Failed to perform opa partial eval", zap.Error(err))
		return mongo.Result{}, err
	}

	if len(pq.Queries) == 0 {
		// always deny
		return mongo.Result{Defined: false}, nil
	}

	for query := range pq.Queries {
		if len(pq.Queries[query]) == 0 {
			// always allow
			return mongo.Result{Defined: true}, nil
		}
	}

	return c.processQuery(pq)
}

func (c *opaClient) processQuery(pq *rego.PartialQueries) (mongo.Result, error) {
	var queries []mongo.Queries
	c.logger.Info("opa-query: ", zap.Any("queries", fmt.Sprintf("%s", pq.Queries)))
	for i := range pq.Queries {
		pipeline := &[]bson.M{}
		for _, expr := range pq.Queries[i] {
			if !expr.IsCall() {
				continue
			}

			if len(expr.Operands()) != 2 {
				return mongo.Result{}, fmt.Errorf("invalid expression: too many arguments")
			}

			var value interface{}
			var processedTerm []string
			var err error
			for _, term := range expr.Operands() {
				if ast.IsConstant(term.Value) {
					value, err = ast.JSON(term.Value)
					if err != nil {
						return mongo.Result{}, fmt.Errorf("error converting term to JSON: %v", err)
					}
				} else {
					processedTerm = processTerm(term.String())
				}
			}

			if processedTerm == nil {
				return mongo.Result{}, nil
			}

			if isEqualityOperator(expr.Operator().String()) {
				mongo.HandleEquals(pipeline, processedTerm[1], value)
			} else if isRangeOperator(expr.Operator().String()) {
				if expr.Operator().String() == "lt" {
					mongo.HandleLessThan(pipeline, processedTerm[1], value)
				} else if expr.Operator().String() == "gt" {
					mongo.HandleGreaterThan(pipeline, processedTerm[1], value)
				} else if expr.Operator().String() == "lte" {
					mongo.HandleLessThanEquals(pipeline, processedTerm[1], value)
				} else if expr.Operator().String() == "gte" {
					mongo.HandleGreaterThanEquals(pipeline, processedTerm[1], value)
				}
			} else if expr.Operator().String() == "neq" {
				mongo.HandleNotEquals(pipeline, processedTerm[1], value)
			} else {
				return mongo.Result{}, fmt.Errorf("invalid expression: operator not supported: %v", expr.Operator().String())
			}
		}
		k1 := mongo.Queries{Pipeline: bson.M{"$and": *pipeline}}
		queries = append(queries, k1)
	}
	result, err := c.mongoClient.QueryMongo(queries)
	if err != nil {
		return result, err
	}

	return result, nil
}

func processTerm(query string) []string {
	splitQ := strings.Split(query, ".")
	var result []string
	for _, term := range splitQ {
		result = append(result, removeOpenBrace(term))
	}

	if result == nil {
		return nil
	}

	indexName := result[1]
	fieldName := result[2]
	if len(result) > 2 {
		fieldName = strings.Join(result[2:], ".")
	}

	return []string{indexName, fieldName}
}

func removeOpenBrace(input string) string {
	return strings.Split(input, "[")[0]
}

func isEqualityOperator(op string) bool {
	return op == "eq" || op == "equal"
}

func isRangeOperator(op string) bool {
	return op == "lt" || op == "gt" || op == "lte" || op == "gte"
}
