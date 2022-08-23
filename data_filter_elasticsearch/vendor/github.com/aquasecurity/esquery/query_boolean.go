package esquery

import "github.com/fatih/structs"

// BoolQuery represents a compound query of type "bool", as described in
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-bool-query.html
type BoolQuery struct {
	must               []Mappable
	filter             []Mappable
	mustNot            []Mappable
	should             []Mappable
	minimumShouldMatch int16
	boost              float32
}

// Bool creates a new compound query of type "bool".
func Bool() *BoolQuery {
	return &BoolQuery{}
}

// Must adds one or more queries of type "must" to the bool query. Must can be
// called multiple times, queries will be appended to existing ones.
func (q *BoolQuery) Must(must ...Mappable) *BoolQuery {
	q.must = append(q.must, must...)
	return q
}

// Filter adds one or more queries of type "filter" to the bool query. Filter
// can be called multiple times, queries will be appended to existing ones.
func (q *BoolQuery) Filter(filter ...Mappable) *BoolQuery {
	q.filter = append(q.filter, filter...)
	return q
}

// Must adds one or more queries of type "must_not" to the bool query. MustNot
// can be called multiple times, queries will be appended to existing ones.
func (q *BoolQuery) MustNot(mustnot ...Mappable) *BoolQuery {
	q.mustNot = append(q.mustNot, mustnot...)
	return q
}

// Should adds one or more queries of type "should" to the bool query. Should can be
// called multiple times, queries will be appended to existing ones.
func (q *BoolQuery) Should(should ...Mappable) *BoolQuery {
	q.should = append(q.should, should...)
	return q
}

// MinimumShouldMatch sets the number or percentage of should clauses returned
// documents must match.
func (q *BoolQuery) MinimumShouldMatch(val int16) *BoolQuery {
	q.minimumShouldMatch = val
	return q
}

// Boost sets the boost value for the query.
func (q *BoolQuery) Boost(val float32) *BoolQuery {
	q.boost = val
	return q
}

// Map returns a map representation of the bool query, thus implementing
// the Mappable interface.
func (q *BoolQuery) Map() map[string]interface{} {
	var data struct {
		Must               []map[string]interface{} `structs:"must,omitempty"`
		Filter             []map[string]interface{} `structs:"filter,omitempty"`
		MustNot            []map[string]interface{} `structs:"must_not,omitempty"`
		Should             []map[string]interface{} `structs:"should,omitempty"`
		MinimumShouldMatch int16                    `structs:"minimum_should_match,omitempty"`
		Boost              float32                  `structs:"boost,omitempty"`
	}

	data.MinimumShouldMatch = q.minimumShouldMatch
	data.Boost = q.boost

	if len(q.must) > 0 {
		data.Must = make([]map[string]interface{}, len(q.must))
		for i, m := range q.must {
			data.Must[i] = m.Map()
		}
	}

	if len(q.filter) > 0 {
		data.Filter = make([]map[string]interface{}, len(q.filter))
		for i, m := range q.filter {
			data.Filter[i] = m.Map()
		}
	}

	if len(q.mustNot) > 0 {
		data.MustNot = make([]map[string]interface{}, len(q.mustNot))
		for i, m := range q.mustNot {
			data.MustNot[i] = m.Map()
		}
	}

	if len(q.should) > 0 {
		data.Should = make([]map[string]interface{}, len(q.should))
		for i, m := range q.should {
			data.Should[i] = m.Map()
		}
	}

	return map[string]interface{}{
		"bool": structs.Map(data),
	}
}
