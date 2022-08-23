package esquery

import (
	"bytes"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// DeleteRequest represents a request to ElasticSearch's Delete By Query API,
// described in
// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-delete-by-query.html
type DeleteRequest struct {
	index []string
	query Mappable
}

// Delete creates a new DeleteRequest object, to be filled via method chaining.
func Delete() *DeleteRequest {
	return &DeleteRequest{}
}

// Index sets the index names for the request
func (req *DeleteRequest) Index(index ...string) *DeleteRequest {
	req.index = index
	return req
}

// Query sets a query for the request.
func (req *DeleteRequest) Query(q Mappable) *DeleteRequest {
	req.query = q
	return req
}

// Run executes the request using the provided ElasticSearch client.
func (req *DeleteRequest) Run(
	api *elasticsearch.Client,
	o ...func(*esapi.DeleteByQueryRequest),
) (res *esapi.Response, err error) {
	return req.RunDelete(api.DeleteByQuery, o...)
}

// RunDelete is the same as the Run method, except that it accepts a value of
// type esapi.DeleteByQuery (usually this is the DeleteByQuery field of an
// elasticsearch.Client object). Since the ElasticSearch client does not provide
// an interface type for its API (which would allow implementation of mock
// clients), this provides a workaround. The Delete function in the ES client is
// actually a field of a function type.
func (req *DeleteRequest) RunDelete(
	del esapi.DeleteByQuery,
	o ...func(*esapi.DeleteByQueryRequest),
) (res *esapi.Response, err error) {
	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(map[string]interface{}{
		"query": req.query.Map(),
	})
	if err != nil {
		return nil, err
	}

	return del(req.index, &b, o...)
}
