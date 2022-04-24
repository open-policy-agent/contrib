package esquery

import "github.com/fatih/structs"

// ConstantScoreQuery represents a compound query of type "constant_score", as
// described in
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-constant-score-query.html
type ConstantScoreQuery struct {
	filter Mappable
	boost  float32
}

// ConstantScore creates a new query of type "contant_score" with the provided
// filter query.
func ConstantScore(filter Mappable) *ConstantScoreQuery {
	return &ConstantScoreQuery{
		filter: filter,
	}
}

// Boost sets the boost value of the query.
func (q *ConstantScoreQuery) Boost(b float32) *ConstantScoreQuery {
	q.boost = b
	return q
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (q *ConstantScoreQuery) Map() map[string]interface{} {
	return map[string]interface{}{
		"constant_score": structs.Map(struct {
			Filter map[string]interface{} `structs:"filter"`
			Boost  float32                `structs:"boost,omitempty"`
		}{q.filter.Map(), q.boost}),
	}
}
