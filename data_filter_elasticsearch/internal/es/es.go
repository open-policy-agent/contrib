// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package es

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic"
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
	return elastic.NewClient()
}

// GetIndexMapping returns Elasticsearch mapping.
func GetIndexMapping() string {
	return mapping
}

// ExecuteEsSearch executes ES query.
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

// GetPrettyESResult returns formatted ES results.
func GetPrettyESResult(searchResult *elastic.SearchResult) []Post {

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
