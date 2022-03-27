package esquery

import "github.com/fatih/structs"

// MatchAllQuery represents a query of type "match_all" or "match_none", as
// described in
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-match-all-query.html
type MatchAllQuery struct {
	all    bool
	params matchAllParams
}

type matchAllParams struct {
	Boost float32 `structs:"boost,omitempty"`
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (q *MatchAllQuery) Map() map[string]interface{} {
	var mType string
	switch q.all {
	case true:
		mType = "match_all"
	default:
		mType = "match_none"
	}

	return map[string]interface{}{
		mType: structs.Map(q.params),
	}
}

// MatchAll creates a new query of type "match_all".
func MatchAll() *MatchAllQuery {
	return &MatchAllQuery{all: true}
}

// Boost assigns a score boost for documents matching the query.
func (q *MatchAllQuery) Boost(b float32) *MatchAllQuery {
	if q.all {
		q.params.Boost = b
	}
	return q
}

// MatchNone creates a new query of type "match_none".
func MatchNone() *MatchAllQuery {
	return &MatchAllQuery{all: false}
}
