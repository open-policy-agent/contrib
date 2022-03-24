package esquery

import (
	"bytes"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// CountRequest represents a request to get the number of matches for a search
// query, as described in:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-count.html
type CountRequest struct {
	query Mappable
}

// Count creates a new count request with the provided query.
func Count(q Mappable) *CountRequest {
	return &CountRequest{
		query: q,
	}
}

// Map returns a map representation of the request, thus implementing the
// Mappable interface.
func (req *CountRequest) Map() map[string]interface{} {
	return map[string]interface{}{
		"query": req.query.Map(),
	}
}

// Run executes the request using the provided ElasticCount client. Zero or
// more search options can be provided as well. It returns the standard Response
// type of the official Go client.
func (req *CountRequest) Run(
	api *elasticsearch.Client,
	o ...func(*esapi.CountRequest),
) (res *esapi.Response, err error) {
	return req.RunCount(api.Count, o...)
}

// RunCount is the same as the Run method, except that it accepts a value of
// type esapi.Count (usually this is the Count field of an elasticsearch.Client
// object). Since the ElasticCount client does not provide an interface type
// for its API (which would allow implementation of mock clients), this provides
// a workaround. The Count function in the ES client is actually a field of a
// function type.
func (req *CountRequest) RunCount(
	count esapi.Count,
	o ...func(*esapi.CountRequest),
) (res *esapi.Response, err error) {
	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(req.Map())
	if err != nil {
		return nil, err
	}

	opts := append([]func(*esapi.CountRequest){count.WithBody(&b)}, o...)

	return count(opts...)
}
