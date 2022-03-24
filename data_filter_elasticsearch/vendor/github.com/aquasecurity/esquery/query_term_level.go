package esquery

import (
	"github.com/fatih/structs"
)

// ExistsQuery represents a query of type "exists", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-exists-query.html
type ExistsQuery struct {
	// Field is the name of the field to check for existence
	Field string `structs:"field"`
}

// Exists creates a new query of type "exists" on the provided field.
func Exists(field string) *ExistsQuery {
	return &ExistsQuery{field}
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (q *ExistsQuery) Map() map[string]interface{} {
	return map[string]interface{}{
		"exists": structs.Map(q),
	}
}

//----------------------------------------------------------------------------//

// IDsQuery represents a query of type "ids", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-ids-query.html
type IDsQuery struct {
	// IDs is the "ids" component of the query
	IDs struct {
		// Values is the list of ID values
		Values []string `structs:"values"`
	} `structs:"ids"`
}

// IDs creates a new query of type "ids" with the provided values.
func IDs(vals ...string) *IDsQuery {
	q := &IDsQuery{}
	q.IDs.Values = vals
	return q
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (q *IDsQuery) Map() map[string]interface{} {
	return structs.Map(q)
}

//----------------------------------------------------------------------------//

// PrefixQuery represents query of type "prefix", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-prefix-query.html
type PrefixQuery struct {
	field  string
	params prefixQueryParams
}

type prefixQueryParams struct {
	// Value is the prefix value to look for
	Value string `structs:"value"`

	// Rewrite is the method used to rewrite the query
	Rewrite string `structs:"rewrite,omitempty"`
}

// Prefix creates a new query of type "prefix", on the provided field and using
// the provided prefix value.
func Prefix(field, value string) *PrefixQuery {
	return &PrefixQuery{
		field:  field,
		params: prefixQueryParams{Value: value},
	}
}

// Rewrite sets the rewrite method for the query
func (q *PrefixQuery) Rewrite(s string) *PrefixQuery {
	q.params.Rewrite = s
	return q
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (q *PrefixQuery) Map() map[string]interface{} {
	return map[string]interface{}{
		"prefix": map[string]interface{}{
			q.field: structs.Map(q.params),
		},
	}
}

//----------------------------------------------------------------------------//

// RangeQuery represents a query of type "range", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-range-query.html
type RangeQuery struct {
	field  string
	params rangeQueryParams
}

type rangeQueryParams struct {
	Gt       interface{}   `structs:"gt,omitempty"`
	Gte      interface{}   `structs:"gte,omitempty"`
	Lt       interface{}   `structs:"lt,omitempty"`
	Lte      interface{}   `structs:"lte,omitempty"`
	Format   string        `structs:"format,omitempty"`
	Relation RangeRelation `structs:"relation,string,omitempty"`
	TimeZone string        `structs:"time_zone,omitempty"`
	Boost    float32       `structs:"boost,omitempty"`
}

// Range creates a new query of type "range" on the provided field
func Range(field string) *RangeQuery {
	return &RangeQuery{field: field}
}

// Gt sets that the value of field must be greater than the provided value
func (a *RangeQuery) Gt(val interface{}) *RangeQuery {
	a.params.Gt = val
	return a
}

// Gt sets that the value of field must be greater than or equal to the provided
// value
func (a *RangeQuery) Gte(val interface{}) *RangeQuery {
	a.params.Gte = val
	return a
}

// Lt sets that the value of field must be lower than the provided value
func (a *RangeQuery) Lt(val interface{}) *RangeQuery {
	a.params.Lt = val
	return a
}

// Lte sets that the value of field must be lower than or equal to the provided
// value
func (a *RangeQuery) Lte(val interface{}) *RangeQuery {
	a.params.Lte = val
	return a
}

// Format sets the date format for date values
func (a *RangeQuery) Format(f string) *RangeQuery {
	a.params.Format = f
	return a
}

// Relation sets how the query matches values for range fields
func (a *RangeQuery) Relation(r RangeRelation) *RangeQuery {
	a.params.Relation = r
	return a
}

// TimeZone sets the time zone used for date values.
func (a *RangeQuery) TimeZone(zone string) *RangeQuery {
	a.params.TimeZone = zone
	return a
}

// Boost sets the boost value of the query.
func (a *RangeQuery) Boost(b float32) *RangeQuery {
	a.params.Boost = b
	return a
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (a *RangeQuery) Map() map[string]interface{} {
	return map[string]interface{}{
		"range": map[string]interface{}{
			a.field: structs.Map(a.params),
		},
	}
}

// RangeRelation is an enumeration type for a range query's "relation" field
type RangeRelation uint8

const (
	_ RangeRelation = iota

	// RangeIntersects is the "INTERSECTS" relation
	RangeIntersects

	// RangeContains is the "CONTAINS" relation
	RangeContains

	// RangeWithin is the "WITHIN" relation
	RangeWithin
)

// String returns a string representation of the RangeRelation value, as
// accepted by ElasticSearch
func (a RangeRelation) String() string {
	switch a {
	case RangeIntersects:
		return "INTERSECTS"
	case RangeContains:
		return "CONTAINS"
	case RangeWithin:
		return "WITHIN"
	default:
		return ""
	}
}

//----------------------------------------------------------------------------//

// RegexpQuery represents a query of type "regexp", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-regexp-query.html
type RegexpQuery struct {
	field    string
	wildcard bool
	params   regexpQueryParams
}

type regexpQueryParams struct {
	Value                 string `structs:"value"`
	Flags                 string `structs:"flags,omitempty"`
	MaxDeterminizedStates uint16 `structs:"max_determinized_states,omitempty"`
	Rewrite               string `structs:"rewrite,omitempty"`
}

// Regexp creates a new query of type "regexp" on the provided field and using
// the provided regular expression.
func Regexp(field, value string) *RegexpQuery {
	return &RegexpQuery{
		field: field,
		params: regexpQueryParams{
			Value: value,
		},
	}
}

// Value changes the regular expression value of the query.
func (q *RegexpQuery) Value(v string) *RegexpQuery {
	q.params.Value = v
	return q
}

// Flags sets the regular expression's optional flags.
func (q *RegexpQuery) Flags(f string) *RegexpQuery {
	if !q.wildcard {
		q.params.Flags = f
	}
	return q
}

// MaxDeterminizedStates sets the maximum number of automaton states required
// for the query.
func (q *RegexpQuery) MaxDeterminizedStates(m uint16) *RegexpQuery {
	if !q.wildcard {
		q.params.MaxDeterminizedStates = m
	}
	return q
}

// Rewrite sets the method used to rewrite the query.
func (q *RegexpQuery) Rewrite(r string) *RegexpQuery {
	q.params.Rewrite = r
	return q
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (q *RegexpQuery) Map() map[string]interface{} {
	var qType string
	if q.wildcard {
		qType = "wildcard"
	} else {
		qType = "regexp"
	}
	return map[string]interface{}{
		qType: map[string]interface{}{
			q.field: structs.Map(q.params),
		},
	}
}

//----------------------------------------------------------------------------//

// Wildcard creates a new query of type "wildcard" on the provided field and
// using the provided regular expression value. Internally, wildcard queries
// are simply specialized RegexpQuery values.
// Wildcard queries are described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-wildcard-query.html
func Wildcard(field, value string) *RegexpQuery {
	return &RegexpQuery{
		field:    field,
		wildcard: true,
		params: regexpQueryParams{
			Value: value,
		},
	}
}

//----------------------------------------------------------------------------//

// FuzzyQuery represents a query of type "fuzzy", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-fuzzy-query.html
type FuzzyQuery struct {
	field  string
	params fuzzyQueryParams
}

type fuzzyQueryParams struct {
	Value          string `structs:"value"`
	Fuzziness      string `structs:"fuzziness,omitempty"`
	MaxExpansions  uint16 `structs:"max_expansions,omitempty"`
	PrefixLength   uint16 `structs:"prefix_length,omitempty"`
	Transpositions *bool  `structs:"transpositions,omitempty"`
	Rewrite        string `structs:"rewrite,omitempty"`
}

// Fuzzy creates a new query of type "fuzzy" on the provided field and using
// the provided value
func Fuzzy(field, value string) *FuzzyQuery {
	return &FuzzyQuery{
		field: field,
		params: fuzzyQueryParams{
			Value: value,
		},
	}
}

// Value sets the value of the query.
func (q *FuzzyQuery) Value(val string) *FuzzyQuery {
	q.params.Value = val
	return q
}

// Fuzziness sets the maximum edit distance allowed for matching.
func (q *FuzzyQuery) Fuzziness(fuzz string) *FuzzyQuery {
	q.params.Fuzziness = fuzz
	return q
}

// MaxExpansions sets the maximum number of variations created.
func (q *FuzzyQuery) MaxExpansions(m uint16) *FuzzyQuery {
	q.params.MaxExpansions = m
	return q
}

// PrefixLength sets the number of beginning characters left unchanged when
// creating expansions
func (q *FuzzyQuery) PrefixLength(l uint16) *FuzzyQuery {
	q.params.PrefixLength = l
	return q
}

// Transpositions sets whether edits include transpositions of two adjacent
// characters.
func (q *FuzzyQuery) Transpositions(b bool) *FuzzyQuery {
	q.params.Transpositions = &b
	return q
}

// Rewrite sets the method used to rewrite the query.
func (q *FuzzyQuery) Rewrite(s string) *FuzzyQuery {
	q.params.Rewrite = s
	return q
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (q *FuzzyQuery) Map() map[string]interface{} {
	return map[string]interface{}{
		"fuzzy": map[string]interface{}{
			q.field: structs.Map(q.params),
		},
	}
}

//----------------------------------------------------------------------------//

// TermQuery represents a query of type "term", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-term-query.html
type TermQuery struct {
	field  string
	params termQueryParams
}

type termQueryParams struct {
	Value interface{} `structs:"value"`
	Boost float32     `structs:"boost,omitempty"`
}

// Term creates a new query of type "term" on the provided field and using the
// provide value
func Term(field string, value interface{}) *TermQuery {
	return &TermQuery{
		field: field,
		params: termQueryParams{
			Value: value,
		},
	}
}

// Value sets the term value for the query.
func (q *TermQuery) Value(val interface{}) *TermQuery {
	q.params.Value = val
	return q
}

// Boost sets the boost value of the query.
func (q *TermQuery) Boost(b float32) *TermQuery {
	q.params.Boost = b
	return q
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (q *TermQuery) Map() map[string]interface{} {
	return map[string]interface{}{
		"term": map[string]interface{}{
			q.field: structs.Map(q.params),
		},
	}
}

//----------------------------------------------------------------------------//

// TermsQuery represents a query of type "terms", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-terms-query.html
type TermsQuery struct {
	field  string
	values []interface{}
	boost  float32
}

// Terms creates a new query of type "terms" on the provided field, and
// optionally with the provided term values.
func Terms(field string, values ...interface{}) *TermsQuery {
	return &TermsQuery{
		field:  field,
		values: values,
	}
}

// Values sets the term values for the query.
func (q *TermsQuery) Values(values ...interface{}) *TermsQuery {
	q.values = values
	return q
}

// Boost sets the boost value of the query.
func (q *TermsQuery) Boost(b float32) *TermsQuery {
	q.boost = b
	return q
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (q TermsQuery) Map() map[string]interface{} {
	innerMap := map[string]interface{}{q.field: q.values}
	if q.boost > 0 {
		innerMap["boost"] = q.boost
	}

	return map[string]interface{}{"terms": innerMap}
}

//----------------------------------------------------------------------------//

// TermsSetQuery represents a query of type "terms_set", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-terms-set-query.html
type TermsSetQuery struct {
	field  string
	params termsSetQueryParams
}

type termsSetQueryParams struct {
	Terms                    []string `structs:"terms"`
	MinimumShouldMatchField  string   `structs:"minimum_should_match_field,omitempty"`
	MinimumShouldMatchScript string   `structs:"minimum_should_match_script,omitempty"`
}

// TermsSet creates a new query of type "terms_set" on the provided field and
// optionally using the provided terms.
func TermsSet(field string, terms ...string) *TermsSetQuery {
	return &TermsSetQuery{
		field: field,
		params: termsSetQueryParams{
			Terms: terms,
		},
	}
}

// Terms sets the terms for the query.
func (q *TermsSetQuery) Terms(terms ...string) *TermsSetQuery {
	q.params.Terms = terms
	return q
}

// MinimumShouldMatchField sets the name of the field containing the number of
// matching terms required to return a document.
func (q *TermsSetQuery) MinimumShouldMatchField(field string) *TermsSetQuery {
	q.params.MinimumShouldMatchField = field
	return q
}

// MinimumShouldMatchScript sets the custom script containing the number of
// matching terms required to return a document.
func (q *TermsSetQuery) MinimumShouldMatchScript(script string) *TermsSetQuery {
	q.params.MinimumShouldMatchScript = script
	return q
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (q TermsSetQuery) Map() map[string]interface{} {
	return map[string]interface{}{
		"terms_set": map[string]interface{}{
			q.field: structs.Map(q.params),
		},
	}
}
