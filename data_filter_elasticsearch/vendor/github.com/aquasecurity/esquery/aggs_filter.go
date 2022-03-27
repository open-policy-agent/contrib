package esquery

type FilterAggregation struct {
	name   string
	filter Mappable
	aggs   []Aggregation
}

// Filter creates a new aggregation of type "filter". The method name includes
// the "Agg" suffix to prevent conflict with the "filter" query.
func FilterAgg(name string, filter Mappable) *FilterAggregation {
	return &FilterAggregation{
		name:   name,
		filter: filter,
	}
}

// Name returns the name of the aggregation.
func (agg *FilterAggregation) Name() string {
	return agg.name
}

// Filter sets the filter items
func (agg *FilterAggregation) Filter(filter Mappable) *FilterAggregation {
	agg.filter = filter
	return agg
}

// Aggs sets sub-aggregations for the aggregation.
func (agg *FilterAggregation) Aggs(aggs ...Aggregation) *FilterAggregation {
	agg.aggs = aggs
	return agg
}

func (agg *FilterAggregation) Map() map[string]interface{} {
	outerMap := map[string]interface{}{
		"filter": agg.filter.Map(),
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
