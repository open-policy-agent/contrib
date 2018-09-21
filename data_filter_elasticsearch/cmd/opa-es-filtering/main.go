// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/api"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/es"
)

func main() {
	ctx := context.Background()

	// Create an ES client.
	client, err := es.NewESClient()
	if err != nil {
		panic(err)
	}

	indexName := "posts"

	// Check if a specified index exists.
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(indexName).BodyString(es.GetIndexMapping()).Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			panic("Index creation not acknowledged")
		}
	}

	// Index posts.
	createTestPosts(ctx, client, indexName)

	// Start server.
	if err := api.New(client, indexName).Run(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Shutting down.")
}

func createTestPosts(ctx context.Context, client *elastic.Client, indexName string) {

	testLikesMap := []map[string]string{}
	testFollowers := []es.People{}
	testStats := []es.Stat{}
	testConditionMap := []map[string]string{}

	// Post-1
	post1 := es.NewPost("post1", "bob", "My first post", "dev", "bob@abc.com", 2, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)
	indexPost(ctx, client, indexName, post1)

	// Post-2
	post2 := es.NewPost("post2", "bob", "My second post", "dev", "bob@abc.com", 2, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)
	indexPost(ctx, client, indexName, post2)

	// Post-3
	post3 := es.NewPost("post3", "charlie", "Hello world", "it", "charlie@xyz.com", 1, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)
	indexPost(ctx, client, indexName, post3)

	// Post-4
	post4 := es.NewPost("post4", "alice", "Hii world", "hr", "alice@xyz.com", 3, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)
	indexPost(ctx, client, indexName, post4)

	// Post-5
	post5 := es.NewPost("post5", "ben", "Hii from Ben", "ceo", "ben@opa.com", 10, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)
	indexPost(ctx, client, indexName, post5)

	// Post-6
	post6 := es.NewPost("post6", "ken", "Hii form Ken", "ceo", "ken@opa.com", 5, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)
	indexPost(ctx, client, indexName, post6)

	// Post-7
	post7 := es.NewPost("post7", "john", "OPA Good", "dev", "john@blah.com", 6, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)
	indexPost(ctx, client, indexName, post7)

	// Post-8
	post8 := es.NewPost("post8", "ben", "This is OPA's time", "ceo", "ben@opa.com", 10, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)
	indexPost(ctx, client, indexName, post8)

	// Post-9
	post9 := es.NewPost("post9", "jane", "Hello from Jane", "it", "jane@opa.org", 7, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)
	indexPost(ctx, client, indexName, post9)

	// Post-10: Nested Query 1 level
	testLikes := make(map[string]string)
	testLikes["name"] = "bob"
	testLikesMap = append(testLikesMap, testLikes)
	post10 := es.NewPost("post10", "ross", "Hello from Ross", "it", "ross@opal.eu", 9, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)
	indexPost(ctx, client, indexName, post10)

	// Post-11: Nested Query 2 levels
	testName := es.Name{
		First: "bob",
		Last:  "doe",
	}
	testFollower := es.People{
		Info: testName,
	}
	testFollowers = append(testFollowers, testFollower)
	post11 := es.NewPost("post11", "rach", "Hello from Rach", "it", "rach@opal.eu", 9, "read", "", testConditionMap, []map[string]string{}, testFollowers, testStats)
	indexPost(ctx, client, indexName, post11)

	// Post-12: Nested Query 3 levels
	authorBio := es.AuthorBioData{
		Country: "US",
		State:   "CA",
		City:    "San Fran",
	}

	authorStat := es.AuthorStatData{
		AuthorBio: authorBio,
	}

	stat := es.Stat{
		AuthorStat: authorStat,
	}
	testStats = append(testStats, stat)

	post12 := es.NewPost("post12", "chan", "Hello from Chan", "it", "chan@opal.eu", 9, "read", "cfgmgmt:nodes", testConditionMap, []map[string]string{}, []es.People{}, testStats)
	indexPost(ctx, client, indexName, post12)
}

func indexPost(ctx context.Context, client *elastic.Client, indexName string, post *es.Post) {
	_, err := client.Index().
		Index(indexName).
		Type("_doc").
		Id(post.ID).
		BodyJson(post).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
}
