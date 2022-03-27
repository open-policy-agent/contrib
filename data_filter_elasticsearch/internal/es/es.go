// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aquasecurity/esquery"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"io"
	"log"
	"strings"
	"time"
)

// Post is a structure used for serializing/deserializing data in Elasticsearch.
type Post struct {
	ID         string              `json:"id"`
	Author     string              `json:"author"`
	Message    string              `json:"message"`
	Department string              `json:"department"`
	Email      string              `json:"email"`
	Clearance  int                 `json:"clearance"`
	Action     string              `json:"action"`
	Resource   string              `json:"resource"`
	Conditions []map[string]string `json:"conditions"`
	Likes      []map[string]string `json:"likes"`
	Followers  []People            `json:"followers"`
	Stats      []Stat              `json:"stats"`
}

// People describes a person.
type People struct {
	Info Name `json:"info"`
}

// Name describes a person's first and last name.
type Name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

// Stat decribes author's stats.
type Stat struct {
	AuthorStat AuthorStatData `json:"authorstat"`
}

// AuthorStatData decribes author's stat data.
type AuthorStatData struct {
	AuthorBio AuthorBioData `json:"authorbio"`
}

// AuthorBioData describes author's bio.
type AuthorBioData struct {
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
}

// InnerHit is a single search result hit
type InnerHit struct {
	Source Post `json:"_source"`
}

// Total number of search result hits
type Total struct {
	Value int `json:"value"`
}

// Hits are the search result hits + metadata like total hits number
type Hits struct {
	Hits  []InnerHit `json:"hits"`
	Total Total      `json:"total"`
}

// SearchResult is Elasticsearch's search query result
type SearchResult struct {
	Hits Hits `json:"hits"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		 "properties": {
			 "id": {
				 "type": "keyword"
			 },
			 "author": {
				 "type": "keyword"
			 },
			 "message": {
				 "type": "text",
				 "fields": {
					  "raw": {
						  "type": "keyword"
					  }
				 }
			 },
			 "department": {
				 "type": "keyword"
			 },
			 "email": {
				 "type": "keyword"
			 },
			 "clearance": {
				 "type": "integer"
			 },
			 "action": {
				 "type": "keyword"
			 },
			 "resource": {
				 "type": "keyword"
			 },
			 "conditions": {
				 "type": "nested",
				 "properties": {
					  "type": {
						  "type": "keyword"
					  },
					  "field": {
						  "type": "keyword"
					  },
					  "value": {
						  "type": "keyword"
					  }
				 }
			 },
			 "likes": {
				 "type": "nested",
				 "properties": {
					  "name": {
						  "type": "keyword"
					  }
				 }
			 },
			 "followers": {
				 "type": "nested",
				 "properties": {
					  "info": {
						  "type": "nested",
						  "properties": {
							  "first": {"type": "keyword"},
							  "last":  {"type": "keyword"}
						  }
					  }
				 }
			 },
			 "stats": {
				 "type": "nested",
				 "properties": {
					  "authorstat": {
						  "type": "nested",
						  "properties": {
							  "authorbio": {
								   "type": "nested",
								   "properties": {
									   "country": {"type": "keyword"},
									   "state":   {"type": "keyword"},
									   "city":    {"type": "keyword"}
								   }
							  }
						  }
					  }
				 }
			 }
		 }
	}
}`

// NewPost returns a post.
func NewPost(id, author, message, department, email string, clearance int, action, resource string, conditions []map[string]string, likes []map[string]string, followers []People, stats []Stat) *Post {
	post := &Post{}
	post.ID = id
	post.Author = author
	post.Message = message
	post.Department = department
	post.Email = email
	post.Clearance = clearance
	post.Action = action
	post.Resource = resource
	post.Conditions = conditions
	post.Likes = likes
	post.Followers = followers
	post.Stats = stats

	return post
}

// NewESClient returns an Elasticsearch client.
func NewESClient() (*elastic.Client, error) {
	return elastic.NewDefaultClient()
}

// GetIndexMapping returns Elasticsearch mapping.
func GetIndexMapping() string {
	return mapping
}

// Elasticsearch queries

// GenerateTermQuery returns an ES Term Query.
func GenerateTermQuery(fieldName string, fieldValue interface{}) *esquery.TermQuery {
	return esquery.Term(fieldName, fieldValue)

}

// GenerateNestedQuery returns an ES Nested Query.
func GenerateNestedQuery(path string, query esquery.Mappable) *esquery.CustomQueryMap {
	return esquery.CustomQuery(map[string]interface{}{"nested": map[string]interface{}{
		"path":  path,
		"query": query.Map()}})
}

// GenerateBoolFilterQuery returns an ES Filter Bool Query.
func GenerateBoolFilterQuery(filters []esquery.Mappable) *esquery.BoolQuery {
	q := esquery.Bool()
	for _, filter := range filters {
		q = q.Filter(filter)
	}
	return q

}

// GenerateBoolShouldQuery returns an ES Should Bool Query.
func GenerateBoolShouldQuery(queries []esquery.Mappable) *esquery.BoolQuery {
	q := esquery.Bool()
	for _, query := range queries {
		q = q.Should(query)
	}
	return q

}

// GenerateBoolMustNotQuery returns an ES Must Not Bool Query.
func GenerateBoolMustNotQuery(fieldName string, fieldValue interface{}) *esquery.BoolQuery {
	q := esquery.Bool()
	q = q.MustNot(esquery.Term(fieldName, fieldValue))
	return q

}

// GenerateMatchAllQuery returns an ES MatchAll Query.
func GenerateMatchAllQuery() *esquery.MatchAllQuery {
	return esquery.MatchAll()
}

// GenerateMatchQuery returns an ES Match Query.
func GenerateMatchQuery(fieldName string, fieldValue interface{}) *esquery.MatchQuery {
	return esquery.Match(fieldName, fieldValue)
}

// GenerateQueryStringQuery returns an ES Query String Query.
func GenerateQueryStringQuery(fieldName string, fieldValue interface{}) *esquery.CustomQueryMap {
	return esquery.CustomQuery(map[string]interface{}{"query_string": map[string]interface{}{
		"query":         fmt.Sprintf("*%s*", fieldValue),
		"default_field": fieldName}})
}

// GenerateRegexpQuery returns an ES Regexp Query.
func GenerateRegexpQuery(fieldName string, fieldValue interface{}) *esquery.RegexpQuery {
	return esquery.Regexp(fieldName, fieldValue.(string))
}

// GenerateRangeQueryLt returns an ES Less Than Range Query.
func GenerateRangeQueryLt(fieldName string, val interface{}) *esquery.RangeQuery {
	return esquery.Range(fieldName).Lt(val)
}

// GenerateRangeQueryLte returns an ES Less Than or Equal Range Query.
func GenerateRangeQueryLte(fieldName string, val interface{}) *esquery.RangeQuery {
	return esquery.Range(fieldName).Lte(val)
}

// GenerateRangeQueryGt returns an ES Greater Than Range Query.
func GenerateRangeQueryGt(fieldName string, val interface{}) *esquery.RangeQuery {
	return esquery.Range(fieldName).Gt(val)
}

// GenerateRangeQueryGte returns an ES Greater Than or Equal Range Query.
func GenerateRangeQueryGte(fieldName string, val interface{}) *esquery.RangeQuery {
	return esquery.Range(fieldName).Gte(val)
}

// ExecuteEsSearch executes ES query.
func ExecuteEsSearch(ctx context.Context, client *elastic.Client, indexName string, query esquery.Mappable) ([]byte, error) {
	queryStr, err := esquery.Query(query).MarshalJSON()
	if err != nil {
		return nil, err
	}

	searchResult, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithPretty(),
		client.Search.WithBody(strings.NewReader(string(queryStr))),
		client.Search.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer searchResult.Body.Close()

	if searchResult.StatusCode != 200 {
		fmt.Println("Search failed: " + searchResult.String())
		return nil, err
	}

	bytes, err := io.ReadAll(searchResult.Body)
	if err != nil {
		fmt.Println("Failed reading search result" + searchResult.String())
		return nil, err
	}
	return bytes, nil
}

// GetPrettyESResult returns formatted ES results.
func GetPrettyESResult(searchResultBytes []byte) []Post {

	result := []Post{}
	var searchResult SearchResult
	err := json.Unmarshal(searchResultBytes, &searchResult)
	if err != nil {
		panic(err)
	}
	fmt.Println(searchResult)
	if searchResult.Hits.Total.Value > 0 {
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			result = append(result, hit.Source)
		}
	}
	return result
}

// IndexPosts indexes posts to ES.
func IndexPosts(ctx context.Context, client *elastic.Client, indexName string, posts []*Post) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,        // The default index name
		Client:        client,           // The Elasticsearch client
		NumWorkers:    3,                // The number of worker goroutines
		FlushBytes:    5e+6,             // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}
	for _, a := range posts {
		// Prepare the data payload: encode post to JSON
		data, err := json.Marshal(a)
		if err != nil {
			log.Fatalf("Cannot encode post %s: %s", a.ID, err)
		}
		// Add an item to the BulkIndexer
		err = bi.Add(
			ctx,
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",
				// DocumentID is the (optional) document ID
				DocumentID: a.ID,
				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(data),
				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
			panic(err)
		}
	}
	// Close the indexer
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
		panic(err)
	}
	biStats := bi.Stats()
	if biStats.NumFailed > 0 {
		log.Printf("Failed to index %d documents", biStats.NumFailed)
		panic(err)
	}
}
