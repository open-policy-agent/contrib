package esquery

import (
	"github.com/fatih/structs"
)

type matchType uint8

const (
	// TypeMatch denotes a query of type "match"
	TypeMatch matchType = iota

	// TypeMatchBool denotes a query of type "match_bool_prefix"
	TypeMatchBoolPrefix

	// TypeMatchPhrase denotes a query of type "match_phrase"
	TypeMatchPhrase

	// TypeMatchPhrasePrefix denotes a query of type "match_phrase_prefix"
	TypeMatchPhrasePrefix
)

// MatchQuery represents a query of type "match", "match_bool_prefix",
// "match_phrase" and "match_phrase_prefix". While all four share the same
// general structure, they don't necessarily support all the same options. The
// library does not attempt to verify provided options are supported.
// See the ElasticSearch documentation for more information:
// - https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-match-query.html
// - https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-match-bool-prefix-query.html
// - https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-match-query-phrase.html
// - https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-match-query-phrase-prefix.html
type MatchQuery struct {
	field  string
	mType  matchType
	params matchParams
}

// Map returns a map representation of the query, thus implementing the
// Mappable interface.
func (q *MatchQuery) Map() map[string]interface{} {
	var mType string
	switch q.mType {
	case TypeMatch:
		mType = "match"
	case TypeMatchBoolPrefix:
		mType = "match_bool_prefix"
	case TypeMatchPhrase:
		mType = "match_phrase"
	case TypeMatchPhrasePrefix:
		mType = "match_phrase_prefix"
	}

	return map[string]interface{}{
		mType: map[string]interface{}{
			q.field: structs.Map(q.params),
		},
	}
}

type matchParams struct {
	Qry          interface{}   `structs:"query"`
	Anl          string        `structs:"analyzer,omitempty"`
	AutoGenerate *bool         `structs:"auto_generate_synonyms_phrase_query,omitempty"`
	Fuzz         string        `structs:"fuzziness,omitempty"`
	MaxExp       uint16        `structs:"max_expansions,omitempty"`
	PrefLen      uint16        `structs:"prefix_length,omitempty"`
	Trans        *bool         `structs:"transpositions,omitempty"`
	FuzzyRw      string        `structs:"fuzzy_rewrite,omitempty"`
	Lent         bool          `structs:"lenient,omitempty"`
	Op           MatchOperator `structs:"operator,string,omitempty"`
	MinMatch     string        `structs:"minimum_should_match,omitempty"`
	ZeroTerms    ZeroTerms     `structs:"zero_terms_query,string,omitempty"`
	Slp          uint16        `structs:"slop,omitempty"` // only relevant for match_phrase query
}

// Match creates a new query of type "match" with the provided field name.
// A comparison value can optionally be provided to quickly create a simple
// query such as { "match": { "message": "this is a test" } }
func Match(fieldName string, simpleQuery ...interface{}) *MatchQuery {
	return newMatch(TypeMatch, fieldName, simpleQuery...)
}

// MatchBoolPrefix creates a new query of type "match_bool_prefix" with the
// provided field name. A comparison value can optionally be provided to quickly
// create a simple query such as { "match": { "message": "this is a test" } }
func MatchBoolPrefix(fieldName string, simpleQuery ...interface{}) *MatchQuery {
	return newMatch(TypeMatchBoolPrefix, fieldName, simpleQuery...)
}

// MatchPhrase creates a new query of type "match_phrase" with the
// provided field name. A comparison value can optionally be provided to quickly
// create a simple query such as { "match": { "message": "this is a test" } }
func MatchPhrase(fieldName string, simpleQuery ...interface{}) *MatchQuery {
	return newMatch(TypeMatchPhrase, fieldName, simpleQuery...)
}

// MatchPhrasePrefix creates a new query of type "match_phrase_prefix" with the
// provided field name. A comparison value can optionally be provided to quickly
// create a simple query such as { "match": { "message": "this is a test" } }
func MatchPhrasePrefix(fieldName string, simpleQuery ...interface{}) *MatchQuery {
	return newMatch(TypeMatchPhrasePrefix, fieldName, simpleQuery...)
}

func newMatch(mType matchType, fieldName string, simpleQuery ...interface{}) *MatchQuery {
	var qry interface{}
	if len(simpleQuery) > 0 {
		qry = simpleQuery[len(simpleQuery)-1]
	}

	return &MatchQuery{
		field: fieldName,
		mType: mType,
		params: matchParams{
			Qry: qry,
		},
	}
}

// Query sets the data to find in the query's field (it is the "query" component
// of the query).
func (q *MatchQuery) Query(data interface{}) *MatchQuery {
	q.params.Qry = data
	return q
}

// Analyzer sets the analyzer used to convert the text in the "query" value into
// tokens.
func (q *MatchQuery) Analyzer(a string) *MatchQuery {
	q.params.Anl = a
	return q
}

// AutoGenerateSynonymsPhraseQuery sets the "auto_generate_synonyms_phrase_query"
// boolean.
func (q *MatchQuery) AutoGenerateSynonymsPhraseQuery(b bool) *MatchQuery {
	q.params.AutoGenerate = &b
	return q
}

// Fuzziness set the maximum edit distance allowed for matching.
func (q *MatchQuery) Fuzziness(f string) *MatchQuery {
	q.params.Fuzz = f
	return q
}

// MaxExpansions sets the maximum number of terms to which the query will expand.
func (q *MatchQuery) MaxExpansions(e uint16) *MatchQuery {
	q.params.MaxExp = e
	return q
}

// PrefixLength sets the number of beginning characters left unchanged for fuzzy
// matching.
func (q *MatchQuery) PrefixLength(l uint16) *MatchQuery {
	q.params.PrefLen = l
	return q
}

// Transpositions sets whether edits for fuzzy matching include transpositions
// of two adjacent characters.
func (q *MatchQuery) Transpositions(b bool) *MatchQuery {
	q.params.Trans = &b
	return q
}

// FuzzyRewrite sets the method used to rewrite the query.
func (q *MatchQuery) FuzzyRewrite(s string) *MatchQuery {
	q.params.FuzzyRw = s
	return q
}

// Lenient sets whether format-based errors should be ignored.
func (q *MatchQuery) Lenient(b bool) *MatchQuery {
	q.params.Lent = b
	return q
}

// Operator sets the boolean logic used to interpret text in the query value.
func (q *MatchQuery) Operator(op MatchOperator) *MatchQuery {
	q.params.Op = op
	return q
}

// MinimumShouldMatch sets the minimum number of clauses that must match for a
// document to be returned.
func (q *MatchQuery) MinimumShouldMatch(s string) *MatchQuery {
	q.params.MinMatch = s
	return q
}

// Slop sets the maximum number of positions allowed between matching tokens.
func (q *MatchQuery) Slop(n uint16) *MatchQuery {
	q.params.Slp = n
	return q
}

// ZeroTermsQuery sets the "zero_terms_query" option to use. This indicates
// whether no documents are returned if the analyzer removes all tokens, such as
// when using a stop filter.
func (q *MatchQuery) ZeroTermsQuery(s ZeroTerms) *MatchQuery {
	q.params.ZeroTerms = s
	return q
}

// MatchOperator is an enumeration type representing supported values for a
// match query's "operator" parameter.
type MatchOperator uint8

const (
	// OperatorOr is the "or" operator
	OperatorOr MatchOperator = iota

	// OperatorAnd is the "and" operator
	OperatorAnd
)

// String returns a string representation of the match operator, as known to
// ElasticSearch.
func (a MatchOperator) String() string {
	switch a {
	case OperatorOr:
		return "OR"
	case OperatorAnd:
		return "AND"
	default:
		return ""
	}
}

// ZeroTerms is an enumeration type representing supported values for a match
// query's "zero_terms_query" parameter.
type ZeroTerms uint8

const (
	// ZeroTermsNone is the "none" value
	ZeroTermsNone ZeroTerms = iota

	// ZeroTermsAll is the "all" value
	ZeroTermsAll
)

// String returns a string representation of the zero_terms_query parameter, as
// known to ElasticSearch.
func (a ZeroTerms) String() string {
	switch a {
	case ZeroTermsNone:
		return "none"
	case ZeroTermsAll:
		return "all"
	default:
		return ""
	}
}
