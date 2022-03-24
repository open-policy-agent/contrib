// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/api"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/es"
	"io"
	"io/ioutil"
	"strings"
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
	exists, err := client.Indices.Exists([]string{indexName},
		client.Indices.Exists.WithContext(ctx))
	if err != nil {
		panic(err)
	}

	if exists.StatusCode != 200 {
		// Create a new index.
		createIndex, err := client.Indices.Create(indexName,
			client.Indices.Create.WithBody(strings.NewReader(es.GetIndexMapping())),
			client.Indices.Create.WithContext(ctx))
		if err != nil {
			panic(err)
		}
		if createIndex.StatusCode != 200 {
			panic("Index creation not acknowledged")
		}
		io.Copy(ioutil.Discard, createIndex.Body)
		createIndex.Body.Close()
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

	// Post-2
	post2 := es.NewPost("post2", "bob", "My second post", "dev", "bob@abc.com", 2, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)

	// Post-3
	post3 := es.NewPost("post3", "charlie", "Hello world", "it", "charlie@xyz.com", 1, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)

	// Post-4
	post4 := es.NewPost("post4", "alice", "Hii world", "hr", "alice@xyz.com", 3, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)

	// Post-5
	post5 := es.NewPost("post5", "ben", "Hii from Ben", "ceo", "ben@opa.com", 10, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)

	// Post-6
	post6 := es.NewPost("post6", "ken", "Hii form Ken", "ceo", "ken@opa.com", 5, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)

	// Post-7
	post7 := es.NewPost("post7", "john", "OPA Good", "dev", "john@blah.com", 6, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)

	// Post-8
	post8 := es.NewPost("post8", "ben", "This is OPA's time", "ceo", "ben@opa.com", 10, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)

	// Post-9
	post9 := es.NewPost("post9", "jane", "Hello from Jane", "it", "jane@opa.org", 7, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)

	// Post-10: Nested Query 1 level
	testLikes := make(map[string]string)
	testLikes["name"] = "bob"
	testLikesMap = append(testLikesMap, testLikes)
	post10 := es.NewPost("post10", "ross", "Hello from Ross", "it", "ross@opal.eu", 9, "read", "", testConditionMap, testLikesMap, testFollowers, testStats)

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

	es.IndexPosts(ctx, client, indexName, []*es.Post{post1, post2, post3, post4, post5, post6, post7, post8, post9, post10, post11, post12})
}
