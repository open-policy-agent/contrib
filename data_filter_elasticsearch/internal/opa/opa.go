// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aquasecurity/esquery"
	"github.com/open-policy-agent/opa/sdk"

	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/es"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

const defaultQuery = "data.example.allow == true"

// Result contains ES queries after partially evaluating OPA queries.
type Result struct {
	Defined bool
	Query   esquery.Mappable
}

// Compile compiles OPA query and partially evaluates it.
func Compile(ctx context.Context, opa *sdk.OPA, input map[string]interface{}) (Result, error) {

	unknowns := []string{"data.elastic"}

	options := sdk.PartialOptions{
		Now:      time.Now(),
		Input:    input,
		Query:    defaultQuery,
		Unknowns: unknowns,
		Mapper:   &sdk.RawMapper{},
	}

	p, err := opa.Partial(ctx, options)
	if err != nil {
		return Result{}, err
	}

	if len(p.AST.Queries) == 0 {
		// always deny
		return Result{Defined: false}, nil
	}

	for _, query := range p.AST.Queries {
		if len(query) == 0 {
			// always allow
			return Result{Defined: true}, nil
		}
	}

	return processQuery(p.AST)
}

func processQuery(pq *rego.PartialQueries) (Result, error) {

	queries := make([]esquery.Mappable, 0, 100)
	for i := range pq.Queries {
		exprQueries := make([]esquery.Mappable, 0, 100)
		for _, expr := range pq.Queries[i] {
			if !expr.IsCall() {
				continue
			}

			if len(expr.Operands()) != 2 {
				return Result{}, fmt.Errorf("invalid expression: too many arguments")
			}

			var value interface{}
			var processedTerm []string
			var err error
			for _, term := range expr.Operands() {
				if ast.IsConstant(term.Value) {
					value, err = ast.JSON(term.Value)
					if err != nil {
						return Result{}, fmt.Errorf("error converting term to JSON: %v", err)
					}
				} else {
					processedTerm = processTerm(term.String())
				}
			}

			var esQuery esquery.Mappable

			if isEqualityOperator(expr.Operator().String()) {
				// generate ES Term query
				esQuery = esquery.Term(processedTerm[1], value)

				// check if nested query
				terms := strings.Split(processedTerm[1], ".")
				if len(terms) > 1 {
					path := strings.Join(terms[:len(terms)-1], ".")
					esQuery = es.GenerateNestedQuery(path, esQuery)
				}

			} else if isRangeOperator(expr.Operator().String()) {
				// generate ES Range query
				if expr.Operator().String() == "lt" {
					esQuery = es.GenerateRangeQueryLt(processedTerm[1], value)
				} else if expr.Operator().String() == "gt" {
					esQuery = es.GenerateRangeQueryGt(processedTerm[1], value)
				} else if expr.Operator().String() == "lte" {
					esQuery = es.GenerateRangeQueryLte(processedTerm[1], value)
				} else {
					esQuery = es.GenerateRangeQueryGte(processedTerm[1], value)
				}
			} else if expr.Operator().String() == "neq" {
				// generate ES Must Not query
				esQuery = es.GenerateBoolMustNotQuery(processedTerm[1], value)
			} else if isContainsOperator(expr.Operator().String()) {
				// generate ES Query String query
				esQuery = es.GenerateQueryStringQuery(processedTerm[1], value)
			} else if isRegexpMatchOperator(expr.Operator().String()) {
				// generate ES Regexp query
				esQuery = es.GenerateRegexpQuery(processedTerm[1], value)
			} else {
				return Result{}, fmt.Errorf("invalid expression: operator not supported: %v", expr.Operator().String())
			}

			fmt.Printf("OPA Query #%d: %v\n", i+1, pq.Queries[i])
			strEsQuery, err := json.Marshal(esQuery.Map())
			if err != nil {
				return Result{}, err
			}
			fmt.Printf("ES  Query #%d: %+v\n\n", i+1, string(strEsQuery))
			exprQueries = append(exprQueries, esQuery)
		}

		if len(exprQueries) == 1 {
			queries = append(queries, exprQueries[0])
		} else {
			// ES queries generated within a rule are And'ed
			boolQuery := es.GenerateBoolFilterQuery(exprQueries)
			queries = append(queries, boolQuery)
		}
	}

	// ES queries generated from partial eval queries
	// are Or'ed
	combinedQuery := es.GenerateBoolShouldQuery(queries)
	return Result{Defined: true, Query: combinedQuery}, nil

}

// Eg. data.elastic.posts[_].<some_field>
// indexName => posts
// fieldName => some_field
func processTerm(query string) []string {
	splitQ := strings.Split(query, ".")
	result := []string{}
	for _, term := range splitQ {
		result = append(result, removeOpenBrace(term))
	}

	indexName := result[2]
	fieldName := result[3]
	if len(result) > 3 {
		fieldName = strings.Join(result[3:], ".")
	}

	return []string{indexName, fieldName}
}

func removeOpenBrace(input string) string {
	return strings.Split(input, "[")[0]
}

func isEqualityOperator(op string) bool {
	return op == "eq" || op == "equal"
}

func isContainsOperator(op string) bool {
	return op == "contains"
}

func isRegexpMatchOperator(op string) bool {
	return op == "re_match"
}

func isRangeOperator(op string) bool {
	return op == "lt" || op == "gt" || op == "lte" || op == "gte"
}
