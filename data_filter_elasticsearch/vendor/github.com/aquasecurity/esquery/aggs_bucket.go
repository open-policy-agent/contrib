package esquery

//----------------------------------------------------------------------------//

// TermsAggregation represents an aggregation of type "terms", as described in
// https://www.elastic.co/guide/en/elasticsearch/reference/current/
//      search-aggregations-bucket-terms-aggregation.html
type TermsAggregation struct {
	name        string
	field       string
	size        *uint64
	shardSize   *float64
	showTermDoc *bool
	aggs        []Aggregation
	order       map[string]string
	include     []string
}

// TermsAgg creates a new aggregation of type "terms". The method name includes
// the "Agg" suffix to prevent conflict with the "terms" query.
func TermsAgg(name, field string) *TermsAggregation {
	return &TermsAggregation{
		name:  name,
		field: field,
	}
}

// Name returns the name of the aggregation.
func (agg *TermsAggregation) Name() string {
	return agg.name
}

// Size sets the number of term buckets to return.
func (agg *TermsAggregation) Size(size uint64) *TermsAggregation {
	agg.size = &size
	return agg
}

// ShardSize sets how many terms to request from each shard.
func (agg *TermsAggregation) ShardSize(size float64) *TermsAggregation {
	agg.shardSize = &size
	return agg
}

// ShowTermDocCountError sets whether to show an error value for each term
// returned by the aggregation which represents the worst case error in the
// document count.
func (agg *TermsAggregation) ShowTermDocCountError(b bool) *TermsAggregation {
	agg.showTermDoc = &b
	return agg
}

// Aggs sets sub-aggregations for the aggregation.
func (agg *TermsAggregation) Aggs(aggs ...Aggregation) *TermsAggregation {
	agg.aggs = aggs
	return agg
}

// Order sets the sort for terms agg
func (agg *TermsAggregation) Order(order map[string]string) *TermsAggregation {
	agg.order = order
	return agg
}

// Include filter the values for  buckets
func (agg *TermsAggregation) Include(include ...string) *TermsAggregation {
	agg.include = include
	return agg
}

// Map returns a map representation of the aggregation, thus implementing the
// Mappable interface.
func (agg *TermsAggregation) Map() map[string]interface{} {
	innerMap := map[string]interface{}{
		"field": agg.field,
	}

	if agg.size != nil {
		innerMap["size"] = *agg.size
	}
	if agg.shardSize != nil {
		innerMap["shard_size"] = *agg.shardSize
	}
	if agg.showTermDoc != nil {
		innerMap["show_term_doc_count_error"] = *agg.showTermDoc
	}
	if agg.order != nil {
		innerMap["order"] = agg.order
	}

	if agg.include != nil {
		if len(agg.include) <= 1 {
			innerMap["include"] = agg.include[0]
		} else {
			innerMap["include"] = agg.include
		}

	}

	outerMap := map[string]interface{}{
		"terms": innerMap,
	}
	if len(agg.aggs) > 0 {
		subAggs := make(map[string]map[string]interface{})
		for _, sub := range agg.aggs {
			subAggs[sub.Name()] = sub.Map()
		}
		outerMap["aggs"] = subAggs
	}

	return outerMap
}
