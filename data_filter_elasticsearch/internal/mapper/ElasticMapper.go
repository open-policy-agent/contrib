package mapper

import (
	"fmt"
	"strings"

	"github.com/olivere/elastic"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

type ElasticMapper struct {
}

type Result struct {
	Defined bool
	Query   elastic.Query
}

func (e *ElasticMapper) MapResults(pq *rego.PartialQueries) (interface{}, error) {
	if len(pq.Queries) == 0 {
		// always deny
		return Result{Defined: false}, nil
	}

	for _, query := range pq.Queries {
		if len(query) == 0 {
			// always allow
			return Result{Defined: true}, nil
		}
	}
	result, err := processQuery(pq)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (e *ElasticMapper) ResultToJSON(results interface{}) (interface{}, error) {
	result := results.(Result)
	log, err := result.Query.Source()
	return map[string]interface{}{"Defined": result.Defined, "Query": log}, err
}

func processQuery(pq *rego.PartialQueries) (Result, error) {

	queries := []elastic.Query{}
	for i := range pq.Queries {
		exprQueries := []elastic.Query{}
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

			var esQuery elastic.Query

			if isEqualityOperator(expr.Operator().String()) {
				// generate ES Term query
				esQuery = GenerateTermQuery(processedTerm[1], value)

				// check if nested query
				terms := strings.Split(processedTerm[1], ".")
				if len(terms) > 1 {
					path := strings.Join(terms[:len(terms)-1], ".")
					esQuery = GenerateNestedQuery(path, esQuery)
				}

			} else if isRangeOperator(expr.Operator().String()) {
				// generate ES Range query
				if expr.Operator().String() == "lt" {
					esQuery = GenerateRangeQueryLt(processedTerm[1], value)
				} else if expr.Operator().String() == "gt" {
					esQuery = GenerateRangeQueryGt(processedTerm[1], value)
				} else if expr.Operator().String() == "lte" {
					esQuery = GenerateRangeQueryLte(processedTerm[1], value)
				} else {
					esQuery = GenerateRangeQueryGte(processedTerm[1], value)
				}
			} else if expr.Operator().String() == "neq" {
				// generate ES Must Not query
				esQuery = GenerateBoolMustNotQuery(processedTerm[1], value)
			} else if isContainsOperator(expr.Operator().String()) {
				// generate ES Query String query
				esQuery = GenerateQueryStringQuery(processedTerm[1], value)
			} else if isRegexpMatchOperator(expr.Operator().String()) {
				// generate ES Regexp query
				esQuery = GenerateRegexpQuery(processedTerm[1], value)
			} else {
				return Result{}, fmt.Errorf("invalid expression: operator not supported: %v", expr.Operator().String())
			}

			// fmt.Printf("OPA Query #%d: %v\n", i+1, pq.Queries[i])
			// fmt.Printf("ES  Query #%d: %+v\n\n", i+1, esQuery)
			exprQueries = append(exprQueries, esQuery)
		}

		if len(exprQueries) == 1 {
			queries = append(queries, exprQueries[0])
		} else {
			// ES queries generated within a rule are And'ed
			boolQuery := GenerateBoolFilterQuery(exprQueries)
			queries = append(queries, boolQuery)
		}
	}

	// ES queries generated from partial eval queries
	// are Or'ed
	combinedQuery := GenerateBoolShouldQuery(queries)
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

func GenerateBoolShouldQuery(queries []elastic.Query) *elastic.BoolQuery {
	q := elastic.NewBoolQuery().QueryName("BoolShouldQuery")
	for _, query := range queries {
		q = q.Should(query)
	}
	return q
}

// GenerateBoolMustNotQuery returns an ES Must Not Bool Query.
func GenerateBoolMustNotQuery(fieldName string, fieldValue interface{}) *elastic.BoolQuery {
	q := elastic.NewBoolQuery().QueryName("BoolMustNotQuery")
	q = q.MustNot(elastic.NewTermQuery(fieldName, fieldValue))
	return q
}

// GenerateMatchAllQuery returns an ES MatchAll Query.
func GenerateMatchAllQuery() *elastic.MatchAllQuery {
	return elastic.NewMatchAllQuery().QueryName("MatchAllQuery")
}

// GenerateMatchQuery returns an ES Match Query.
func GenerateMatchQuery(fieldName string, fieldValue interface{}) *elastic.MatchQuery {
	return elastic.NewMatchQuery(fieldName, fieldValue).QueryName("MatchQuery")
}

// GenerateQueryStringQuery returns an ES Query String Query.
func GenerateQueryStringQuery(fieldName string, fieldValue interface{}) *elastic.QueryStringQuery {
	queryString := fmt.Sprintf("*%s*", fieldValue)
	q := elastic.NewQueryStringQuery(queryString).QueryName("QueryStringQuery")
	q = q.DefaultField(fieldName)
	return q
}

// GenerateRegexpQuery returns an ES Regexp Query.
func GenerateRegexpQuery(fieldName string, fieldValue interface{}) *elastic.RegexpQuery {
	return elastic.NewRegexpQuery(fieldName, fieldValue.(string))
}

// GenerateRangeQueryLt returns an ES Less Than Range Query.
func GenerateRangeQueryLt(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Lt(val)
}

// GenerateRangeQueryLte returns an ES Less Than or Equal Range Query.
func GenerateRangeQueryLte(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Lte(val)
}

// GenerateRangeQueryGt returns an ES Greater Than Range Query.
func GenerateRangeQueryGt(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Gt(val)
}

// GenerateRangeQueryGte returns an ES Greater Than or Equal Range Query.
func GenerateRangeQueryGte(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Gte(val)
}

// GenerateTermQuery returns an ES Term Query.
func GenerateTermQuery(fieldName string, fieldValue interface{}) *elastic.TermQuery {
	return elastic.NewTermQuery(fieldName, fieldValue).QueryName("TermQuery")

}

// GenerateNestedQuery returns an ES Nested Query.
func GenerateNestedQuery(path string, query elastic.Query) *elastic.NestedQuery {
	return elastic.NewNestedQuery(path, query).QueryName("NestedQuery").IgnoreUnmapped(true)

}

// GenerateBoolFilterQuery returns an ES Filter Bool Query.
func GenerateBoolFilterQuery(filters []elastic.Query) *elastic.BoolQuery {
	q := elastic.NewBoolQuery()
	for _, filter := range filters {
		q = q.Filter(filter)
	}
	q = q.QueryName("BoolFilterQuery")
	return q

}
