package esquery

import (
	"github.com/fatih/structs"
)

// Map returns a map representation of the highlight; implementing the
// Mappable interface.
func (q *QueryHighlight) Map() map[string]interface{} {
	results := structs.Map(q.params)
	if q.highlightQuery != nil {
		results["query"] = q.highlightQuery.Map()
	}
	if q.fields != nil && len(q.fields) > 0 {
		fields := make(map[string]interface{})
		for k, v := range q.fields {
			fields[k] = v.Map()
		}
		results["fields"] = fields
	}
	return results
}

type QueryHighlight struct {
	highlightQuery Mappable                   `structs:"highlight_query,omitempty"`
	fields         map[string]*QueryHighlight `structs:"fields"`
	params         highlighParams
}

type highlighParams struct {
	PreTags  []string `structs:"pre_tags,omitempty"`
	PostTags []string `structs:"post_tags,omitempty"`

	FragmentSize          uint16                   `structs:"fragment_size,omitempty"`
	NumberOfFragments     uint16                   `structs:"number_of_fragments,omitempty"`
	Type                  HighlightType            `structs:"type,string,omitempty"`
	BoundaryChars         string                   `structs:"boundary_chars,omitempty"`
	BoundaryMaxScan       uint16                   `structs:"boundary_max_scan,omitempty"`
	BoundaryScanner       HighlightBoundaryScanner `structs:"boundary_scanner,string,omitempty"`
	BoundaryScannerLocale string                   `structs:"boundary_scanner_locale,omitempty"`
	Encoder               HighlightEncoder         `structs:"encoder,string,omitempty"`
	ForceSource           *bool                    `structs:"force_source,omitempty"`
	Fragmenter            HighlightFragmenter      `structs:"fragmenter,string,omitempty"`
	FragmentOffset        uint16                   `structs:"fragment_offset,omitempty"`
	MatchedFields         []string                 `structs:"matched_fields,omitempty"`
	NoMatchSize           uint16                   `structs:"no_match_size,omitempty"`
	Order                 HighlightOrder           `structs:"order,string,omitempty"`
	PhraseLimit           uint16                   `structs:"phrase_limit,omitempty"`
	RequireFieldMatch     *bool                    `structs:"require_field_match,omitempty"`
	TagsSchema            HighlightTagsSchema      `structs:"tags_schema,string,omitempty"`
}

// Highlight creates a new "query" of type "highlight"
func Highlight() *QueryHighlight {
	return newHighlight()
}

func newHighlight() *QueryHighlight {
	return &QueryHighlight{
		fields: make(map[string]*QueryHighlight),
		params: highlighParams{},
	}
}

// PreTags sets the highlight query's pre_tags ignore unmapped field
func (q *QueryHighlight) PreTags(s ...string) *QueryHighlight {
	q.params.PreTags = append(q.params.PreTags, s...)
	return q
}

// PostTags sets the highlight query's post_tags ignore unmapped field
func (q *QueryHighlight) PostTags(s ...string) *QueryHighlight {
	q.params.PostTags = append(q.params.PostTags, s...)
	return q
}

// Field sets an entry the highlight query's fields
func (q *QueryHighlight) Field(name string, h ...*QueryHighlight) *QueryHighlight {
	var fld *QueryHighlight
	if len(h) > 0 {
		fld = h[len(h)-1]
	} else {
		fld = &QueryHighlight{}
	}
	q.fields[name] = fld
	return q
}

// Fields sets all entries for the highlight query's fields
func (q *QueryHighlight) Fields(h map[string]*QueryHighlight) *QueryHighlight {
	q.fields = h
	return q
}

// FragmentSize sets the highlight query's fragment_size ignore unmapped field
func (q *QueryHighlight) FragmentSize(i uint16) *QueryHighlight {
	q.params.FragmentSize = i
	return q
}

// NumberOfFragments sets the highlight query's number_of_fragments ignore unmapped field
func (q *QueryHighlight) NumberOfFragments(i uint16) *QueryHighlight {
	q.params.NumberOfFragments = i
	return q
}

// Type sets the highlight query's type ignore unmapped field
func (q *QueryHighlight) Type(t HighlightType) *QueryHighlight {
	q.params.Type = t
	return q
}

// BoundaryChars sets the highlight query's boundary_chars ignore unmapped field
func (q *QueryHighlight) BoundaryChars(s string) *QueryHighlight {
	q.params.BoundaryChars = s
	return q
}

// BoundaryMaxScan sets the highlight query's boundary_max_scan ignore unmapped field
func (q *QueryHighlight) BoundaryMaxScan(i uint16) *QueryHighlight {
	q.params.BoundaryMaxScan = i
	return q
}

// BoundaryScanner sets the highlight query's boundary_scanner ignore unmapped field
func (q *QueryHighlight) BoundaryScanner(t HighlightBoundaryScanner) *QueryHighlight {
	q.params.BoundaryScanner = t
	return q
}

// BoundaryScannerLocale sets the highlight query's boundary_scanner_locale ignore unmapped field
func (q *QueryHighlight) BoundaryScannerLocale(l string) *QueryHighlight {
	q.params.BoundaryScannerLocale = l
	return q
}

// Encoder sets the highlight query's encoder ignore unmapped field
func (q *QueryHighlight) Encoder(e HighlightEncoder) *QueryHighlight {
	q.params.Encoder = e
	return q
}

// ForceSource sets the highlight query's force_source ignore unmapped field
func (q *QueryHighlight) ForceSource(b bool) *QueryHighlight {
	q.params.ForceSource = &b
	return q
}

// Fragmenter sets the highlight query's fragmenter ignore unmapped field
func (q *QueryHighlight) Fragmenter(f HighlightFragmenter) *QueryHighlight {
	q.params.Fragmenter = f
	return q
}

// FragmentOffset sets the highlight query's fragment_offset ignore unmapped field
func (q *QueryHighlight) FragmentOffset(i uint16) *QueryHighlight {
	q.params.FragmentOffset = i
	return q
}

// HighlightQuery sets the highlight query's highlight_query ignore unmapped field
func (q *QueryHighlight) HighlightQuery(b Mappable) *QueryHighlight {
	q.highlightQuery = b
	return q
}

// MatchedFields sets the highlight query's matched_fields ignore unmapped field
func (q *QueryHighlight) MatchedFields(s ...string) *QueryHighlight {
	q.params.MatchedFields = append(q.params.MatchedFields, s...)
	return q
}

// NoMatchSize sets the highlight query's no_match_size ignore unmapped field
func (q *QueryHighlight) NoMatchSize(i uint16) *QueryHighlight {
	q.params.NoMatchSize = i
	return q
}

// Order sets the nested highlight's score order unmapped field
func (q *QueryHighlight) Order(o HighlightOrder) *QueryHighlight {
	q.params.Order = o
	return q
}

// PhraseLimit sets the highlight query's phrase_limit ignore unmapped field
func (q *QueryHighlight) PhraseLimit(i uint16) *QueryHighlight {
	q.params.PhraseLimit = i
	return q
}

// RequireFieldMatch sets the highlight query's require_field_match ignore unmapped field
func (q *QueryHighlight) RequireFieldMatch(b bool) *QueryHighlight {
	q.params.RequireFieldMatch = &b
	return q
}

// TagsSchema sets the highlight query's tags_schema ignore unmapped field
func (q *QueryHighlight) TagsSchema(s HighlightTagsSchema) *QueryHighlight {
	q.params.TagsSchema = s
	return q
}

type HighlightType uint8

const (
	// HighlighterUnified is the "unified" value
	HighlighterUnified HighlightType = iota

	// HighlighterPlain is the "plain" value
	HighlighterPlain

	// HighlighterFvh is the "fvh" value
	HighlighterFvh
)

// String returns a string representation of the type parameter, as
// known to ElasticSearch.
func (a HighlightType) String() string {
	switch a {
	case HighlighterUnified:
		return "unified"
	case HighlighterPlain:
		return "plain"
	case HighlighterFvh:
		return "fvh"
	}
	return ""
}

type HighlightBoundaryScanner uint8

const (
	BoundaryScannerDefault HighlightBoundaryScanner = iota

	// BoundaryScannerChars is the "chars" value
	BoundaryScannerChars

	// BoundaryScannerSentence is the "sentence" value
	BoundaryScannerSentence

	// BoundaryScannerWord is the "word" value
	BoundaryScannerWord
)

// String returns a string representation of the boundary_scanner parameter, as
// known to ElasticSearch.
func (a HighlightBoundaryScanner) String() string {
	switch a {
	case BoundaryScannerChars:
		return "chars"
	case BoundaryScannerSentence:
		return "sentence"
	case BoundaryScannerWord:
		return "word"
	}
	return ""
}

type HighlightEncoder uint8

const (
	// EncoderDefault is the "default" value
	EncoderDefault HighlightEncoder = iota

	// EncoderHtml is the "html" value
	EncoderHtml
)

// String returns a string representation of the encoder parameter, as
// known to ElasticSearch.
func (a HighlightEncoder) String() string {
	switch a {
	case EncoderDefault:
		return "default"
	case EncoderHtml:
		return "html"
	}
	return ""
}

type HighlightFragmenter uint8

const (
	// FragmentSpan is the "span" value
	FragmenterSpan HighlightFragmenter = iota

	// FragmenterSimple is the "simple" value
	FragmenterSimple
)

// String returns a string representation of the fragmenter parameter, as
// known to ElasticSearch.
func (a HighlightFragmenter) String() string {
	switch a {
	case FragmenterSpan:
		return "span"
	case FragmenterSimple:
		return "simple"
	}
	return ""
}

type HighlightOrder uint8

const (
	// OrderNone is the "none" value
	OrderNone HighlightOrder = iota

	// OrderScore is the "score" value
	OrderScore
)

// String returns a string representation of the order parameter, as
// known to ElasticSearch.
func (a HighlightOrder) String() string {
	switch a {
	case OrderNone:
		return "none"
	case OrderScore:
		return "score"
	}
	return ""
}

type HighlightTagsSchema uint8

const (
	TagsSchemaDefault HighlightTagsSchema = iota
	// TagsSchemaStyled is the "styled" value
	TagsSchemaStyled
)

// String returns a string representation of the tags_schema parameter, as
// known to ElasticSearch.
func (a HighlightTagsSchema) String() string {
	switch a {
	case TagsSchemaStyled:
		return "styled"
	}
	return ""
}
