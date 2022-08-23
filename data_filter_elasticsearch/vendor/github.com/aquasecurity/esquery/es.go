// Package esquery provides a non-obtrusive, idiomatic and easy-to-use query
// and aggregation builder for the official Go client
// (https://github.com/elastic/go-elasticsearch) for the ElasticSearch
// database (https://www.elastic.co/products/elasticsearch).
//
// esquery alleviates the need to use extremely nested maps
// (map[string]interface{}) and serializing queries to JSON manually. It also
// helps eliminating common mistakes such as misspelling query types, as
// everything is statically typed.
//
// Using `esquery` can make your code much easier to write, read and maintain,
// and significantly reduce the amount of code you write.
//
//
//
// Usage
//
//
//
// esquery provides a method chaining-style API for building and executing
// queries and aggregations. It does not wrap the official Go client nor does it
// require you to change your existing code in order to integrate the library.
// Queries can be directly built with `esquery`, and executed by passing an
// `*elasticsearch.Client` instance (with optional search parameters). Results
// are returned as-is from the official client (e.g. `*esapi.Response` objects).
//
// Getting started is extremely simple:
//
//     package main
//
//     import (
//         "context"
//         "log"
//
//         "github.com/aquasecurity/esquery"
//         "github.com/elastic/go-elasticsearch/v7"
//     )
//
//     func main() {
//         // connect to an ElasticSearch instance
//         es, err := elasticsearch.NewDefaultClient()
//         if err != nil {
//             log.Fatalf("Failed creating client: %s", err)
//         }
//
//         // run a boolean search query
//         qRes, err := esquery.Query(
//             esquery.
//                 Bool().
//                 Must(esquery.Term("title", "Go and Stuff")).
//                 Filter(esquery.Term("tag", "tech")),
//             ).Run(
//                 es,
//                 es.Search.WithContext(context.TODO()),
//                 es.Search.WithIndex("test"),
//             )
//         if err != nil {
//             log.Fatalf("Failed searching for stuff: %s", err)
//         }
//
//         defer qRes.Body.Close()
//
//         // run an aggregation
//         aRes, err := esquery.Aggregate(
//             esquery.Avg("average_score", "score"),
//             esquery.Max("max_score", "score"),
//         ).Run(
//             es,
//             es.Search.WithContext(context.TODO()),
//             es.Search.WithIndex("test"),
//         )
//         if err != nil {
//             log.Fatalf("Failed searching for stuff: %s", err)
//         }
//
//         defer aRes.Body.Close()
//
//         // ...
//     }
//
//
//
// Notes
//
//
//
//* esquery currently supports version 7 of the ElasticSearch Go client.
//* The library cannot currently generate "short queries". For example,
//  whereas ElasticSearch can accept this:
//
//     { "query": { "term": { "user": "Kimchy" } } }
//
// The library will always generate this:
//
//     { "query": { "term": { "user": { "value": "Kimchy" } } } }
//
// This is also true for queries such as "bool", where fields like "must" can
// either receive one query object, or an array of query objects. `esquery` will
// generate an array even if there's only one query object.
package esquery

// Mappable is the interface implemented by the various query and aggregation
// types provided by the package. It allows the library to easily transform the
// different queries to "generic" maps that can be easily encoded to JSON.
type Mappable interface {
	Map() map[string]interface{}
}

// Aggregation is an interface that each aggregation type must implement. It
// is simply an extension of the Mappable interface to include a Named function,
// which returns the name of the aggregation.
type Aggregation interface {
	Mappable
	Name() string
}
