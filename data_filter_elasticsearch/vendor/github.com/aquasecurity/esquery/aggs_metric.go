package esquery

import "github.com/fatih/structs"

// BaseAgg contains several fields that are common for all aggregation types.
type BaseAgg struct {
	name           string
	apiName        string
	*BaseAggParams `structs:",flatten"`
}

// BaseAggParams contains fields that are common to most metric-aggregation
// types.
type BaseAggParams struct {
	// Field is the name of the field to aggregate on.
	Field string `structs:"field"`
	// Miss is a value to provide for documents that are missing a value for the
	// field.
	Miss interface{} `structs:"missing,omitempty"`
}

func newBaseAgg(apiName, name, field string) *BaseAgg {
	return &BaseAgg{
		name:    name,
		apiName: apiName,
		BaseAggParams: &BaseAggParams{
			Field: field,
		},
	}
}

// Name returns the name of the aggregation, allowing implementation of the
// Aggregation interface.
func (agg *BaseAgg) Name() string {
	return agg.name
}

// Map returns a map representation of the aggregation, implementing the
// Mappable interface.
func (agg *BaseAgg) Map() map[string]interface{} {
	return map[string]interface{}{
		agg.apiName: structs.Map(agg.BaseAggParams),
	}
}

// AvgAgg represents an aggregation of type "avg", as described in
// https://www.elastic.co/guide/en/elasticsearch/reference/
//     current/search-aggregations-metrics-avg-aggregation.html
type AvgAgg struct {
	*BaseAgg `structs:",flatten"`
}

// Avg creates an aggregation of type "avg", with the provided name and on the
// provided field.
func Avg(name, field string) *AvgAgg {
	return &AvgAgg{
		BaseAgg: newBaseAgg("avg", name, field),
	}
}

// Missing sets the value to provide for documents missing a value for the
// selected field.
func (agg *AvgAgg) Missing(val interface{}) *AvgAgg {
	agg.Miss = val
	return agg
}

//----------------------------------------------------------------------------//

// WeightedAvgAgg represents an aggregation of type "weighted_avg", as described
// in https://www.elastic.co/guide/en/elasticsearch/reference/
//     current/search-aggregations-metrics-weight-avg-aggregation.html
type WeightedAvgAgg struct {
	name    string
	apiName string

	// Val is the value component of the aggregation
	Val *BaseAggParams `structs:"value"`

	// Weig is the weight component of the aggregation
	Weig *BaseAggParams `structs:"weight"`
}

// WeightedAvg creates a new aggregation of type "weighted_agg" with the
// provided name.
func WeightedAvg(name string) *WeightedAvgAgg {
	return &WeightedAvgAgg{
		name:    name,
		apiName: "weighted_avg",
	}
}

// Name returns the name of the aggregation.
func (agg *WeightedAvgAgg) Name() string {
	return agg.name
}

// Value sets the value field and optionally a value to use when records are
// missing a value for the field.
func (agg *WeightedAvgAgg) Value(field string, missing ...interface{}) *WeightedAvgAgg {
	agg.Val = new(BaseAggParams)
	agg.Val.Field = field
	if len(missing) > 0 {
		agg.Val.Miss = missing[len(missing)-1]
	}
	return agg
}

// Value sets the weight field and optionally a value to use when records are
// missing a value for the field.
func (agg *WeightedAvgAgg) Weight(field string, missing ...interface{}) *WeightedAvgAgg {
	agg.Weig = new(BaseAggParams)
	agg.Weig.Field = field
	if len(missing) > 0 {
		agg.Weig.Miss = missing[len(missing)-1]
	}
	return agg
}

// Map returns a map representation of the aggregation, thus implementing the
// Mappable interface.
func (agg *WeightedAvgAgg) Map() map[string]interface{} {
	return map[string]interface{}{
		agg.apiName: structs.Map(agg),
	}
}

//----------------------------------------------------------------------------//

// CardinalityAgg represents an aggregation of type "cardinality", as described
// in https://www.elastic.co/guide/en/elasticsearch/reference/
//     current/search-aggregations-metrics-cardinality-aggregation.html
type CardinalityAgg struct {
	*BaseAgg `structs:",flatten"`

	// PrecisionThr is the precision threshold of the aggregation
	PrecisionThr uint16 `structs:"precision_threshold,omitempty"`
}

// Cardinality creates a new aggregation of type "cardinality" with the provided
// name and on the provided field.
func Cardinality(name, field string) *CardinalityAgg {
	return &CardinalityAgg{
		BaseAgg: newBaseAgg("cardinality", name, field),
	}
}

// Missing sets the value to provide for records that are missing a value for
// the field.
func (agg *CardinalityAgg) Missing(val interface{}) *CardinalityAgg {
	agg.Miss = val
	return agg
}

// PrecisionThreshold sets the precision threshold of the aggregation.
func (agg *CardinalityAgg) PrecisionThreshold(val uint16) *CardinalityAgg {
	agg.PrecisionThr = val
	return agg
}

// Map returns a map representation of the aggregation, thus implementing the
// Mappable interface
func (agg *CardinalityAgg) Map() map[string]interface{} {
	return map[string]interface{}{
		agg.apiName: structs.Map(agg),
	}
}

//----------------------------------------------------------------------------//

// MaxAgg represents an aggregation of type "max", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/
//     current/search-aggregations-metrics-max-aggregation.html
type MaxAgg struct {
	*BaseAgg `structs:",flatten"`
}

// Max creates a new aggregation of type "max", with the provided name and on
// the provided field.
func Max(name, field string) *MaxAgg {
	return &MaxAgg{
		BaseAgg: newBaseAgg("max", name, field),
	}
}

// Missing sets the value to provide for records that are missing a value for
// the field.
func (agg *MaxAgg) Missing(val interface{}) *MaxAgg {
	agg.Miss = val
	return agg
}

//----------------------------------------------------------------------------//

// MinAgg represents an aggregation of type "min", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/
//     current/search-aggregations-metrics-min-aggregation.html
type MinAgg struct {
	*BaseAgg `structs:",flatten"`
}

// Min creates a new aggregation of type "min", with the provided name and on
// the provided field.
func Min(name, field string) *MinAgg {
	return &MinAgg{
		BaseAgg: newBaseAgg("min", name, field),
	}
}

// Missing sets the value to provide for records that are missing a value for
// the field.
func (agg *MinAgg) Missing(val interface{}) *MinAgg {
	agg.Miss = val
	return agg
}

//----------------------------------------------------------------------------//

// SumAgg represents an aggregation of type "sum", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/
//     current/search-aggregations-metrics-sum-aggregation.html
type SumAgg struct {
	*BaseAgg `structs:",flatten"`
}

// Sum creates a new aggregation of type "sum", with the provided name and on
// the provided field.
func Sum(name, field string) *SumAgg {
	return &SumAgg{
		BaseAgg: newBaseAgg("sum", name, field),
	}
}

// Missing sets the value to provide for records that are missing a value for
// the field.
func (agg *SumAgg) Missing(val interface{}) *SumAgg {
	agg.Miss = val
	return agg
}

//----------------------------------------------------------------------------//

// ValueCountAgg represents an aggregation of type "value_count", as described
// in https://www.elastic.co/guide/en/elasticsearch/reference/
//     current/search-aggregations-metrics-valuecount-aggregation.html
type ValueCountAgg struct {
	*BaseAgg `structs:",flatten"`
}

// ValueCount creates a new aggregation of type "value_count", with the provided
// name and on the provided field
func ValueCount(name, field string) *ValueCountAgg {
	return &ValueCountAgg{
		BaseAgg: newBaseAgg("value_count", name, field),
	}
}

//----------------------------------------------------------------------------//

// PercentilesAgg represents an aggregation of type "percentiles", as described
// in https://www.elastic.co/guide/en/elasticsearch/reference/
//     current/search-aggregations-metrics-percentile-aggregation.html
type PercentilesAgg struct {
	*BaseAgg `structs:",flatten"`

	// Prcnts is the aggregation's percentages
	Prcnts []float32 `structs:"percents,omitempty"`

	// Key denotes whether the aggregation is keyed or not
	Key *bool `structs:"keyed,omitempty"`

	// TDigest includes options for the TDigest algorithm
	TDigest struct {
		// Compression is the compression level to use
		Compression uint16 `structs:"compression,omitempty"`
	} `structs:"tdigest,omitempty"`

	// HDR includes options for the HDR implementation
	HDR struct {
		// NumHistogramDigits defines the resolution of values for the histogram
		// in number of significant digits
		NumHistogramDigits uint8 `structs:"number_of_significant_value_digits,omitempty"`
	} `structs:"hdr,omitempty"`
}

// Percentiles creates a new aggregation of type "percentiles" with the provided
// name and on the provided field.
func Percentiles(name, field string) *PercentilesAgg {
	return &PercentilesAgg{
		BaseAgg: newBaseAgg("percentiles", name, field),
	}
}

// Percents sets the aggregation's percentages
func (agg *PercentilesAgg) Percents(percents ...float32) *PercentilesAgg {
	agg.Prcnts = percents
	return agg
}

// Missing sets the value to provide for records that are missing a value for
// the field.
func (agg *PercentilesAgg) Missing(val interface{}) *PercentilesAgg {
	agg.Miss = val
	return agg
}

// Keyed sets whether the aggregate is keyed or not.
func (agg *PercentilesAgg) Keyed(b bool) *PercentilesAgg {
	agg.Key = &b
	return agg
}

// Compression sets the compression level for the aggregation.
func (agg *PercentilesAgg) Compression(val uint16) *PercentilesAgg {
	agg.TDigest.Compression = val
	return agg
}

// NumHistogramDigits specifies the resolution of values for the histogram in
// number of significant digits.
func (agg *PercentilesAgg) NumHistogramDigits(val uint8) *PercentilesAgg {
	agg.HDR.NumHistogramDigits = val
	return agg
}

// Map returns a map representation of the aggregation, thus implementing the
// Mappable interface.
func (agg *PercentilesAgg) Map() map[string]interface{} {
	return map[string]interface{}{
		agg.apiName: structs.Map(agg),
	}
}

//----------------------------------------------------------------------------//

// StatsAgg represents an aggregation of type "stats", as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/
//     current/search-aggregations-metrics-stats-aggregation.html
type StatsAgg struct {
	*BaseAgg `structs:",flatten"`
}

// Stats creates a new "stats" aggregation with the provided name and on the
// provided field.
func Stats(name, field string) *StatsAgg {
	return &StatsAgg{
		BaseAgg: newBaseAgg("stats", name, field),
	}
}

// Missing sets the value to provide for records missing a value for the field.
func (agg *StatsAgg) Missing(val interface{}) *StatsAgg {
	agg.Miss = val
	return agg
}

// ---------------------------------------------------------------------------//

// StringStatsAgg represents an aggregation of type "string_stats", as described
// in https://www.elastic.co/guide/en/elasticsearch/reference/
//     current/search-aggregations-metrics-string-stats-aggregation.html
type StringStatsAgg struct {
	*BaseAgg `structs:",flatten"`

	// ShowDist indicates whether to ask ElasticSearch to return a probability
	// distribution for all characters
	ShowDist *bool `structs:"show_distribution,omitempty"`
}

// StringStats creates a new "string_stats" aggregation with the provided name
// and on the provided field.
func StringStats(name, field string) *StringStatsAgg {
	return &StringStatsAgg{
		BaseAgg: newBaseAgg("string_stats", name, field),
	}
}

// Missing sets the value to provide for records missing a value for the field.
func (agg *StringStatsAgg) Missing(val interface{}) *StringStatsAgg {
	agg.Miss = val
	return agg
}

// ShowDistribution sets whether to show the probability distribution for all
// characters
func (agg *StringStatsAgg) ShowDistribution(b bool) *StringStatsAgg {
	agg.ShowDist = &b
	return agg
}

// Map returns a map representation of the aggregation, thus implementing the
// Mappable interface.
func (agg *StringStatsAgg) Map() map[string]interface{} {
	return map[string]interface{}{
		agg.apiName: structs.Map(agg),
	}
}

// ---------------------------------------------------------------------------//

// TopHitsAgg represents an aggregation of type "top_hits", as described
// in https://www.elastic.co/guide/en/elasticsearch/reference/
//     current/search-aggregations-metrics-top-hits-aggregation.html
type TopHitsAgg struct {
	name   string
	from   uint64
	size   uint64
	sort   []map[string]interface{}
	source Source
}

// TopHits creates an aggregation of type "top_hits".
func TopHits(name string) *TopHitsAgg {
	return &TopHitsAgg{
		name: name,
	}
}

// Name returns the name of the aggregation.
func (agg *TopHitsAgg) Name() string {
	return agg.name
}

// From sets an offset from the first result to return.
func (agg *TopHitsAgg) From(offset uint64) *TopHitsAgg {
	agg.from = offset
	return agg
}

// Size sets the maximum number of top matching hits to return per bucket (the
// default is 3).
func (agg *TopHitsAgg) Size(size uint64) *TopHitsAgg {
	agg.size = size
	return agg
}

// Sort sets how the top matching hits should be sorted. By default the hits are
// sorted by the score of the main query.
func (agg *TopHitsAgg) Sort(name string, order Order) *TopHitsAgg {
	agg.sort = append(agg.sort, map[string]interface{}{
		name: map[string]interface{}{
			"order": order,
		},
	})

	return agg
}

// SourceIncludes sets the keys to return from the top matching documents.
func (agg *TopHitsAgg) SourceIncludes(keys ...string) *TopHitsAgg {
	agg.source.includes = keys
	return agg
}

// Map returns a map representation of the aggregation, thus implementing the
// Mappable interface.
func (agg *TopHitsAgg) Map() map[string]interface{} {
	innerMap := make(map[string]interface{})

	if agg.from > 0 {
		innerMap["from"] = agg.from
	}
	if agg.size > 0 {
		innerMap["size"] = agg.size
	}
	if len(agg.sort) > 0 {
		innerMap["sort"] = agg.sort
	}
	if len(agg.source.includes) > 0 {
		innerMap["_source"] = agg.source.Map()
	}

	return map[string]interface{}{
		"top_hits": innerMap,
	}
}
