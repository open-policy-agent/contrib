// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package es

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"
)

// Posts is a structure used for serializing/deserializing data in Elasticsearch.
type Post struct {
	Id         string              `json:"id"`
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

type People struct {
	Info Name `json:"info"`
}

type Name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type Stat struct {
	AuthorStat AuthorStatData `json:"authorstat"`
}

type AuthorStatData struct {
	AuthorBio AuthorBioData `json:"authorbio"`
}

type AuthorBioData struct {
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"_doc":{
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
	}
}`

func NewPost(id, author, message, department, email string, clearance int, action, resource string, conditions []map[string]string, likes []map[string]string, followers []People, stats []Stat) *Post {
	post := &Post{}
	post.Id = id
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

func NewESClient() (*elastic.Client, error) {
	return elastic.NewClient()
}

func GetIndexMapping() string {
	return mapping
}

// Elasticsearch queries

func GenerateTermQuery(fieldName string, fieldValue interface{}) *elastic.TermQuery {
	return elastic.NewTermQuery(fieldName, fieldValue).QueryName("TermQuery")

}

func GenerateNestedQuery(path string, query elastic.Query) *elastic.NestedQuery {
	return elastic.NewNestedQuery(path, query).QueryName("NestedQuery")

}

func GenerateBoolFilterQuery(filters []elastic.Query) *elastic.BoolQuery {
	q := elastic.NewBoolQuery()
	for _, filter := range filters {
		q = q.Filter(filter)
	}
	q = q.QueryName("BoolFilterQuery")
	return q

}

func GenerateBoolShouldQuery(queries []elastic.Query) *elastic.BoolQuery {
	q := elastic.NewBoolQuery().QueryName("BoolShouldQuery")
	for _, query := range queries {
		q = q.Should(query)
	}
	return q
}

func GenerateBoolMustNotQuery(fieldName string, fieldValue interface{}) *elastic.BoolQuery {
	q := elastic.NewBoolQuery().QueryName("BoolMustNotQuery")
	q = q.MustNot(elastic.NewTermQuery(fieldName, fieldValue))
	return q
}

func GenerateMatchAllQuery() *elastic.MatchAllQuery {
	return elastic.NewMatchAllQuery().QueryName("MatchAllQuery")
}

func GenerateMatchQuery(fieldName string, fieldValue interface{}) *elastic.MatchQuery {
	return elastic.NewMatchQuery(fieldName, fieldValue).QueryName("MatchQuery")
}

func GenerateQueryStringQuery(fieldName string, fieldValue interface{}) *elastic.QueryStringQuery {
	queryString := fmt.Sprintf("*%s*", fieldValue)
	q := elastic.NewQueryStringQuery(queryString).QueryName("QueryStringQuery")
	q = q.DefaultField(fieldName)
	return q
}

func GenerateRegexpQuery(fieldName string, fieldValue interface{}) *elastic.RegexpQuery {
	return elastic.NewRegexpQuery(fieldName, fieldValue.(string))
}

func GenerateRangeQueryLt(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Lt(val)
}

func GenerateRangeQueryLte(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Lte(val)
}

func GenerateRangeQueryGt(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Gt(val)
}

func GenerateRangeQueryGte(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Gte(val)
}

func ExecuteEsSearch(ctx context.Context, client *elastic.Client, indexName string, query elastic.Query) (*elastic.SearchResult, error) {
	searchResult, err := client.Search().
		Index(indexName).
		Query(query). // specify the query
		Pretty(true). // pretty print request and response JSON
		Do(ctx)       // execute
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func AnalyzeSearchResult(searchResult *elastic.SearchResult) {

	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d posts\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// Deserialize hit
			var t Post
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				panic(err)
			}

			// Print with post
			fmt.Printf("\nPost ID: %s\nAuthor: %s\nMessage: %s\nDepartment: %s\nClearance: %d\n", t.Id, t.Author, t.Message, t.Department, t.Clearance)
		}
	} else {
		// No hits
		fmt.Print("Found no posts\n")
	}
}

func GetPrettyResult(searchResult *elastic.SearchResult) []Post {

	result := []Post{}
	if searchResult.Hits.TotalHits > 0 {
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// Deserialize hit
			var t Post
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				panic(err)
			}
			result = append(result, t)
		}
	}
	return result
}
