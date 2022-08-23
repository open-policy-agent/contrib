package esquery

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// CustomQueryMap represents an arbitrary query map for custom queries.
type CustomQueryMap map[string]interface{}

// CustomQuery generates a custom request of type "query" from an arbitrary map
// provided by the user. It is useful for issuing a search request with a syntax
// that is not yet supported by the library. CustomQuery values are versatile,
// they can either be used as parameters for the library's Query function, or
// standlone by invoking their Run method.
func CustomQuery(m map[string]interface{}) *CustomQueryMap {
	q := CustomQueryMap(m)
	return &q
}

// Map returns the custom query as a map[string]interface{}, thus implementing
// the Mappable interface.
func (m *CustomQueryMap) Map() map[string]interface{} {
	return map[string]interface{}(*m)
}

// Run executes the custom query using the provided ElasticSearch client. Zero
// or more search options can be provided as well. It returns the standard
// Response type of the official Go client.
func (m *CustomQueryMap) Run(
	api *elasticsearch.Client,
	o ...func(*esapi.SearchRequest),
) (res *esapi.Response, err error) {
	return Search().Query(m).Run(api, o...)
}

//----------------------------------------------------------------------------//

// CustomAggMap represents an arbitrary aggregation map for custom aggregations.
type CustomAggMap struct {
	name string
	agg  map[string]interface{}
}

// CustomAgg generates a custom aggregation from an arbitrary map provided by
// the user.
func CustomAgg(name string, m map[string]interface{}) *CustomAggMap {
	return &CustomAggMap{
		name: name,
		agg:  m,
	}
}

// Name returns the name of the aggregation
func (agg *CustomAggMap) Name() string {
	return agg.name
}

// Map returns a map representation of the custom aggregation, thus implementing
// the Mappable interface
func (agg *CustomAggMap) Map() map[string]interface{} {
	return agg.agg
}
